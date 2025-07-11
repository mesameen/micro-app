package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mesameen/micro-app/rating/internal/controller"
	"github.com/mesameen/micro-app/rating/pkg/model"
	"github.com/mesameen/micro-app/src/pkg/logger"
)

// Handler defines a rating service controller
type Handler struct {
	ctrl *controller.Controller
}

// New creates a new rating service HTTP handler
func New(ctrl *controller.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

// Handle handles GET /rating requests
func (h *Handler) GetRatings(c *gin.Context) {
	recordID := model.RecordID(c.Query("id"))
	recordType := model.RecordType(c.Query("type"))
	val, err := h.ctrl.GetAggregatedRating(c.Request.Context(), recordID, recordType)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, val)
}

// SaveRatings handles PUT /rating requests
func (h *Handler) SaveRatings(c *gin.Context) {
	recordID := model.RecordID(c.Query("id"))
	recordType := model.RecordType(c.Query("type"))
	userID := model.UserID(c.Query("userId"))
	v, err := strconv.ParseFloat(c.Query("value"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value for ratino"})
		return
	}
	if err := h.ctrl.PutRating(c.Request.Context(), recordID, recordType, &model.Rating{
		UserID: userID,
		Value:  model.RatingValue(v),
	}); err != nil {
		logger.Errorf("Failed to put rating. Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully stored"})
}
