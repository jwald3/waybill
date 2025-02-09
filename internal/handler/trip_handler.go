package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jwald3/waybill/internal/domain"
	"github.com/jwald3/waybill/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripHandler struct {
	tripService service.TripService
}

func NewTripHandler(tripService service.TripService) *TripHandler {
	return &TripHandler{tripService: tripService}
}

var (
	invalidTripId   = "invalid trip id"
	invalidTripPath = "invalid path"
)

// DTOS =======================================================

type TripCreateRequest struct {
	TripNumber      string              `json:"trip_number"`
	DriverID        *primitive.ObjectID `json:"driver_id"`
	TruckID         *primitive.ObjectID `json:"truck_id"`
	StartFacilityID *primitive.ObjectID `json:"start_facility_id"`
	EndFacilityID   *primitive.ObjectID `json:"end_facility_id"`
	Route           domain.Route        `json:"route"`
	StartTime       primitive.DateTime  `json:"start_time"`
	EndTime         primitive.DateTime  `json:"end_time"`
	Status          domain.TripStatus   `json:"status"`
	Cargo           domain.Cargo        `json:"cargo"`
	FuelUsage       float64             `json:"fuel_usage_gallons"`
	DistanceMiles   int                 `json:"distance_miles"`
}

type TripUpdateRequest struct {
	TripNumber      string              `json:"trip_number"`
	DriverID        *primitive.ObjectID `json:"driver_id"`
	TruckID         *primitive.ObjectID `json:"truck_id"`
	StartFacilityID *primitive.ObjectID `json:"start_facility_id"`
	EndFacilityID   *primitive.ObjectID `json:"end_facility_id"`
	Route           domain.Route        `json:"route"`
	StartTime       primitive.DateTime  `json:"start_time"`
	EndTime         primitive.DateTime  `json:"end_time"`
	Status          domain.TripStatus   `json:"status"`
	Cargo           domain.Cargo        `json:"cargo"`
	FuelUsage       float64             `json:"fuel_usage_gallons"`
	DistanceMiles   int                 `json:"distance_miles"`
}

type TripResponse struct {
	ID              primitive.ObjectID  `json:"id,omitempty"`
	TripNumber      string              `json:"trip_number"`
	DriverID        *primitive.ObjectID `json:"driver_id,omitempty"`
	Driver          *domain.Driver      `json:"driver,omitempty"`
	TruckID         *primitive.ObjectID `json:"truck_id,omitempty"`
	Truck           *domain.Truck       `json:"truck,omitempty"`
	StartFacilityID *primitive.ObjectID `json:"start_facility_id,omitempty"`
	StartFacility   *domain.Facility    `json:"start_facility,omitempty"`
	EndFacilityID   *primitive.ObjectID `json:"end_facility_id,omitempty"`
	EndFacility     *domain.Facility    `json:"end_facility,omitempty"`
	Route           domain.Route        `json:"route"`
	StartTime       primitive.DateTime  `json:"start_time"`
	EndTime         primitive.DateTime  `json:"end_time"`
	Status          domain.TripStatus   `json:"status"`
	Cargo           domain.Cargo        `json:"cargo"`
	FuelUsage       float64             `json:"fuel_usage_gallons"`
	DistanceMiles   int                 `json:"distance_miles"`
	CreatedAt       primitive.DateTime  `json:"created_at"`
	UpdatedAt       primitive.DateTime  `json:"updated_at"`
}

type ListTripsResponse struct {
	Trips []TripResponse `json:"trips"`
}

func tripRequestToDomainCreate(req TripCreateRequest) (*domain.Trip, error) {
	return domain.NewTrip(
		req.TripNumber,
		req.Status,
		req.DriverID,
		req.TruckID,
		req.StartFacilityID,
		req.EndFacilityID,
		req.Route,
		req.StartTime,
		req.EndTime,
		req.Cargo,
		req.FuelUsage,
		req.DistanceMiles,
	)
}

func tripRequestToDomainUpdate(req TripUpdateRequest) (*domain.Trip, error) {
	if !req.Status.IsValid() {
		return nil, fmt.Errorf("invalid status provided: %s", req.Status)
	}

	return &domain.Trip{
		TripNumber:      req.TripNumber,
		DriverID:        req.DriverID,
		TruckID:         req.TruckID,
		StartFacilityID: req.StartFacilityID,
		EndFacilityID:   req.EndFacilityID,
		Route:           req.Route,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		Status:          req.Status,
		Cargo:           req.Cargo,
		FuelUsage:       req.FuelUsage,
		DistanceMiles:   req.DistanceMiles,
	}, nil
}

func tripDomainToResponse(t *domain.Trip) TripResponse {
	return TripResponse{
		ID:              t.ID,
		TripNumber:      t.TripNumber,
		DriverID:        t.DriverID,
		Driver:          t.Driver,
		TruckID:         t.TruckID,
		Truck:           t.Truck,
		StartFacilityID: t.StartFacilityID,
		StartFacility:   t.StartFacility,
		EndFacilityID:   t.EndFacilityID,
		EndFacility:     t.EndFacility,
		Route:           t.Route,
		StartTime:       t.StartTime,
		EndTime:         t.EndTime,
		Status:          t.Status,
		Cargo:           t.Cargo,
		FuelUsage:       t.FuelUsage,
		DistanceMiles:   t.DistanceMiles,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}

// =================================================================

func (h *TripHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req TripCreateRequest

	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	trip, err := tripRequestToDomainCreate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.tripService.Create(r.Context(), trip); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	WriteJSON(w, http.StatusCreated, tripDomainToResponse(trip))
}

func (h *TripHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	trip, err := h.tripService.GetById(r.Context(), objectID)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, Response{Error: "trip not found"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: tripDomainToResponse(trip)})
}

func (h *TripHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	var req TripUpdateRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	trip, err := tripRequestToDomainUpdate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	trip.ID = objectID

	if err := h.tripService.Update(r.Context(), trip); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to update user"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: tripDomainToResponse(trip)})
}

func (h *TripHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	err = h.tripService.Delete(r.Context(), objectID)
	if err != nil {
		if err == domain.ErrTripNotFound {
			WriteJSON(w, http.StatusNotFound, Response{Error: "trip not found"})
			return
		}
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to delete trip"})
		return
	}

	WriteJSON(w, http.StatusNoContent, nil)
}

func (h *TripHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := getQueryIntParam(r, "limit", 10)
	offset := getQueryIntParam(r, "offset", 0)

	result, err := h.tripService.List(r.Context(), int64(limit), int64(offset))
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to fetch trips"})
		return
	}

	tripResponses := make([]TripResponse, len(result.Trips))
	for i, t := range result.Trips {
		tripResponses[i] = tripDomainToResponse(t)
	}

	var nextOffset *int64
	if int64(offset)+int64(limit) < result.Total {
		next := int64(offset + limit)
		nextOffset = &next
	}

	response := PaginatedResponse{
		Items:      tripResponses,
		Total:      result.Total,
		Limit:      int64(limit),
		Offset:     int64(offset),
		NextOffset: nextOffset,
	}

	WriteJSON(w, http.StatusOK, response)
}
