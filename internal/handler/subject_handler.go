package handler

import (
	"caps_influx/internal/dto"
	"caps_influx/internal/service"
	"caps_influx/pkg/rsp"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubjectHandler struct {
	serv service.SubjectService
}

func NewSubjectHandler(ss service.SubjectService) *SubjectHandler {
	return &SubjectHandler{
		serv: ss,
	}
}

const (
	subjectsFetchFailed  = "Failed to fetch subjects"
	subjectsFetchSuccess = "Subjects successfully retrieved"
	subjectCreateFailed  = "Failed to create subject"
	subjectCreateSuccess = "Subject successfully created"
	subjectDeleteFailed  = "Failed to delete subject"
	subjectDeleteSuccess = "Subject successfully deleted"
)

func (h *SubjectHandler) GetAllSubjects(ctx *gin.Context) {
	resp, err := h.serv.GetSubjects(ctx)
	if err != nil {
		rsp.FailOrError(ctx, http.StatusInternalServerError, subjectsFetchFailed, err)
		return
	}

	rsp.Success(ctx, http.StatusOK, subjectsFetchSuccess, resp)
}

func (h *SubjectHandler) AddSubject(ctx *gin.Context) {
	var req dto.SubjectCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		rsp.FailOrError(ctx, http.StatusBadRequest, subjectCreateFailed, err)
		return
	}

	if err := h.serv.AddSubject(ctx, req); err != nil {
		rsp.FailOrError(ctx, http.StatusInternalServerError, subjectCreateFailed, err)
		return
	}

	rsp.Success(ctx, http.StatusCreated, subjectCreateSuccess, nil)
}

func (h *SubjectHandler) DeleteSubject(ctx *gin.Context) {
	subjectParam := ctx.Param("subjectId")
	subjectID, err := strconv.Atoi(subjectParam)
	if err != nil {
		rsp.FailOrError(ctx, http.StatusBadRequest, subjectDeleteFailed, err)
		return
	}

	if err := h.serv.DeleteSubject(ctx, int64(subjectID)); err != nil {
		rsp.FailOrError(ctx, http.StatusInternalServerError, subjectDeleteFailed, err)
		return
	}

	rsp.Success(ctx, http.StatusOK, subjectDeleteSuccess, nil)
}

func (h *SubjectHandler) GetSubjectsWithDevice(ctx *gin.Context) {
	resp, err := h.serv.GetSubjectsWithDevice(ctx)
	if err != nil {
		rsp.FailOrError(ctx, http.StatusInternalServerError, subjectsFetchFailed, err)
		return
	}

	rsp.Success(ctx, http.StatusOK, subjectsFetchSuccess, resp)
}
