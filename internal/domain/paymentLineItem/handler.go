package payment

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
	r.POST("/payment-line-items", h.Create)
	r.GET("/payment-line-items/:uid", h.GetByUID)
	r.PUT("/payment-line-items/:uid", h.Update)
	r.DELETE("/payment-line-items/:uid", h.Delete)
	r.GET("/contractors/:id/payment-line-items", h.GetByContractor)
}

func (h *Handler) Create(c *gin.Context) {
	var req PaymentLineItem
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
		Status string  `json:"status"`
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var resp PaymentLineItem
	var err error
	
	if updateData.Status != "" {
		resp, err = h.svc.UpdateStatus(c.Param("uid"), updateData.Status)
	} else if updateData.Amount > 0 {
		resp, err = h.svc.UpdateAmount(c.Param("uid"), updateData.Amount)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either status or amount must be provided"})
		return
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Delete(c *gin.Context) {
	// For SCD, we don't actually delete but mark as cancelled
	_, err := h.svc.UpdateStatus(c.Param("uid"), "cancelled")
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
