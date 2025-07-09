package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mesameen/micro-app/metadata/internal/controller/metadata"
	"github.com/mesameen/micro-app/metadata/internal/repository"
)

// Handler defines a movie metadata HTTP handler
type Handler struct {
	ctrl *metadata.Controller
}

// New creates a new movie metadata HTTP handler
func New(ctrl *metadata.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

// GetMetadata handles GET /metadata requests
func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	m, err := h.ctrl.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error for mobie:%s: %v\n", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetMetadata handles GET /metadata requests
func (h *Handler) GetMetadata1(c *gin.Context) {
	fmt.Println("isnide GetMetadata1")
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id isn't presents"})
		return
	}
	// ctx := r.Context()
	m, err := h.ctrl.Get(c.Request.Context(), id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": repository.ErrNotFound})
		return
	} else if err != nil {
		log.Printf("Repository get error for mobie:%s: %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, m)
}
