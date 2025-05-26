package handler

import (
	"caps_influx/internal/dto"
	"caps_influx/internal/service"
	"caps_influx/pkg/rsp"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	serv service.DeviceService
}

func NewDeviceHandler(ds service.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		serv: ds,
	}
}

const (
	devicesFetchFailed  = "Failed to fetch devices"
	devicesFetchSuccess = "Devices successfully retrieved"
	deviceCreateFailed  = "Failed to create device"
	deviceCreateSuccess = "Device successfully created"
	deviceDeleteFailed  = "Failed to delete device"
	deviceDeleteSuccess = "Device successfully deleted"

	deviceSubjectUpdateFailed  = "Failed to update device subject"
	deviceSubjectUpdateSuccess = "Device subject successfully updated"
)

func (h *DeviceHandler) GetAllDevices(ctx *gin.Context) {
	resp, err := h.serv.GetDevices(ctx)
	if err != nil {
		rsp.FailOrError(ctx, http.StatusInternalServerError, devicesFetchFailed, err)
		return
	}

	rsp.Success(ctx, http.StatusOK, devicesFetchSuccess, resp)
}

func (h *DeviceHandler) AddDevice(ctx *gin.Context) {
	var req dto.DeviceCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		rsp.FailOrError(ctx, http.StatusBadRequest, deviceCreateFailed, err)
		return
	}

	if err := h.serv.AddDevice(ctx, req); err != nil {
		rsp.FailOrError(ctx, http.StatusInternalServerError, deviceCreateFailed, err)
		return
	}

	rsp.Success(ctx, http.StatusCreated, deviceCreateSuccess, nil)
}

func (h *DeviceHandler) DeleteDevice(ctx *gin.Context) {
	deviceParam := ctx.Param("deviceId")
	deviceID, err := strconv.Atoi(deviceParam)
	if err != nil {
		rsp.FailOrError(ctx, http.StatusBadRequest, deviceDeleteFailed, err)
		return
	}

	if err := h.serv.DeleteDevice(ctx, int64(deviceID)); err != nil {
		rsp.FailOrError(ctx, http.StatusInternalServerError, deviceDeleteFailed, err)
		return
	}

	rsp.Success(ctx, http.StatusOK, deviceDeleteSuccess, nil)
}

func (h *DeviceHandler) GetDevicesWithSubject(ctx *gin.Context) {
	resp, err := h.serv.GetDevicesWithSubject(ctx)
	if err != nil {
		rsp.FailOrError(ctx, http.StatusInternalServerError, devicesFetchFailed, err)
		return
	}

	rsp.Success(ctx, http.StatusOK, devicesFetchSuccess, resp)
}

func (h *DeviceHandler) SetDeviceSubject(ctx *gin.Context) {
	deviceParam := ctx.Param("deviceId")
	deviceID, err := strconv.Atoi(deviceParam)
	if err != nil {
		rsp.FailOrError(ctx, http.StatusBadRequest, deviceSubjectUpdateFailed, err)
		return
	}

	var req dto.SetDeviceSubjectReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		rsp.FailOrError(ctx, http.StatusBadRequest, deviceSubjectUpdateFailed, err)
		return
	}

	if err := h.serv.SetDeviceSubject(ctx, int64(deviceID), req); err != nil {
		rsp.FailOrError(ctx, http.StatusInternalServerError, deviceSubjectUpdateFailed, err)
		return
	}

	rsp.Success(ctx, http.StatusCreated, deviceSubjectUpdateSuccess, nil)
}
