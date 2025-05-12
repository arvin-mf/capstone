package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type SubjectRepository interface {
	AddSubject(ctx context.Context, params Subject) (sql.Result, error)
	FindSubjects(ctx context.Context) ([]Subject, error)
	DeleteSubject(ctx context.Context, params Subject) (sql.Result, error)
	FindSubjectByDeviceID(ctx context.Context, dID int64) (*Subject, error)
}

type subjectRepository struct {
	db *sqlx.DB
}

func NewSubjectRepository(db *sqlx.DB) SubjectRepository {
	return &subjectRepository{
		db: db,
	}
}

const addSubject = `INSERT INTO subjects (name) VALUES (:name);`

func (r *subjectRepository) AddSubject(ctx context.Context, params Subject) (sql.Result, error) {
	result, err := r.db.NamedExec(addSubject, params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

const findSubjects = `SELECT * FROM subjects;`

type Subject struct {
	ID         int64     `db:"id"`
	Name       string    `db:"name"`
	IsFatigued bool      `db:"is_fatigued"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (r *subjectRepository) FindSubjects(ctx context.Context) ([]Subject, error) {
	var rows []Subject
	err := r.db.Select(&rows, findSubjects)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Subject{}, nil
		}
		return nil, err
	}
	return rows, nil
}

const deleteSubject = `DELETE FROM subjects WHERE id = :id`

func (r *subjectRepository) DeleteSubject(ctx context.Context, params Subject) (sql.Result, error) {
	result, err := r.db.NamedExec(deleteSubject, params)
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

const findSubjectByDeviceID = `SELECT s.* FROM subjects s
INNER JOIN device_subjects ds ON ds.subject_id = s.id
WHERE ds.device_id = ?`

func (r *subjectRepository) FindSubjectByDeviceID(ctx context.Context, dID int64) (*Subject, error) {
	var row Subject
	if err := r.db.Get(&row, findSubjectByDeviceID, dID); err != nil {
		return nil, err
	}
	return &row, nil
}
