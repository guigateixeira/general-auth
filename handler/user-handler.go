package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/guigateixeira/general-auth/services"
	"github.com/guigateixeira/general-auth/util"
)

type UserHandler struct {
	service *services.UserService
}

func New(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
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

	// Validate required fields
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Missing required fields: email and password")
		return
	}

	ctx := r.Context()
	userID, err := h.service.CreateUser(ctx, req.Email, req.Password)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	response := map[string]string{"userId": userID}
	util.RespondWithJSON(w, http.StatusCreated, response)
}
