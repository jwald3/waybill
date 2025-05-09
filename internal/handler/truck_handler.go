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

type TruckHandler struct {
	truckService service.TruckService
}

func NewTruckHandler(truckService service.TruckService) *TruckHandler {
	return &TruckHandler{truckService: truckService}
}

var (
	invalidTruckId = "invalid truck id"
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

type TruckUpdateStatusRequest struct {
	Status domain.TruckStatus `json:"status"`
}

type TruckUpdateMileageRequest struct {
	Mileage int `json:"mileage"`
}

type TruckUpdateLastMaintenanceRequest struct {
	LastMaintenance string `json:"last_maintenance"`
}

type TruckResponse struct {
	ID               primitive.ObjectID  `json:"id,omitempty"`
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

func truckRequestToDomainCreate(userID primitive.ObjectID, req TruckCreateRequest) (*domain.Truck, error) {
	return domain.NewTruck(
		userID,
		req.TruckNumber,
		req.VIN,
		req.Make,
		req.Model,
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
	if !req.FuelType.IsValid() {
		return nil, fmt.Errorf("invalid fuel type provided: %s", req.FuelType)
	}

	if !req.TrailerType.IsValid() {
		return nil, fmt.Errorf("invalid trailer type provided: %s", req.TrailerType)
	}

	return &domain.Truck{
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

	var req TruckCreateRequest

	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	truck, err := truckRequestToDomainCreate(userID, req)
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

	truck, err := h.truckService.GetById(r.Context(), objectID, userID)
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
		WriteJSON(w, http.StatusBadRequest, Response{Error: invalidTruckId})
		return
	}

	err = h.truckService.Delete(r.Context(), objectID, userID)
	if err != nil {
		if err == domain.ErrTruckNotFound {
			WriteJSON(w, http.StatusNotFound, Response{Error: "truck not found"})
			return
		}

		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to delete truck"})
		return
	}

	WriteJSON(w, http.StatusNoContent, nil)
}

func (h *TruckHandler) List(w http.ResponseWriter, r *http.Request) {
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

	filter := domain.NewTruckFilter()
	filter.UserID = userID

	if vin := r.URL.Query().Get("vin"); vin != "" {
		filter.VIN = vin
	}

	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = domain.TruckStatus(status)
	}

	if assignedDriver := r.URL.Query().Get("assignedDriver"); assignedDriver != "" {
		if id, err := primitive.ObjectIDFromHex(assignedDriver); err == nil {
			filter.AssignedDriverID = &id
		}
	}

	if trailerType := r.URL.Query().Get("trailerType"); trailerType != "" {
		filter.TrailerType = domain.TrailerType(trailerType)
	}

	if fuelType := r.URL.Query().Get("fuelType"); fuelType != "" {
		filter.FuelType = domain.FuelType(fuelType)
	}

	filter.Limit = int64(getQueryIntParam(r, "limit", 10))
	filter.Offset = int64(getQueryIntParam(r, "offset", 0))

	result, err := h.truckService.List(r.Context(), filter)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to fetch trucks"})
		return
	}

	truckResponses := make([]TruckResponse, len(result.Trucks))
	for i, t := range result.Trucks {
		truckResponses[i] = truckDomainToResponse(t)
	}

	var nextOffset *int64
	if filter.Offset+filter.Limit < result.Total {
		next := filter.Offset + filter.Limit
		nextOffset = &next
	}

	response := PaginatedResponse{
		Items:      truckResponses,
		Total:      result.Total,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
		NextOffset: nextOffset,
	}

	WriteJSON(w, http.StatusOK, response)
}

// atomic methods
func (h *TruckHandler) MakeTruckAvailable(w http.ResponseWriter, r *http.Request) {
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

	if err := h.truckService.MakeTruckAvailable(r.Context(), objectID, userID); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	updatedTruck, err := h.truckService.GetById(r.Context(), objectID, userID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "status updated but failed to fetch updated truck"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: truckDomainToResponse(updatedTruck)})
}

func (h *TruckHandler) RetireTruck(w http.ResponseWriter, r *http.Request) {
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

	if err := h.truckService.RetireTruck(r.Context(), objectID, userID); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	updatedTruck, err := h.truckService.GetById(r.Context(), objectID, userID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "status updated but failed to fetch updated truck"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: truckDomainToResponse(updatedTruck)})
}

func (h *TruckHandler) SetTruckInTransit(w http.ResponseWriter, r *http.Request) {
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

	if err := h.truckService.SetTruckInTransit(r.Context(), objectID, userID); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	updatedTruck, err := h.truckService.GetById(r.Context(), objectID, userID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "status updated but failed to fetch updated truck"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: truckDomainToResponse(updatedTruck)})
}

func (h *TruckHandler) SetTruckInMaintenance(w http.ResponseWriter, r *http.Request) {
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

	if err := h.truckService.SetTruckInMaintenance(r.Context(), objectID, userID); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	updatedTruck, err := h.truckService.GetById(r.Context(), objectID, userID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "status updated but failed to fetch updated truck"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: truckDomainToResponse(updatedTruck)})
}

func (h *TruckHandler) UpdateTruckMileage(w http.ResponseWriter, r *http.Request) {
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

	var req TruckUpdateMileageRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	if err := h.truckService.UpdateTruckMileage(r.Context(), objectID, userID, req.Mileage); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	updatedTruck, err := h.truckService.GetById(r.Context(), objectID, userID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "mileage updated but failed to fetch updated truck"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: truckDomainToResponse(updatedTruck)})
}

func (h *TruckHandler) UpdateTruckLastMaintenance(w http.ResponseWriter, r *http.Request) {
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

	var req TruckUpdateLastMaintenanceRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	if err := h.truckService.UpdateTruckMaintenance(r.Context(), objectID, userID, req.LastMaintenance); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	updatedTruck, err := h.truckService.GetById(r.Context(), objectID, userID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "maintenance updated but failed to fetch updated truck"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: truckDomainToResponse(updatedTruck)})
}
