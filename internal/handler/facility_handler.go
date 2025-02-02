package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jwald3/waybill/internal/domain"
	"github.com/jwald3/waybill/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FacilityHandler struct {
	facilityService service.FacilityService
}

func NewFacilityHandler(facilityService service.FacilityService) *FacilityHandler {
	return &FacilityHandler{facilityService: facilityService}
}

var (
	invalidFacilityId   = "invalid facility id"
	invalidFacilityPath = "invalid path"
)

// DTOS =======================================================

type FacilityCreateRequest struct {
	FacilityNumber    string                   `json:"facility_number"`
	Name              string                   `json:"name"`
	Type              string                   `json:"type"`
	Address           domain.Address           `json:"address"`
	ContactInfo       domain.ContactInfo       `json:"contact_info"`
	ParkingCapacity   int                      `json:"parking_capacity"`
	ServicesAvailable []domain.FacilityService `json:"services_available"`
}

type FacilityUpdateRequest struct {
	ID                primitive.ObjectID       `json:"id,omitempty"`
	FacilityNumber    string                   `json:"facility_number"`
	Name              string                   `json:"name"`
	Type              string                   `json:"type"`
	Address           domain.Address           `json:"address"`
	ContactInfo       domain.ContactInfo       `json:"contact_info"`
	ParkingCapacity   int                      `json:"parking_capacity"`
	ServicesAvailable []domain.FacilityService `json:"services_available"`
}

type FacilityResponse struct {
	ID                primitive.ObjectID       `json:"id,omitempty"`
	FacilityNumber    string                   `json:"facility_number"`
	Name              string                   `json:"name"`
	Type              string                   `json:"type"`
	Address           domain.Address           `json:"address"`
	ContactInfo       domain.ContactInfo       `json:"contact_info"`
	ParkingCapacity   int                      `json:"parking_capacity"`
	ServicesAvailable []domain.FacilityService `json:"services_available"`
	CreatedAt         primitive.DateTime       `json:"created_at"`
	UpdatedAt         primitive.DateTime       `json:"updated_at"`
}

type ListFacilitiesResponse struct {
	Facilities []FacilityResponse `json:"facilities"`
}

func facilityRequestToDomainCreate(req FacilityCreateRequest) (*domain.Facility, error) {
	return domain.NewFacility(
		req.FacilityNumber,
		req.Name,
		req.Type,
		req.Address,
		req.ContactInfo,
		req.ParkingCapacity,
		req.ServicesAvailable,
	)
}

func facilityRequestToDomainUpdate(req FacilityUpdateRequest) (*domain.Facility, error) {
	now := time.Now()

	return &domain.Facility{
		ID:                req.ID,
		FacilityNumber:    req.FacilityNumber,
		Name:              req.Name,
		Type:              req.Type,
		Address:           req.Address,
		ContactInfo:       req.ContactInfo,
		ParkingCapacity:   req.ParkingCapacity,
		ServicesAvailable: req.ServicesAvailable,
		UpdatedAt:         primitive.NewDateTimeFromTime(now),
	}, nil
}

func facilityDomainToResponse(f *domain.Facility) FacilityResponse {
	return FacilityResponse{
		ID:                f.ID,
		FacilityNumber:    f.FacilityNumber,
		Name:              f.Name,
		Type:              f.Type,
		Address:           f.Address,
		ContactInfo:       f.ContactInfo,
		ParkingCapacity:   f.ParkingCapacity,
		ServicesAvailable: f.ServicesAvailable,
		CreatedAt:         f.CreatedAt,
		UpdatedAt:         f.UpdatedAt,
	}
}

// =================================================================

func (h *FacilityHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req FacilityCreateRequest

	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	facility, err := facilityRequestToDomainCreate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.facilityService.Create(r.Context(), facility); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	WriteJSON(w, http.StatusCreated, facilityDomainToResponse(facility))
}

func (h *FacilityHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	facility, err := h.facilityService.GetById(r.Context(), objectID)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, Response{Error: "facility not found"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: facilityDomainToResponse(facility)})
}

func (h *FacilityHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	var req FacilityUpdateRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	facility, err := facilityRequestToDomainUpdate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	facility.ID = objectID

	if err := h.facilityService.Update(r.Context(), facility); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to update user"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: facilityDomainToResponse(facility)})
}

func (h *FacilityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.facilityService.Delete(r.Context(), objectID); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to delete facility"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Message: "facility deleted successfully"})
}

func (h *FacilityHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := getQueryIntParam(r, "limit", 10)
	offset := getQueryIntParam(r, "offset", 0)

	facilities, err := h.facilityService.List(r.Context(), int64(limit), int64(offset))
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to fetch facilities"})
		return
	}

	facilityResponses := make([]FacilityResponse, len(facilities))
	for i, d := range facilities {
		facilityResponses[i] = facilityDomainToResponse(d)
	}

	WriteJSON(w, http.StatusOK, Response{Data: ListFacilitiesResponse{Facilities: facilityResponses}})
}
