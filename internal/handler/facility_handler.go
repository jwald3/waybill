package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	invalidFacilityId = "invalid facility id"
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
	FacilityNumber    string                   `json:"facility_number"`
	Name              string                   `json:"name"`
	Type              string                   `json:"type"`
	Address           domain.Address           `json:"address"`
	ContactInfo       domain.ContactInfo       `json:"contact_info"`
	ParkingCapacity   int                      `json:"parking_capacity"`
	ServicesAvailable []domain.FacilityService `json:"services_available"`
}

type FacilityUpdateAvailableServicesRequest struct {
	AvailableServices []domain.FacilityService `json:"services_available"`
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
	for _, service := range req.ServicesAvailable {
		if !service.IsValid() {
			return nil, fmt.Errorf("invalid facility service: %s", service)
		}
	}

	return &domain.Facility{
		FacilityNumber:    req.FacilityNumber,
		Name:              req.Name,
		Type:              req.Type,
		Address:           req.Address,
		ContactInfo:       req.ContactInfo,
		ParkingCapacity:   req.ParkingCapacity,
		ServicesAvailable: req.ServicesAvailable,
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
		WriteJSON(w, http.StatusBadRequest, Response{Error: invalidFacilityId})
		return
	}

	err = h.facilityService.Delete(r.Context(), objectID)
	if err != nil {
		if err == domain.ErrFacilityNotFound {
			WriteJSON(w, http.StatusNotFound, Response{Error: "facility not found"})
			return
		}

		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to delete facility"})
		return
	}

	WriteJSON(w, http.StatusNoContent, nil)
}

func (h *FacilityHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := domain.NewFacilityFilter()

	// Parse query parameters, adding them to the filter if they're present
	// any unrecognized query params will be ignored and not added to the filter
	// we can expand the filter options as needed by adding more to the domain.FacilityFilter struct
	if stateCode := r.URL.Query().Get("stateCode"); stateCode != "" {
		filter.StateCode = stateCode
	}

	if facilityType := r.URL.Query().Get("type"); facilityType != "" {
		filter.Type = facilityType
	}

	if services := r.URL.Query().Get("services"); services != "" {
		// the API expects a comma-separated list of services, so we split it into a list and check each one
		// the parameter could look like this: ?services=REPAIRS,LOADING_UNLOADING,LODGING
		servicesList := strings.Split(services, ",")
		for _, s := range servicesList {
			service := domain.FacilityService(strings.TrimSpace(s))
			if service.IsValid() {
				filter.ServicesInclude = append(filter.ServicesInclude, service)
			}
		}
	}

	if minCapStr := r.URL.Query().Get("minCapacity"); minCapStr != "" {
		if minCap, err := strconv.Atoi(minCapStr); err == nil {
			filter.MinCapacity = &minCap
		}
	}

	if maxCapStr := r.URL.Query().Get("maxCapacity"); maxCapStr != "" {
		if maxCap, err := strconv.Atoi(maxCapStr); err == nil {
			filter.MaxCapacity = &maxCap
		}
	}

	filter.Limit = int64(getQueryIntParam(r, "limit", 10))
	filter.Offset = int64(getQueryIntParam(r, "offset", 0))

	result, err := h.facilityService.ListWithFilter(r.Context(), filter)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to fetch facilities"})
		return
	}

	facilityResponses := make([]FacilityResponse, len(result.Facilities))
	for i, d := range result.Facilities {
		facilityResponses[i] = facilityDomainToResponse(d)
	}

	var nextOffset *int64
	if filter.Offset+filter.Limit < result.Total {
		next := filter.Offset + filter.Limit
		nextOffset = &next
	}

	response := PaginatedResponse{
		Items:      facilityResponses,
		Total:      result.Total,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
		NextOffset: nextOffset,
	}

	WriteJSON(w, http.StatusOK, response)
}

func (h *FacilityHandler) UpdateAvailableFacilityServices(w http.ResponseWriter, r *http.Request) {
	// Get facility ID from URL
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid facility ID"})
		return
	}

	// Parse request body
	var req FacilityUpdateAvailableServicesRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	// Validate services before updating
	for _, service := range req.AvailableServices {
		if !service.IsValid() {
			WriteJSON(w, http.StatusBadRequest, Response{Error: fmt.Sprintf("invalid facility service: %s", service)})
			return
		}
	}

	// Update services
	err = h.facilityService.UpdateAvailableFacilityServices(r.Context(), objectID, req.AvailableServices)
	if err != nil {
		switch {
		case err == domain.ErrFacilityNotFound:
			WriteJSON(w, http.StatusNotFound, Response{Error: "facility not found"})
		default:
			WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		}
		return
	}

	// Fetch updated facility to return in response
	updatedFacility, err := h.facilityService.GetById(r.Context(), objectID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "services updated but failed to fetch updated facility"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: facilityDomainToResponse(updatedFacility)})
}
