package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jwald3/go_rest_template/internal/domain"
	"github.com/jwald3/go_rest_template/internal/service"
)

// DTOs are a useful way to specify which data a user will be privy to during a request/response cycle
// DTOs can be helpful for abstracting away details revealed to a user, such as passwords or internal use only fields
type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   string `json:"status"`
}

type UserResponse struct {
	ID     int64  `json:"id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

// a convenient packaging for the methods that return multiple user objects under the `users` property
type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

// turn the request object into a `User` object - this handles the logic used when creating a user for the first time
func requestToDomainCreate(req UserRequest) (*domain.User, error) {
	return domain.NewUser(req.Email, req.Password)
}

// turn the request object into a `User` object - this ensures that you have control over what changes you make instead
// of using the default constructor
func requestToDomainUpdate(req UserRequest) (*domain.User, error) {
	pass, err := domain.NewPassword(req.Password)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		Email:    req.Email,
		Password: pass,
		Status:   domain.Status(req.Status),
	}, nil
}

// repackage the domain object into a response object
func domainToResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:     u.ID,
		Email:  u.Email,
		Status: string(u.Status),
	}
}

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// you can define commonly used error messages here
var (
	invalidId   = "invalid user ID"
	invalidPath = "invalid path"
)

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	// ensure that the request fits the DTO format before going any further
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	// by using the `requestToDomain` function, you're leveraging the constructor that hashes the
	// password and assigns the default status. If this fails, throw an error
	user, err := requestToDomainCreate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
	}

	// attempt creating the user. If this fails, throw an error
	if err := h.userService.Create(r.Context(), user); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	// repackage the user object into the DTO
	WriteJSON(w, http.StatusCreated, domainToResponse(user))
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	// use gorilla mux to attempt to extract out the ID from the URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	user, err := h.userService.Get(r.Context(), int64(id))
	if err != nil {
		WriteJSON(w, http.StatusNotFound, Response{Error: "user not found"})
		return
	}

	// repackage the user object into a DTO
	WriteJSON(w, http.StatusOK, Response{Data: domainToResponse(user)})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	var req UserRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	// use the update-specific function to ensure you aren't overwriting any of the properties that get
	// set when using the default constructor
	user, err := requestToDomainUpdate(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	user.ID = int64(id)

	if err := h.userService.Update(r.Context(), user); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to update user"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Data: domainToResponse(user)})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := h.userService.Delete(r.Context(), int64(id)); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to delete user"})
		return
	}

	WriteJSON(w, http.StatusOK, Response{Message: "user deleted successfully"})
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := getQueryIntParam(r, "limit", 10)
	offset := getQueryIntParam(r, "offset", 0)

	users, err := h.userService.List(r.Context(), limit, offset)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: "failed to fetch users"})
		return
	}

	// use the packaged response type to send a list of users
	userResponses := make([]UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = domainToResponse(u)
	}

	WriteJSON(w, http.StatusOK, Response{Data: ListUsersResponse{Users: userResponses}})
}
