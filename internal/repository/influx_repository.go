package repository

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type InfluxRepository interface {
	WritePeriodic(ctx context.Context, params InfluxPeriodicPointParam) error
	WritePerpetual(ctx context.Context, params InfluxPerpetualPointParam) error
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
	deviceTagName               string = "device"
	subjectTagName              string = "subject"
	bpmFieldName                string = "bpm"
	bodyTemperatureFieldName    string = "body_temperature"
	ambientTemperatureFieldName string = "ambient_temperature"
	statusFieldName             string = "status"
	rawEcgFieldName             string = "raw_ecg"
	timestampFieldName          string = "timestamp"
)

type SubjectStatus bool

const (
	StatusFatigued    SubjectStatus = true
	StatusNotFatigued SubjectStatus = false
)

type InfluxPeriodicPointParam struct {
	DeviceID           string
	SubjectID          string
	Bpm                float32
	BodyTemperature    float32
	AmbientTemperature float32
	Status             SubjectStatus
}

func (r *influxRepository) WritePeriodic(ctx context.Context, params InfluxPeriodicPointParam) error {
	point := write.NewPoint(
		"subject_fatigue",
		map[string]string{
			deviceTagName:  params.DeviceID,
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

type InfluxPerpetualPointParam struct {
	DeviceID  string
	SubjectID string
	RawEcg    float32
	Timestamp int64
}

func (r *influxRepository) WritePerpetual(ctx context.Context, params InfluxPerpetualPointParam) error {
	point := write.NewPoint(
		"subject_raw_ecg",
		map[string]string{
			deviceTagName:  params.DeviceID,
			subjectTagName: params.SubjectID,
		},
		map[string]interface{}{
			rawEcgFieldName:    params.RawEcg,
			timestampFieldName: params.Timestamp,
		},
		time.Now(),
	)

	if err := r.write.WritePoint(ctx, point); err != nil {
		return err
	}

	return nil
}
