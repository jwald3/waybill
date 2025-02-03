package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jwald3/waybill/internal/domain"
	"github.com/jwald3/waybill/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MaintenanceLogHandler struct {
	maintenanceLogService service.MaintenanceLogService
}

func NewMaintenanceLogHandler(maintenanceLogService service.MaintenanceLogService) *MaintenanceLogHandler {
	return &MaintenanceLogHandler{maintenanceLogService: maintenanceLogService}
}

var (
	invalidMaintenanceLogId   = "invalid maintenance log id"
	invalidMaintenanceLogPath = "invalid path"
)

// DTOS =======================================================

type MaintenanceLogCreateRequest struct {
	TruckID     *primitive.ObjectID           `json:"truck_id"`
	Date        string                        `json:"date"`
	ServiceType domain.MaintenanceServiceType `json:"service_type"`
	Cost        float64                       `json:"cost"`
	Notes       string                        `json:"notes"`
	Mechanic    string                        `json:"mechanic"`
	Location    string                        `json:"location"`
}

type MaintenanceLogUpdateRequest struct {
	ID          primitive.ObjectID            `json:"id,omitempty"`
	TruckID     *primitive.ObjectID           `json:"truck_id"`
	Date        string                        `json:"date"`
	ServiceType domain.MaintenanceServiceType `json:"service_type"`
	Cost        float64                       `json:"cost"`
	Notes       string                        `json:"notes"`
	Mechanic    string                        `json:"mechanic"`
	Location    string                        `json:"location"`
}

type MaintenanceLogResponse struct {
	ID          primitive.ObjectID            `json:"id,omitempty"`
	TruckID     *primitive.ObjectID           `json:"truck_id,omitempty"`
	Truck       *domain.Truck                 `json:"truck,omitempty"`
	Date        string                        `json:"date"`
	ServiceType domain.MaintenanceServiceType `json:"service_type"`
	Cost        float64                       `json:"cost"`
	Notes       string                        `json:"notes"`
	Mechanic    string                        `json:"mechanic"`
	Location    string                        `json:"location"`
	CreatedAt   primitive.DateTime            `json:"created_at"`
	UpdatedAt   primitive.DateTime            `json:"updated_at"`
}

type ListMaintenanceLogsResponse struct {
	MaintenanceLogs []MaintenanceLogResponse `json:"maintenance_logs"`
}

func maintenanceLogRequestToDomainCreate(req MaintenanceLogCreateRequest) (*domain.MaintenanceLog, error) {
	return domain.NewMaintenanceLog(
		req.TruckID,
		req.Date,
		req.ServiceType,
		req.Notes,
		req.Mechanic,
		req.Location,
		req.Cost,
	)
}

func maintenanceLogRequestToDomainUpdate(req MaintenanceLogUpdateRequest) (*domain.MaintenanceLog, error) {
	if !req.ServiceType.IsValid() {
		return nil, fmt.Errorf("invalid service type provided: %s", req.ServiceType)
	}

	now := time.Now()

	return &domain.MaintenanceLog{
		ID:          req.ID,
		TruckID:     req.TruckID,
		Date:        req.Date,
		ServiceType: req.ServiceType,
		Notes:       req.Notes,
		Mechanic:    req.Mechanic,
		Location:    req.Location,
		Cost:        req.Cost,
		CreatedAt:   primitive.NewDateTimeFromTime(now),
		UpdatedAt:   primitive.NewDateTimeFromTime(now),
	}, nil
}

func maintenanceLogDomainToResponse(m *domain.MaintenanceLog) MaintenanceLogResponse {
	return MaintenanceLogResponse{
		ID:          m.ID,
		TruckID:     m.TruckID,
		Truck:       m.Truck,
		Date:        m.Date,
		ServiceType: m.ServiceType,
		Notes:       m.Notes,
		Mechanic:    m.Mechanic,
		Location:    m.Location,
		Cost:        m.Cost,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// =================================================================

func (h *MaintenanceLogHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req MaintenanceLogCreateRequest

	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	maintenanceLog, err := maintenanceLogRequestToDomainCreate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.maintenanceLogService.Create(r.Context(), maintenanceLog); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	WriteJSON(w, http.StatusCreated, maintenanceLogDomainToResponse(maintenanceLog))
}

func (h *MaintenanceLogHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	maintenanceLog, err := h.maintenanceLogService.GetById(r.Context(), objectID)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, Response{Error: "maintenance log not found"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: maintenanceLogDomainToResponse(maintenanceLog)})
}

func (h *MaintenanceLogHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	var req MaintenanceLogUpdateRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	maintenanceLog, err := maintenanceLogRequestToDomainUpdate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	maintenanceLog.ID = objectID

	if err := h.maintenanceLogService.Update(r.Context(), maintenanceLog); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to update user"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: maintenanceLogDomainToResponse(maintenanceLog)})
}

func (h *MaintenanceLogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.maintenanceLogService.Delete(r.Context(), objectID); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to delete maintenance log"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Message: "maintenance log deleted successfully"})
}

func (h *MaintenanceLogHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := getQueryIntParam(r, "limit", 10)
	offset := getQueryIntParam(r, "offset", 0)

	maintenanceLogs, err := h.maintenanceLogService.List(r.Context(), int64(limit), int64(offset))
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to fetch maintenance logs"})
		return
	}

	maintenanceLogResponses := make([]MaintenanceLogResponse, len(maintenanceLogs))
	for i, d := range maintenanceLogs {
		maintenanceLogResponses[i] = maintenanceLogDomainToResponse(d)
	}

	WriteJSON(w, http.StatusOK, Response{Data: ListMaintenanceLogsResponse{MaintenanceLogs: maintenanceLogResponses}})
}
