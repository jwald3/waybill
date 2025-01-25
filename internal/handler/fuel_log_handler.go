package handler

import (
	"net/http"
	"time"

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
	invalidFuelLogId   = "invalid fuel log id"
	invalidFuelLogPath = "invalid path"
)

// DTOS =======================================================

type FuelLogCreateRequest struct {
	TruckID          primitive.ObjectID `json:"truck_id"`
	DriverID         primitive.ObjectID `json:"driver_id"`
	Date             string             `json:"date"`
	GallonsPurchased float64            `json:"gallons_purchased"`
	PricePerGallon   float64            `json:"price_per_gallon"`
	TotalCost        float64            `json:"total_cost"`
	Location         string             `json:"location"`
	OdometerReading  int                `json:"odometer_reading"`
}

type FuelLogUpdateRequest struct {
	ID               primitive.ObjectID `json:"_id,omitempty"`
	TruckID          primitive.ObjectID `json:"truck_id"`
	DriverID         primitive.ObjectID `json:"driver_id"`
	Date             string             `json:"date"`
	GallonsPurchased float64            `json:"gallons_purchased"`
	PricePerGallon   float64            `json:"price_per_gallon"`
	TotalCost        float64            `json:"total_cost"`
	Location         string             `json:"location"`
	OdometerReading  int                `json:"odometer_reading"`
}

type FuelLogResponse struct {
	ID               primitive.ObjectID `json:"_id,omitempty"`
	TruckID          primitive.ObjectID `json:"truck_id"`
	DriverID         primitive.ObjectID `json:"driver_id"`
	Date             string             `json:"date"`
	GallonsPurchased float64            `json:"gallons_purchased"`
	PricePerGallon   float64            `json:"price_per_gallon"`
	TotalCost        float64            `json:"total_cost"`
	Location         string             `json:"location"`
	OdometerReading  int                `json:"odometer_reading"`
	CreatedAt        primitive.DateTime `json:"created_at"`
	UpdatedAt        primitive.DateTime `json:"updated_at"`
}

type ListFuelLogsResponse struct {
	FuelLogs []FuelLogResponse `json:"fuel_logs"`
}

func fuelLogRequestToDomainCreate(req FuelLogCreateRequest) (*domain.FuelLog, error) {
	return domain.NewFuelLog(
		req.TruckID,
		req.DriverID,
		req.Date,
		req.Location,
		req.GallonsPurchased,
		req.PricePerGallon,
		req.TotalCost,
		req.OdometerReading,
	)
}

func fuelLogRequestToDomainUpdate(req FuelLogUpdateRequest) (*domain.FuelLog, error) {
	now := time.Now()

	return &domain.FuelLog{
		ID:               req.ID,
		TruckID:          req.TruckID,
		DriverID:         req.DriverID,
		Date:             req.Date,
		GallonsPurchased: req.GallonsPurchased,
		PricePerGallon:   req.PricePerGallon,
		TotalCost:        req.TotalCost,
		Location:         req.Location,
		OdometerReading:  req.OdometerReading,
		UpdatedAt:        primitive.NewDateTimeFromTime(now),
	}, nil
}

func fuelLogDomainToResponse(f *domain.FuelLog) FuelLogResponse {
	return FuelLogResponse{
		TruckID:          f.TruckID,
		DriverID:         f.DriverID,
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
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to update user"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: fuelLogDomainToResponse(fuelLog)})
}

func (h *FuelLogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.fuelLogService.Delete(r.Context(), objectID); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to delete fuel log"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Message: "fuel log deleted successfully"})
}

func (h *FuelLogHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := getQueryIntParam(r, "limit", 10)
	offset := getQueryIntParam(r, "offset", 0)

	fuelLogs, err := h.fuelLogService.List(r.Context(), int64(limit), int64(offset))
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to fetch fuel logs"})
		return
	}

	fuelLogResponses := make([]FuelLogResponse, len(fuelLogs))
	for i, d := range fuelLogs {
		fuelLogResponses[i] = fuelLogDomainToResponse(d)
	}

	WriteJSON(w, http.StatusOK, Response{Data: ListFuelLogsResponse{FuelLogs: fuelLogResponses}})
}
