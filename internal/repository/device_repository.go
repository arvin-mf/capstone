package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type DeviceRepository interface {
	FindDevices(ctx context.Context) ([]Device, error)
	FindDeviceByClientID(ctx context.Context, cID string) (*Device, error)
	AddDevice(ctx context.Context, params Device) (sql.Result, error)
	DeleteDevice(ctx context.Context, params Device) (sql.Result, error)
	FindDevicesWithSubject(ctx context.Context) ([]DeviceWithSubject, error)
}

type deviceRepository struct {
	db *sqlx.DB
}

func NewDeviceRepository(db *sqlx.DB) DeviceRepository {
	return &deviceRepository{
		db: db,
	}
}

const findDevices = `SELECT * FROM devices;`

type Device struct {
	ID        int64     `db:"id"`
	ClientID  string    `db:"client_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *deviceRepository) FindDevices(ctx context.Context) ([]Device, error) {
	var rows []Device
	if err := r.db.Select(&rows, findDevices); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Device{}, nil
		}
		return nil, err
	}
	return rows, nil
}

const findDeviceByClientID = `SELECT id FROM devices WHERE client_id = ?`

func (r *deviceRepository) FindDeviceByClientID(ctx context.Context, cID string) (*Device, error) {
	var row Device
	if err := r.db.Get(&row, findDeviceByClientID, cID); err != nil {
		return nil, err
	}
	return &row, nil
}

const addDevice = `INSERT INTO devices (client_id) VALUES (:client_id);`

func (r *deviceRepository) AddDevice(ctx context.Context, params Device) (sql.Result, error) {
	result, err := r.db.NamedExec(addDevice, params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

const deleteDevice = `DELETE FROM devices WHERE id = :id;`

func (r *deviceRepository) DeleteDevice(ctx context.Context, params Device) (sql.Result, error) {
	result, err := r.db.NamedExec(deleteDevice, params)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected < 1 {
		return nil, errors.New("no rows affected")
	}

	return result, nil
}

const findDevicesWithSubject = `SELECT
	d.id AS device_id,
	ds.subject_id,
	s.name,
	s.is_fatigued,
	ds.created_at
FROM devices d
LEFT JOIN device_subjects ds ON ds.device_id = d.id
LEFT JOIN subjects s ON ds.subject_id = s.id;`

type DeviceWithSubject struct {
	DeviceID   int64          `db:"device_id"`
	SubjectID  sql.NullInt64  `db:"subject_id"`
	Name       sql.NullString `db:"name"`
	IsFatigued sql.NullBool   `db:"is_fatigued"`
	CreatedAt  sql.NullTime   `db:"created_at"`
}

func (r *deviceRepository) FindDevicesWithSubject(ctx context.Context) ([]DeviceWithSubject, error) {
	var rows []DeviceWithSubject
	if err := r.db.Select(&rows, findDevicesWithSubject); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []DeviceWithSubject{}, nil
		}
		fmt.Println(err)
		return nil, err
	}

	return rows, nil
}
