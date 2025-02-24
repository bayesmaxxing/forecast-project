package handlers

import (
	"backend/internal/auth"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.service.CreateUser(r.Context(), &user)
	respondJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get requested user ID to delete
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Only allow users to delete their own account or require admin role
	if userID != claims.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.service.DeleteUser(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, users)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var changePasswordRequest struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&changePasswordRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.ChangePassword(r.Context(), claims.UserID, changePasswordRequest.OldPassword, changePasswordRequest.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Password changed successfully"})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByUsername(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.VerifyPassword(r.Context(), loginRequest.Username, loginRequest.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"token":    token,
		"username": user.Username,
	})
}
