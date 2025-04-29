package handler

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jwald3/waybill/internal/domain"
	"github.com/jwald3/waybill/internal/middleware"
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
	invalidIncidentReportId = "invalid incident report id"
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

func incidentReportRequestToDomainCreate(userID primitive.ObjectID, req IncidentReportCreateRequest) (*domain.IncidentReport, error) {
	return domain.NewIncidentReport(
		userID,
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
	claims, ok := r.Context().Value(middleware.UserContextKey).(jwt.MapClaims)
	if !ok {
		WriteJSON(w, http.StatusUnauthorized, Response{Error: "unauthorized"})
		return
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		WriteJSON(w, http.StatusUnauthorized, Response{Error: "invalid user id in token"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, Response{Error: "invalid user id format"})
		return
	}

	var req IncidentReportCreateRequest

	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	incidentReport, err := incidentReportRequestToDomainCreate(userID, req)
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
	claims, ok := r.Context().Value(middleware.UserContextKey).(jwt.MapClaims)
	if !ok {
		WriteJSON(w, http.StatusUnauthorized, Response{Error: "unauthorized"})
		return
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		WriteJSON(w, http.StatusUnauthorized, Response{Error: "invalid user id in token"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, Response{Error: "invalid user id format"})
		return
	}

	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	incidentReport, err := h.incidentReportService.GetById(r.Context(), objectID, userID)
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
	claims, ok := r.Context().Value(middleware.UserContextKey).(jwt.MapClaims)
	if !ok {
		WriteJSON(w, http.StatusUnauthorized, Response{Error: "unauthorized"})
		return
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		WriteJSON(w, http.StatusUnauthorized, Response{Error: "invalid user id in token"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, Response{Error: "invalid user id format"})
		return
	}

	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: invalidIncidentReportId})
		return
	}

	err = h.incidentReportService.Delete(r.Context(), objectID, userID)
	if err != nil {
		if err == domain.ErrIncidentReportNotFound {
			WriteJSON(w, http.StatusNotFound, Response{Error: "incident report not found"})
			return
		}

		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to delete incident report"})
		return
	}

	WriteJSON(w, http.StatusNoContent, nil)
}

func (h *IncidentReportHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(jwt.MapClaims)
	if !ok {
		WriteJSON(w, http.StatusUnauthorized, Response{Error: "unauthorized"})
		return
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		WriteJSON(w, http.StatusUnauthorized, Response{Error: "invalid user id in token"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, Response{Error: "invalid user id format"})
		return
	}

	filter := domain.NewIncidentReportFilter()

	filter.UserID = userID

	if tripId := r.URL.Query().Get("tripID"); tripId != "" {
		if id, err := primitive.ObjectIDFromHex(tripId); err != nil {
			filter.TripID = &id
		}
	}

	if truckId := r.URL.Query().Get("truckID"); truckId != "" {
		if id, err := primitive.ObjectIDFromHex(truckId); err != nil {
			filter.TruckID = &id
		}
	}

	if truckId := r.URL.Query().Get("truckID"); truckId != "" {
		if id, err := primitive.ObjectIDFromHex(truckId); err != nil {
			filter.TruckID = &id
		}
	}

	filter.Limit = int64(getQueryIntParam(r, "limit", 10))
	filter.Offset = int64(getQueryIntParam(r, "offset", 0))

	result, err := h.incidentReportService.List(r.Context(), filter)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to fetch incident reports"})
		return
	}

	incidentReportResponses := make([]IncidentReportResponse, len(result.IncidentReports))
	for i, d := range result.IncidentReports {
		incidentReportResponses[i] = incidentReportDomainToResponse(d)
	}

	var nextOffset *int64
	if filter.Offset+filter.Limit < result.Total {
		next := filter.Offset + filter.Limit
		nextOffset = &next
	}

	response := PaginatedResponse{
		Items:      incidentReportResponses,
		Total:      result.Total,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
		NextOffset: nextOffset,
	}

	WriteJSON(w, http.StatusOK, response)
}
