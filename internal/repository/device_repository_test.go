package repository

import (
	"context"
	"testing"
)

var ctx = context.Background()
var device = Device{
	ClientID: "aa:bb:cc:dd:ee:ff",
}
var id int64

func TestAddDevice(t *testing.T) {
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
}

func TestFindDeviceByClientID(t *testing.T) {
	foundDevice, err := deviceRepo.FindDeviceByClientID(ctx, device.ClientID)
	if err != nil {
		t.Fatalf("Failed to find device by client ID: %v", err)
	}
	if foundDevice.ID != id {
		t.Errorf("Expected to find device with ID %d, got %v", id, foundDevice)
	}
}

func TestFindDevices(t *testing.T) {
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
}

func TestDeleteDevice(t *testing.T) {
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
}
