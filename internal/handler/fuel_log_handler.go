package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jwald3/waybill/internal/domain"
	"github.com/jwald3/waybill/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FuelLogHandler struct {
	fuelLogService service.FuelLogService
}

func NewFuelLogHandler(fuelLogService service.FuelLogService) *FuelLogHandler {
	return &FuelLogHandler{fuelLogService: fuelLogService}
}

var (
	invalidFuelLogId = "invalid fuel log id"
)

// DTOS =======================================================

type FuelLogCreateRequest struct {
	TripID           *primitive.ObjectID `json:"trip_id"`
	Date             string              `json:"date"`
	GallonsPurchased float64             `json:"gallons_purchased"`
	PricePerGallon   float64             `json:"price_per_gallon"`
	TotalCost        float64             `json:"total_cost"`
	Location         string              `json:"location"`
	OdometerReading  int                 `json:"odometer_reading"`
}

type FuelLogUpdateRequest struct {
	TripID           *primitive.ObjectID `json:"trip_id"`
	Date             string              `json:"date"`
	GallonsPurchased float64             `json:"gallons_purchased"`
	PricePerGallon   float64             `json:"price_per_gallon"`
	TotalCost        float64             `json:"total_cost"`
	Location         string              `json:"location"`
	OdometerReading  int                 `json:"odometer_reading"`
}

type FuelLogResponse struct {
	ID               primitive.ObjectID  `json:"id,omitempty"`
	TripID           *primitive.ObjectID `json:"trip_id,omitempty"`
	Trip             *domain.Trip        `json:"trip,omitempty"`
	Date             string              `json:"date"`
	GallonsPurchased float64             `json:"gallons_purchased"`
	PricePerGallon   float64             `json:"price_per_gallon"`
	TotalCost        float64             `json:"total_cost"`
	Location         string              `json:"location"`
	OdometerReading  int                 `json:"odometer_reading"`
	CreatedAt        primitive.DateTime  `json:"created_at"`
	UpdatedAt        primitive.DateTime  `json:"updated_at"`
}

type ListFuelLogsResponse struct {
	FuelLogs []FuelLogResponse `json:"fuel_logs"`
}

func fuelLogRequestToDomainCreate(req FuelLogCreateRequest) (*domain.FuelLog, error) {
	return domain.NewFuelLog(
		req.TripID,
		req.Date,
		req.Location,
		req.GallonsPurchased,
		req.PricePerGallon,
		req.TotalCost,
		req.OdometerReading,
	)
}

func fuelLogRequestToDomainUpdate(req FuelLogUpdateRequest) (*domain.FuelLog, error) {

	return &domain.FuelLog{
		TripID:           req.TripID,
		Date:             req.Date,
		GallonsPurchased: req.GallonsPurchased,
		PricePerGallon:   req.PricePerGallon,
		TotalCost:        req.TotalCost,
		Location:         req.Location,
		OdometerReading:  req.OdometerReading,
	}, nil
}

func fuelLogDomainToResponse(f *domain.FuelLog) FuelLogResponse {
	return FuelLogResponse{
		ID:               f.ID,
		TripID:           f.TripID,
		Trip:             f.Trip,
		Date:             f.Date,
		GallonsPurchased: f.GallonsPurchased,
		PricePerGallon:   f.PricePerGallon,
		TotalCost:        f.TotalCost,
		Location:         f.Location,
		OdometerReading:  f.OdometerReading,
		CreatedAt:        f.CreatedAt,
		UpdatedAt:        f.UpdatedAt,
	}
}

// =================================================================

func (h *FuelLogHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req FuelLogCreateRequest

	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	fuelLog, err := fuelLogRequestToDomainCreate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.fuelLogService.Create(r.Context(), fuelLog); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	WriteJSON(w, http.StatusCreated, fuelLogDomainToResponse(fuelLog))
}

func (h *FuelLogHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	fuelLog, err := h.fuelLogService.GetById(r.Context(), objectID)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, Response{Error: "fuel log not found"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: fuelLogDomainToResponse(fuelLog)})
}

func (h *FuelLogHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	var req FuelLogUpdateRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	fuelLog, err := fuelLogRequestToDomainUpdate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	fuelLog.ID = objectID

	if err := h.fuelLogService.Update(r.Context(), fuelLog); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to update fuel log"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: fuelLogDomainToResponse(fuelLog)})
}

func (h *FuelLogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: invalidFuelLogId})
		return
	}

	err = h.fuelLogService.Delete(r.Context(), objectID)
	if err != nil {
		if err == domain.ErrFuelLogNotFound {
			WriteJSON(w, http.StatusNotFound, Response{Error: "fuel log not found"})
			return
		}

		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to delete fuel log"})
		return
	}

	WriteJSON(w, http.StatusNoContent, nil)
}

func (h *FuelLogHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := getQueryIntParam(r, "limit", 10)
	offset := getQueryIntParam(r, "offset", 0)

	result, err := h.fuelLogService.List(r.Context(), int64(limit), int64(offset))
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to fetch fuel logs"})
		return
	}

	fuelLogResponses := make([]FuelLogResponse, len(result.FuelLogs))
	for i, d := range result.FuelLogs {
		fuelLogResponses[i] = fuelLogDomainToResponse(d)
	}

	var nextOffset *int64
	if int64(offset)+int64(limit) < result.Total {
		next := int64(offset + limit)
		nextOffset = &next
	}

	response := PaginatedResponse{
		Items:      fuelLogResponses,
		Total:      result.Total,
		Limit:      int64(limit),
		Offset:     int64(offset),
		NextOffset: nextOffset,
	}

	WriteJSON(w, http.StatusOK, response)
}
