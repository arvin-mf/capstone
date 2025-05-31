package repository

import (
	"context"
	"testing"
)

func TestDeviceOperations(t *testing.T) {
	testDB.Exec("TRUNCATE TABLE devices")

	ctx := context.Background()
	device := Device{ClientID: "aa:bb:cc:dd:ee:ff"}
	var id int64

	t.Run("AddDevice", func(t *testing.T) {
		result, err := deviceRepo.AddDevice(ctx, device)
		if err != nil {
			t.Fatalf("Failed to add device: %v", err)
		}
		id, err = result.LastInsertId()
		if err != nil {
			t.Fatalf("Failed to get last insert ID: %v", err)
		}
		if id != 1 {
			t.Errorf("Expected last insert ID to be 1, got %d", id)
		}
	})

	t.Run("FindDeviceByClientID", func(t *testing.T) {
		foundDevice, err := deviceRepo.FindDeviceByClientID(ctx, device.ClientID)
		if err != nil {
			t.Fatalf("Failed to find device by client ID: %v", err)
		}
		if foundDevice.ID != id {
			t.Errorf("Expected to find device with ID %d, got %v", id, foundDevice)
		}
	})

	t.Run("FindDevices", func(t *testing.T) {
		secondDevice := Device{
			ClientID: "11:22:33:44:55:66",
		}
		_, err := deviceRepo.AddDevice(ctx, secondDevice)
		if err != nil {
			t.Fatalf("Failed to add second device: %v", err)
		}

		devices, err := deviceRepo.FindDevices(ctx)
		if err != nil {
			t.Fatalf("Failed to find devices: %v", err)
		}
		deviceCount := len(devices)
		if deviceCount != 2 {
			t.Errorf("Expected to find 2 devices, got %d", deviceCount)
		}
		if devices[1].ClientID != secondDevice.ClientID {
			t.Errorf("Expected client ID %s, got %s", device.ClientID, devices[1].ClientID)
		}
		if devices[0].ID != id {
			t.Errorf("Expected ID %d, got %d", device.ID, devices[0].ID)
		}
		if devices[0].DeviceStatus != false || devices[1].DeviceStatus != false {
			t.Errorf("Expected device status false, got %t", devices[0].DeviceStatus)
		}
	})

	t.Run("UpdateDeviceStatus", func(t *testing.T) {
		deviceToUpdate := Device{
			ID:           id,
			DeviceStatus: true,
		}
		_, err := deviceRepo.UpdateDeviceStatus(ctx, deviceToUpdate)
		if err != nil {
			t.Fatalf("Failed to update device status: %v", err)
		}

		devices, err := deviceRepo.FindDevices(ctx)
		if err != nil {
			t.Fatalf("Failed to find devices after status update: %v", err)
		}
		if !devices[0].DeviceStatus {
			t.Errorf("Expected device status to be true, got %t", devices[0].DeviceStatus)
		}
	})

	t.Run("DeleteDevice", func(t *testing.T) {
		deviceToDelete := Device{
			ID: id,
		}
		result, err := deviceRepo.DeleteDevice(ctx, deviceToDelete)
		if err != nil {
			t.Fatalf("Failed to delete device: %v", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			t.Fatalf("Failed to get rows affected: %v", err)
		}
		if rowsAffected != 1 {
			t.Errorf("Expected 1 row affected, got %d", rowsAffected)
		}

		foundDevice, err := deviceRepo.FindDeviceByClientID(ctx, device.ClientID)
		if err != nil {
			t.Fatalf("Failed to find device after deletion: %v", err)
		}
		if foundDevice != nil {
			t.Errorf("Expected no device found after deletion, got %v", foundDevice)
		}
	})
}

func TestDeviceSubjectInteractions(t *testing.T) {
	testDB.Exec("TRUNCATE TABLE device_subjects")
	testDB.Exec("TRUNCATE TABLE devices")
	testDB.Exec("TRUNCATE TABLE subjects")

	ctx := context.Background()
	subject := Subject{Name: "Test Subject", IsFatigued: false}

	res, err := testDB.Exec("INSERT INTO devices (client_id) VALUES ('11:22:33:44:55:66')")
	if err != nil {
		t.Fatalf("Failed to insert device: %v", err)
	}
	deviceID, _ := res.LastInsertId()

	res, err = testDB.Exec("INSERT INTO subjects (name) VALUES (?)", subject.Name)
	if err != nil {
		t.Fatalf("Failed to insert subject: %v", err)
	}
	subjectID, _ := res.LastInsertId()

	t.Run("SetDeviceSubject", func(t *testing.T) {
		_, err := deviceRepo.SetDeviceSubject(ctx, SetDeviceSubjectParam{
			SubjectID: subjectID,
			DeviceID:  deviceID,
		})
		if err != nil {
			t.Fatalf("Failed to add device with subject: %v", err)
		}
	})

	t.Run("FindDevicesWithSubject", func(t *testing.T) {
		devices, err := deviceRepo.FindDevicesWithSubject(ctx)
		if err != nil {
			t.Fatalf("Failed to find devices with subject: %v", err)
		}
		if len(devices) == 0 {
			t.Error("Expected to find devices with subject, got none")
		}

		if devices[0].Name.String != subject.Name {
			t.Errorf("Expected subject name %s, got %s", subject.Name, devices[0].Name.String)
		}
		if devices[0].IsFatigued.Valid && devices[0].IsFatigued.Bool != false {
			t.Errorf("Expected subject fatigued status to be false, got %v", devices[0].IsFatigued.Bool)
		}
	})

	t.Run("RemoveDeviceSubject", func(t *testing.T) {
		result, err := deviceRepo.RemoveDeviceSubject(ctx, deviceID)
		if err != nil {
			t.Fatalf("Failed to remove device subject: %v", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			t.Fatalf("Failed to get rows affected: %v", err)
		}
		if rowsAffected != 1 {
			t.Errorf("Expected 1 row affected, got %d", rowsAffected)
		}

		devicesAtLast, err := deviceRepo.FindDevicesWithSubject(ctx)
		if err != nil {
			t.Fatalf("Failed to find devices with subject after removal: %v", err)
		}
		if devicesAtLast[0].Name.Valid != false {
			t.Error("Expected no subject attached after removal, got some")
		}
	})
}
