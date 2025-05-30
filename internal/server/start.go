package server

import (
	"caps_influx/config"
	"caps_influx/internal/handler"
	"caps_influx/internal/repository"
	"caps_influx/internal/service"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func StartEngine(e *gin.Engine, db *sqlx.DB, ic influxdb2.Client, mc mqtt.Client, rc *redis.Client) {
	sh, dh, subscs := initHandler(db, ic, mc, rc)
	go subscs.SubscribePeriodicData()
	go subscs.SubscribePerpetualData()
	go subscs.SubscribeStatusData()
	route(e, sh, dh)
}

func route(r *gin.Engine, sh *handler.SubjectHandler, dh *handler.DeviceHandler) {
	apiRoute(r, sh, dh)
	webRoute(r)
}

func initHandler(db *sqlx.DB, ic influxdb2.Client, mc mqtt.Client, rc *redis.Client) (*handler.SubjectHandler, *handler.DeviceHandler, service.SubscribeService) {
	var (
		subjectRepo = repository.NewSubjectRepository(db)
		deviceRepo  = repository.NewDeviceRepository(db, rc)
		influxRepo  = repository.NewInfluxRepository(
			ic,
			config.GetEnv("INFLUXDB_ORG", ""),
			config.GetEnv("INFLUXDB_BUCKET", ""),
		)
	)

	var (
		subjectServ   = service.NewSubjectService(subjectRepo)
		deviceServ    = service.NewDeviceService(deviceRepo)
		influxServ    = service.NewInfluxService(influxRepo, subjectRepo, deviceRepo)
		subscribeServ = service.NewSubscribeService(mc, influxServ, deviceRepo)
	)

	var (
		subjectHand = handler.NewSubjectHandler(subjectServ)
		deviceHand  = handler.NewDeviceHandler(deviceServ)
	)

	return subjectHand, deviceHand, subscribeServ
}

func apiRoute(r *gin.Engine, sh *handler.SubjectHandler, dh *handler.DeviceHandler) {
	api := r.Group("/api")

	api.GET("/subjects", sh.GetAllSubjects)
	api.POST("/subjects", sh.AddSubject)
	api.DELETE("/subjects/:subjectId", sh.DeleteSubject)

	api.GET("/devices", dh.GetAllDevices)
	api.POST("/devices", dh.AddDevice)
	api.DELETE("/devices/:deviceId", dh.DeleteDevice)
	api.GET("/devices/subjects", dh.GetDevicesWithSubject)
	api.PUT("/devices/:deviceId/subjects", dh.SetDeviceSubject)
}

func webRoute(r *gin.Engine) {
	r.Static("/static", "./web")
	r.LoadHTMLFiles("./web/index.html")
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
}
