package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/guigateixeira/general-auth/errors"
	"github.com/guigateixeira/general-auth/services"
	"github.com/guigateixeira/general-auth/util"
)

type UserHandler struct {
	userService *services.UserService
}

func New(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Missing required fields: email and password")
		return
	}

	ctx := r.Context()
	userID, err := h.userService.CreateUser(ctx, req.Email, req.Password)
	if err != nil {
		serviceErr, ok := err.(*errors.BaseError)
		if !ok {
			serviceErr = errors.NewBaseError("Unknown error occurred", http.StatusInternalServerError)
		}
		util.RespondWithError(w, serviceErr.StatusCode, serviceErr.Message)
		return
	}

	response := map[string]string{"userId": userID}
	util.RespondWithJSON(w, http.StatusCreated, response)
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Email == "" || req.Password == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	ctx := r.Context()
	token, err := h.userService.SignIn(ctx, req.Email, req.Password)
	if err != nil {
		serviceErr, ok := err.(*errors.BaseError)
		if !ok {
			serviceErr = errors.NewBaseError("Unknown error occurred", http.StatusInternalServerError)
		}
		util.RespondWithError(w, serviceErr.StatusCode, serviceErr.Message)
		return
	}

	response := map[string]string{"token": token}
	util.RespondWithJSON(w, http.StatusOK, response)
}
