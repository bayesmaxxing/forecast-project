package handlers

import (
	"encoding/json"
	"backend/internal/models"
	"backend/internal/services"
	"net/http"
)

type BlogpostHandler struct {
	service *services.BlogpostService
}

func NewBlogpostHandler(s *services.BlogpostService) *BlogpostHandler {
	return &BlogpostHandler{service: s}
}

func (h *BlogpostHandler) ListBlogposts(w http.ResponseWriter, r *http.Request) {
	blogposts, err := h.service.GetBlogposts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, blogposts)
}

func (h *BlogpostHandler) GetBlogpostBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")

	blogpost, err := h.service.GetBlogpostBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, blogpost)
}

func (h *BlogpostHandler) CreateBlogpost(w http.ResponseWriter, r *http.Request) {
	var blogpost models.Blogpost
	if err := json.NewDecoder(r.Body).Decode(&blogpost); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.CreateBlogpost(r.Context(), &blogpost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	message := "blogpost created"
	respondJSON(w, http.StatusCreated, message)
}
