package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jwald3/go_rest_template/internal/domain"
	"github.com/jwald3/go_rest_template/internal/service"
)

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

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

func requestToDomain(req UserRequest) (*domain.User, error) {
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

var (
	invalidId   = "invalid user ID"
	invalidPath = "invalid path"
)

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: "invalid request payload"})
		return
	}

	user, err := requestToDomain(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
	}

	if err := h.userService.Create(r.Context(), user); err != nil {
		WriteJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	WriteJSON(w, http.StatusCreated, domainToResponse(user))
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
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

	user, err := requestToDomain(req)
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

	userResponses := make([]UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = domainToResponse(u)
	}

	WriteJSON(w, http.StatusOK, Response{Data: ListUsersResponse{Users: userResponses}})
}
