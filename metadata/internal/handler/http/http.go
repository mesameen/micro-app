package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mesameen/micro-app/metadata/internal/controller"
	"github.com/mesameen/micro-app/metadata/internal/repository"
)

// Handler defines a movie metadata HTTP handler
type Handler struct {
	ctrl *controller.Controller
}

// New creates a new movie metadata HTTP handler
func New(ctrl *controller.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

// GetMetadata handles GET /metadata requests
func (h *Handler) GetMetadata(c *gin.Context) {
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
