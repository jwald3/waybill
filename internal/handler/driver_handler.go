package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jwald3/waybill/internal/domain"
	"github.com/jwald3/waybill/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DriverHandler struct {
	driverService service.DriverService
}

func NewDriverHandler(driverService service.DriverService) *DriverHandler {
	return &DriverHandler{driverService: driverService}
}

var (
	invalidDriverId = "invalid driver id"
)

// DTOS =======================================================

type DriverCreateRequest struct {
	FirstName         string         `json:"first_name"`
	LastName          string         `json:"last_name"`
	DOB               string         `json:"dob"`
	LicenseNumber     string         `json:"license_number"`
	LicenseState      string         `json:"license_state"`
	LicenseExpiration string         `json:"license_expiration"`
	Phone             string         `json:"phone"`
	Email             string         `json:"email"`
	Address           domain.Address `json:"address"`
}

type DriverUpdateRequest struct {
	FirstName         string                  `json:"first_name"`
	LastName          string                  `json:"last_name"`
	DOB               string                  `json:"dob"`
	LicenseNumber     string                  `json:"license_number"`
	LicenseState      string                  `json:"license_state"`
	LicenseExpiration string                  `json:"license_expiration"`
	Phone             string                  `json:"phone"`
	Email             string                  `json:"email"`
	Address           domain.Address          `json:"address"`
	EmploymentStatus  domain.EmploymentStatus `json:"employment_status"`
}

type DriverResponse struct {
	ID                primitive.ObjectID      `json:"id,omitempty"`
	FirstName         string                  `json:"first_name"`
	LastName          string                  `json:"last_name"`
	DOB               string                  `json:"dob"`
	LicenseNumber     string                  `json:"license_number"`
	LicenseState      string                  `json:"license_state"`
	LicenseExpiration string                  `json:"license_expiration"`
	Phone             domain.PhoneNumber      `json:"phone"`
	Email             domain.Email            `json:"email"`
	Address           domain.Address          `json:"address"`
	EmploymentStatus  domain.EmploymentStatus `json:"employment_status"`
	CreatedAt         primitive.DateTime      `json:"created_at"`
	UpdatedAt         primitive.DateTime      `json:"updated_at"`
}

type ListDriversResponse struct {
	Drivers []DriverResponse `json:"drivers"`
}

func driverRequestToDomainCreate(req DriverCreateRequest) (*domain.Driver, error) {
	return domain.NewDriver(
		req.FirstName,
		req.LastName,
		req.DOB,
		req.LicenseNumber,
		req.LicenseState,
		req.LicenseExpiration,
		req.Phone,
		req.Email,
		req.Address,
	)
}

func driverRequestToDomainUpdate(req DriverUpdateRequest) (*domain.Driver, error) {
	validEmail, err := domain.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	validPhone, err := domain.NewPhoneNumber(req.Phone)
	if err != nil {
		return nil, err
	}

	return &domain.Driver{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		DOB:               req.DOB,
		LicenseNumber:     req.LicenseNumber,
		LicenseState:      req.LicenseState,
		LicenseExpiration: req.LicenseExpiration,
		Phone:             validPhone,
		Email:             validEmail,
		Address:           req.Address,
		EmploymentStatus:  req.EmploymentStatus,
	}, nil
}

func driverDomainToResponse(d *domain.Driver) DriverResponse {
	return DriverResponse{
		ID:                d.ID,
		FirstName:         d.FirstName,
		LastName:          d.LastName,
		DOB:               d.DOB,
		LicenseNumber:     d.LicenseNumber,
		LicenseState:      d.LicenseState,
		LicenseExpiration: d.LicenseExpiration,
		Phone:             d.Phone,
		Email:             d.Email,
		Address:           d.Address,
		EmploymentStatus:  d.EmploymentStatus,
		CreatedAt:         d.CreatedAt,
		UpdatedAt:         d.UpdatedAt,
	}
}

// =================================================================

func (h *DriverHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req DriverCreateRequest

	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	driver, err := driverRequestToDomainCreate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.driverService.Create(r.Context(), driver); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	WriteJSON(w, http.StatusCreated, driver)
}

func (h *DriverHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	driver, err := h.driverService.GetById(r.Context(), objectID)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, Response{Error: "driver not found"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: driverDomainToResponse(driver)})
}

func (h *DriverHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	var req DriverUpdateRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	driver, err := driverRequestToDomainUpdate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	driver.ID = objectID

	if err := h.driverService.Update(r.Context(), driver); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to update driver"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: driver})
}

func (h *DriverHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: invalidDriverId})
		return
	}

	err = h.driverService.Delete(r.Context(), objectID)
	if err != nil {
		if err == domain.ErrDriverNotFound {
			WriteJSON(w, http.StatusNotFound, Response{Error: "driver not found"})
			return
		}
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to delete driver"})
		return
	}

	WriteJSON(w, http.StatusNoContent, nil)
}

func (h *DriverHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := domain.NewDriverFilter()

	if licenseState := r.URL.Query().Get("licenseState"); licenseState != "" {
		filter.LicenseState = licenseState
	}

	if phone := r.URL.Query().Get("phone"); phone != "" {
		filter.Phone = domain.PhoneNumber(phone)
	}

	if email := r.URL.Query().Get("email"); email != "" {
		filter.Email = domain.Email(email)
	}

	if employmentStatus := r.URL.Query().Get("employmentStatus"); employmentStatus != "" {
		filter.EmploymentStatus = domain.EmploymentStatus(employmentStatus)
	}

	result, err := h.driverService.List(r.Context(), filter)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to fetch drivers"})
		return
	}

	driverResponses := make([]DriverResponse, len(result.Drivers))
	for i, d := range result.Drivers {
		driverResponses[i] = driverDomainToResponse(d)
	}

	var nextOffset *int64
	if filter.Offset+filter.Limit < result.Total {
		next := filter.Offset + filter.Limit
		nextOffset = &next
	}

	response := PaginatedResponse{
		Items:      driverResponses,
		Total:      result.Total,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
		NextOffset: nextOffset,
	}

	WriteJSON(w, http.StatusOK, response)
}

func (h *DriverHandler) SuspendDriver(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.driverService.SuspendDriver(r.Context(), objectID); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	updatedDriver, err := h.driverService.GetById(r.Context(), objectID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "status updated but failed to fetch updated driver"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: driverDomainToResponse(updatedDriver)})
}

func (h *DriverHandler) TerminateDriver(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.driverService.TerminateDriver(r.Context(), objectID); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	updatedDriver, err := h.driverService.GetById(r.Context(), objectID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "status updated but failed to fetch updated driver"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: driverDomainToResponse(updatedDriver)})
}

func (h *DriverHandler) ActivateDriver(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.driverService.ActivateDriver(r.Context(), objectID); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	updatedDriver, err := h.driverService.GetById(r.Context(), objectID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "status updated but failed to fetch updated driver"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: driverDomainToResponse(updatedDriver)})
}
