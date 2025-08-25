package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mesameen/micro-app/movie/internal/controller"
	"github.com/mesameen/micro-app/src/pkg/telemetryservice"
)

// Handler defines a movie handler
type Handler struct {
	telem telemetryservice.Repo
	ctrl  *controller.Controller
}

// New creates a new movie HTTP handler
func New(telem telemetryservice.Repo, ctrl *controller.Controller) *Handler {
	return &Handler{
		telem: telem,
		ctrl:  ctrl,
	}
}

// GetMovieDetails handles GET /movie requests
func (h *Handler) GetMovieDetails(c *gin.Context) {
	ctx, span := h.telem.TraceStart(c.Request.Context(), "get_movie_details")
	defer span.End()
	id := c.Query("id")
	details, err := h.ctrl.Get(c.Request.Context(), id)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		h.telem.Errorf(ctx, "Requested movie: %v not found", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	} else if err != nil {
		h.telem.Errorf(ctx, "Failed to get movie details. Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, details)
}
