package repository

import (
	"context"
	"testing"
)

func TestSubjectOperations(t *testing.T) {
	testDB.Exec("TRUNCATE TABLE subjects")

	ctx := context.Background()
	subject := Subject{Name: "Subject X"}
	var subjectID int64

	t.Run("AddSubject", func(t *testing.T) {
		result, err := subjectRepo.AddSubject(ctx, subject)
		if err != nil {
			t.Fatalf("Failed to add subject: %v", err)
		}
		subjectID, err = result.LastInsertId()
		if err != nil {
			t.Fatalf("Failed to get last insert ID: %v", err)
		}
		if subjectID != 1 {
			t.Errorf("Expected last insert ID to be 1, got %d", subjectID)
		}
	})

	secondSubject := Subject{Name: "Subject Y"}

	t.Run("FindSubjects", func(t *testing.T) {
		_, err := subjectRepo.AddSubject(ctx, secondSubject)
		if err != nil {
			t.Fatalf("Failed to add second subject: %v", err)
		}

		subjects, err := subjectRepo.FindSubjects(ctx)
		if err != nil {
			t.Fatalf("Failed to find subjects: %v", err)
		}
		subjectCount := len(subjects)
		if subjectCount != 2 {
			t.Errorf("Expected to find 2 subjects, got %d", subjectCount)
		}
		if subjects[1].Name != secondSubject.Name {
			t.Errorf("Expected name %s, got %s", secondSubject.Name, subjects[1].Name)
		}
		if subjects[0].ID != subjectID {
			t.Errorf("Expected ID %d, got %d", subject.ID, subjects[0].ID)
		}
		if subjects[0].IsFatigued || subjects[1].IsFatigued {
			t.Error("Expected IsFatigued to be false, but found true")
		}
	})

	t.Run("UpdateSubjectFatiguedStatus", func(t *testing.T) {
		subjectToUpdate := Subject{
			ID:         subjectID,
			IsFatigued: true,
		}

		_, err := subjectRepo.UpdateSubjectFatiguedStatus(ctx, subjectToUpdate)
		if err != nil {
			t.Fatalf("Failed to update subject fatigued status: %v", err)
		}

		subjects, err := subjectRepo.FindSubjects(ctx)
		if err != nil {
			t.Errorf("Failed to find subjects after update: %v", err)
		}
		if !subjects[0].IsFatigued {
			t.Errorf("Expected subject with ID %d to be fatigued, but it is not", subjectID)
		}
	})

	t.Run("DeleteSubject", func(t *testing.T) {
		subjectToDelete := Subject{ID: subjectID}
		result, err := subjectRepo.DeleteSubject(ctx, subjectToDelete)
		if err != nil {
			t.Fatalf("Failed to delete subject: %v", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			t.Fatalf("Failed to get rows affected: %v", err)
		}
		if rowsAffected != 1 {
			t.Errorf("Expected 1 row affected, got %d", rowsAffected)
		}

		subjects, err := subjectRepo.FindSubjects(ctx)
		if err != nil {
			t.Fatalf("Failed to find subjects after deletion: %v", err)
		}
		if len(subjects) != 1 {
			t.Errorf("Expected 1 subject after deletion, got %d", len(subjects))
		}
		if subjects[0].Name != secondSubject.Name {
			t.Errorf("Expected remaining subject to be %s, got %s", secondSubject.Name, subjects[0].Name)
		}
	})
}

func TestSubjectInteractions(t *testing.T) {
	testDB.Exec("TRUNCATE TABLE subjects")
	testDB.Exec("TRUNCATE TABLE devices")
	testDB.Exec("TRUNCATE TABLE device_subjects")

	ctx := context.Background()
	subject := Subject{Name: "Test Subject"}

	res, err := testDB.Exec("INSERT INTO devices (client_id) VALUES ('AA:BB:CC:DD:EE:FF')")
	if err != nil {
		t.Fatalf("Failed to insert device: %v", err)
	}
	deviceID, _ := res.LastInsertId()

	res, err = testDB.Exec("INSERT INTO subjects (name) VALUES (?)", subject.Name)
	if err != nil {
		t.Fatalf("Failed to insert subject: %v", err)
	}
	subjectID, _ := res.LastInsertId()

	testDB.Exec("INSERT INTO device_subjects (device_id, subject_id) VALUES (?, ?)", deviceID, subjectID)

	t.Run("FindSubjectByDeviceID", func(t *testing.T) {
		foundSubject, err := subjectRepo.FindSubjectByDeviceID(ctx, deviceID)
		if err != nil {
			t.Fatalf("Failed to find subject by client ID: %v", err)
		}
		if foundSubject.ID != subjectID {
			t.Errorf("Expected to find subject with ID %d, got %d", subjectID, foundSubject.ID)
		}
		if foundSubject.Name != subject.Name {
			t.Errorf("Expected subject name 'Test Subject', got '%s'", foundSubject.Name)
		}
	})
}
