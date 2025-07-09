package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mesameen/micro-app/movie/internal/controller/movie"
)

// Handler defines a movie handler
type Handler struct {
	ctrl *movie.Controller
}

// New creates a new movie HTTP handler
func New(ctrl *movie.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

// GetMovieDetails handles GET /movie requests
func (h *Handler) GetMovieDetails(c *gin.Context) {
	id := c.Query("id")
	details, err := h.ctrl.Get(c.Request.Context(), id)
	if err != nil && errors.Is(err, movie.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, details)
}
