package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jwald3/waybill/internal/domain"
	"github.com/jwald3/waybill/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TruckHandler struct {
	truckService service.TruckService
}

func NewTruckHandler(truckService service.TruckService) *TruckHandler {
	return &TruckHandler{truckService: truckService}
}

var (
	invalidTruckId   = "invalid truck id"
	invalidTruckPath = "invalid path"
)

// DTOS =======================================================

type TruckCreateRequest struct {
	TruckNumber      string              `json:"truck_number"`
	VIN              string              `json:"vin"`
	Make             string              `json:"make"`
	Model            string              `json:"model"`
	Year             int                 `json:"year"`
	LicensePlate     domain.LicensePlate `json:"license_plate"`
	Mileage          int                 `json:"mileage"`
	Status           domain.TruckStatus  `json:"status"`
	AssignedDriverID *primitive.ObjectID `json:"assigned_driver_id,omitempty"`
	TrailerType      domain.TrailerType  `json:"trailer_type"`
	CapacityTons     float64             `json:"capacity_tons"`
	FuelType         domain.FuelType     `json:"fuel_type"`
	LastMaintenance  string              `json:"last_maintenance"`
}

type TruckUpdateRequest struct {
	ID               primitive.ObjectID  `json:"id"`
	TruckNumber      string              `json:"truck_number"`
	VIN              string              `json:"vin"`
	Make             string              `json:"make"`
	Model            string              `json:"model"`
	Year             int                 `json:"year"`
	LicensePlate     domain.LicensePlate `json:"license_plate"`
	Mileage          int                 `json:"mileage"`
	Status           domain.TruckStatus  `json:"status"`
	AssignedDriverID *primitive.ObjectID `json:"assigned_driver_id,omitempty"`
	TrailerType      domain.TrailerType  `json:"trailer_type"`
	CapacityTons     float64             `json:"capacity_tons"`
	FuelType         domain.FuelType     `json:"fuel_type"`
	LastMaintenance  string              `json:"last_maintenance"`
}

type TruckResponse struct {
	ID               primitive.ObjectID  `json:"_id,omitempty"`
	TruckNumber      string              `json:"truck_number"`
	VIN              string              `json:"vin"`
	Make             string              `json:"make"`
	Model            string              `json:"model"`
	Year             int                 `json:"year"`
	LicensePlate     domain.LicensePlate `json:"license_plate"`
	Mileage          int                 `json:"mileage"`
	Status           domain.TruckStatus  `json:"status"`
	AssignedDriverID *primitive.ObjectID `json:"assigned_driver_id,omitempty"`
	AssignedDriver   *domain.Driver      `json:"assigned_driver,omitempty"`
	TrailerType      domain.TrailerType  `json:"trailer_type"`
	CapacityTons     float64             `json:"capacity_tons"`
	FuelType         domain.FuelType     `json:"fuel_type"`
	LastMaintenance  string              `json:"last_maintenance"`
	CreatedAt        primitive.DateTime  `json:"created_at"`
	UpdatedAt        primitive.DateTime  `json:"updated_at"`
}

type ListTrucksResponse struct {
	Trucks []TruckResponse `json:"trucks"`
}

func truckRequestToDomainCreate(req TruckCreateRequest) (*domain.Truck, error) {
	return domain.NewTruck(
		req.TruckNumber,
		req.VIN,
		req.Make,
		req.Model,
		req.Status,
		req.TrailerType,
		req.FuelType,
		req.LastMaintenance,
		req.Year,
		req.Mileage,
		req.CapacityTons,
		req.LicensePlate,
	)
}

func truckRequestToDomainUpdate(req TruckUpdateRequest) (*domain.Truck, error) {
	now := time.Now()

	return &domain.Truck{
		ID:               req.ID,
		TruckNumber:      req.TruckNumber,
		VIN:              req.VIN,
		Make:             req.Make,
		Model:            req.Model,
		Year:             req.Year,
		LicensePlate:     req.LicensePlate,
		Mileage:          req.Mileage,
		Status:           req.Status,
		AssignedDriverID: req.AssignedDriverID,
		TrailerType:      req.TrailerType,
		CapacityTons:     req.CapacityTons,
		FuelType:         req.FuelType,
		LastMaintenance:  req.LastMaintenance,
		UpdatedAt:        primitive.NewDateTimeFromTime(now),
	}, nil
}

func truckDomainToResponse(t *domain.Truck) TruckResponse {
	return TruckResponse{
		ID:               t.ID,
		TruckNumber:      t.TruckNumber,
		VIN:              t.VIN,
		Make:             t.Make,
		Model:            t.Model,
		Year:             t.Year,
		LicensePlate:     t.LicensePlate,
		Mileage:          t.Mileage,
		Status:           t.Status,
		AssignedDriverID: t.AssignedDriverID,
		AssignedDriver:   t.AssignedDriver,
		TrailerType:      t.TrailerType,
		CapacityTons:     t.CapacityTons,
		FuelType:         t.FuelType,
		LastMaintenance:  t.LastMaintenance,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}
}

// =================================================================
func (h *TruckHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req TruckCreateRequest

	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	truck, err := truckRequestToDomainCreate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.truckService.Create(r.Context(), truck); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	WriteJSON(w, http.StatusCreated, truckDomainToResponse(truck))
}

func (h *TruckHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	truck, err := h.truckService.GetById(r.Context(), objectID)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, Response{Error: "truck not found"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: truckDomainToResponse(truck)})
}

func (h *TruckHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	var req TruckUpdateRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	truck, err := truckRequestToDomainUpdate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	truck.ID = objectID

	if err := h.truckService.Update(r.Context(), truck); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to update truck"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: truckDomainToResponse(truck)})
}

func (h *TruckHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.truckService.Delete(r.Context(), objectID); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to delete truck"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Message: "truck deleted successfully"})
}

func (h *TruckHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := getQueryIntParam(r, "limit", 10)
	offset := getQueryIntParam(r, "offset", 0)

	trucks, err := h.truckService.List(r.Context(), int64(limit), int64(offset))
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to fetch trucks"})
		return
	}

	truckResponses := make([]TruckResponse, len(trucks))
	for i, t := range trucks {
		truckResponses[i] = truckDomainToResponse(t)
	}

	WriteJSON(w, http.StatusOK, Response{Data: ListTrucksResponse{Trucks: truckResponses}})
}
