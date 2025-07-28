package timelog

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(s Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/timelogs", h.Create)
	r.GET("/timelogs/:uid", h.GetByUID)
	r.PUT("/timelogs/:uid", h.Update)
	r.DELETE("/timelogs/:uid", h.Delete)
	r.GET("/contractors/:id/timelogs", h.GetByContractor)
}

func (h *Handler) Create(c *gin.Context) {
	var req Timelog
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.svc.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetByUID(c *gin.Context) {
	resp, err := h.svc.GetByUID(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Update(c *gin.Context) {
	var updateData struct {
		Duration int64  `json:"duration"`
		Type     string `json:"type"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var resp Timelog
	var err error

	if updateData.Duration > 0 {
		resp, err = h.svc.UpdateDuration(c.Param("uid"), updateData.Duration)
	} else if updateData.Type != "" {
		resp, err = h.svc.UpdateType(c.Param("uid"), updateData.Type)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either duration or type must be provided"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Delete(c *gin.Context) {
	// For SCD, we don't actually delete but mark as invalid by setting type
	_, err := h.svc.UpdateType(c.Param("uid"), "deleted")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetByContractor(c *gin.Context) {
	resp, err := h.svc.GetByContractor(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
