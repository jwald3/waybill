package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jwald3/waybill/internal/domain"
	"github.com/jwald3/waybill/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IncidentReportHandler struct {
	incidentReportService service.IncidentReportService
}

func NewIncidentReportHandler(incidentReportService service.IncidentReportService) *IncidentReportHandler {
	return &IncidentReportHandler{incidentReportService: incidentReportService}
}

var (
	invalidIncidentReportId   = "invalid incident report id"
	invalidIncidentReportPath = "invalid path"
)

// DTOS =======================================================

type IncidentReportCreateRequest struct {
	TripID         *primitive.ObjectID `json:"trip_id"`
	TruckID        *primitive.ObjectID `json:"truck_id"`
	DriverID       *primitive.ObjectID `json:"driver_id"`
	Type           domain.IncidentType `json:"type"`
	Description    string              `json:"description"`
	Date           string              `json:"date"`
	Location       string              `json:"location"`
	DamageEstimate float64             `json:"damage_estimate"`
}

type IncidentReportUpdateRequest struct {
	TripID         *primitive.ObjectID `json:"trip_id"`
	TruckID        *primitive.ObjectID `json:"truck_id"`
	DriverID       *primitive.ObjectID `json:"driver_id"`
	Type           domain.IncidentType `json:"type"`
	Description    string              `json:"description"`
	Date           string              `json:"date"`
	Location       string              `json:"location"`
	DamageEstimate float64             `json:"damage_estimate"`
}

type IncidentReportResponse struct {
	ID             primitive.ObjectID  `json:"id,omitempty"`
	TripID         *primitive.ObjectID `json:"trip_id,omitempty"`
	Trip           *domain.Trip        `json:"trip,omitempty"`
	TruckID        *primitive.ObjectID `json:"truck_id,omitempty"`
	Truck          *domain.Truck       `json:"truck,omitempty"`
	DriverID       *primitive.ObjectID `json:"driver_id,omitempty"`
	Driver         *domain.Driver      `json:"driver,omitempty"`
	Type           domain.IncidentType `json:"type"`
	Description    string              `json:"description"`
	Date           string              `json:"date"`
	Location       string              `json:"location"`
	DamageEstimate float64             `json:"damage_estimate"`
	CreatedAt      primitive.DateTime  `json:"created_at"`
	UpdatedAt      primitive.DateTime  `json:"updated_at"`
}

type ListIncidentReportsResponse struct {
	IncidentReports []IncidentReportResponse `json:"fuel_logs"`
}

func incidentReportRequestToDomainCreate(req IncidentReportCreateRequest) (*domain.IncidentReport, error) {
	return domain.NewIncidentReport(
		req.TripID,
		req.TruckID,
		req.DriverID,
		req.Type,
		req.Description,
		req.Date,
		req.Location,
		req.DamageEstimate,
	)
}

func incidentReportRequestToDomainUpdate(req IncidentReportUpdateRequest) (*domain.IncidentReport, error) {
	if !req.Type.IsValid() {
		return nil, fmt.Errorf("invalid incident report type: %s", req.Type)
	}

	return &domain.IncidentReport{
		TripID:         req.TripID,
		TruckID:        req.TruckID,
		DriverID:       req.DriverID,
		Type:           req.Type,
		Description:    req.Description,
		Date:           req.Date,
		Location:       req.Location,
		DamageEstimate: req.DamageEstimate,
	}, nil
}

func incidentReportDomainToResponse(i *domain.IncidentReport) IncidentReportResponse {
	return IncidentReportResponse{
		ID:             i.ID,
		TripID:         i.TripID,
		Trip:           i.Trip,
		TruckID:        i.TruckID,
		Truck:          i.Truck,
		DriverID:       i.DriverID,
		Driver:         i.Driver,
		Type:           i.Type,
		Description:    i.Description,
		Date:           i.Date,
		Location:       i.Location,
		DamageEstimate: i.DamageEstimate,
		CreatedAt:      i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
	}
}

// =================================================================

func (h *IncidentReportHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req IncidentReportCreateRequest

	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	incidentReport, err := incidentReportRequestToDomainCreate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.incidentReportService.Create(r.Context(), incidentReport); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	WriteJSON(w, http.StatusCreated, incidentReportDomainToResponse(incidentReport))
}

func (h *IncidentReportHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	incidentReport, err := h.incidentReportService.GetById(r.Context(), objectID)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, Response{Error: "incident report not found"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: incidentReportDomainToResponse(incidentReport)})
}

func (h *IncidentReportHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	var req IncidentReportUpdateRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	incidentReport, err := incidentReportRequestToDomainUpdate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	incidentReport.ID = objectID

	if err := h.incidentReportService.Update(r.Context(), incidentReport); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to update user"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: incidentReportDomainToResponse(incidentReport)})
}

func (h *IncidentReportHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.incidentReportService.Delete(r.Context(), objectID); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to delete incident report"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Message: "incident report deleted successfully"})
}

func (h *IncidentReportHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := getQueryIntParam(r, "limit", 10)
	offset := getQueryIntParam(r, "offset", 0)

	incidentReports, err := h.incidentReportService.List(r.Context(), int64(limit), int64(offset))
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to fetch incident reports"})
		return
	}

	incidentReportResponses := make([]IncidentReportResponse, len(incidentReports))
	for i, d := range incidentReports {
		incidentReportResponses[i] = incidentReportDomainToResponse(d)
	}

	WriteJSON(w, http.StatusOK, Response{Data: ListIncidentReportsResponse{IncidentReports: incidentReportResponses}})
}
