package repository

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type InfluxRepository interface {
	Write(ctx context.Context, params InfluxPointParam) error
}

type influxRepository struct {
	write api.WriteAPIBlocking
	// query api.QueryAPI
}

func NewInfluxRepository(ic influxdb2.Client, org, bucket string) InfluxRepository {
	return &influxRepository{
		write: ic.WriteAPIBlocking(org, bucket),
		// query: ic.QueryAPI(org),
	}
}

const (
	clientTagName               string = "client"
	subjectTagName              string = "subject"
	bpmFieldName                string = "bpm"
	bodyTemperatureFieldName    string = "body_temperature"
	ambientTemperatureFieldName string = "ambient_temperature"
	statusFieldName             string = "status"
)

type SubjectStatus bool

const (
	StatusFatigued    SubjectStatus = true
	StatusNotFatigued SubjectStatus = false
)

type InfluxPointParam struct {
	DeviceID           string
	SubjectID          string
	Bpm                float32
	BodyTemperature    float32
	AmbientTemperature float32
	Status             SubjectStatus
}

func (r *influxRepository) Write(ctx context.Context, params InfluxPointParam) error {
	point := write.NewPoint(
		"subject_fatigue",
		map[string]string{
			clientTagName:  params.DeviceID,
			subjectTagName: params.SubjectID,
		},
		map[string]interface{}{
			bpmFieldName:                params.Bpm,
			bodyTemperatureFieldName:    params.BodyTemperature,
			ambientTemperatureFieldName: params.AmbientTemperature,
			statusFieldName:             bool(params.Status),
		},
		time.Now(),
	)

	if err := r.write.WritePoint(ctx, point); err != nil {
		return err
	}

	return nil
}
