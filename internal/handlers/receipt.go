package handlers

import (
	"net/http"
	"receipt-processor/internal/calculator"
	"receipt-processor/internal/models"
	"receipt-processor/internal/store"
	"receipt-processor/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReceiptHandler struct {
	store *store.MemoryStore
}

func NewReceiptHandler(store *store.MemoryStore) *ReceiptHandler {
	return &ReceiptHandler{store: store}
}

func (h *ReceiptHandler) ProcessReceipt(c *gin.Context) {
	var receipt models.Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "The receipt is invalid."})
		return
	}

	// Receipt Validation
	v := validator.New()
	if !v.ValidateReceipt(&receipt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"description": "The receipt is invalid.",
			"details":     v.Errors(),
		})
		return
	}

	// Generating uuid ID
	id := uuid.New().String()

	// Points Calculation
	points := calculator.CalculatePoints(&receipt)

	// Save receipt and points
	h.store.SaveReceipt(id, &receipt, points)

	c.JSON(http.StatusOK, models.ProcessResponse{ID: id})
}

func (h *ReceiptHandler) GetPoints(c *gin.Context) {
	id := c.Param("id")
	validator := validator.New()

	// Receipt ID Validation
	if !validator.ValidateReceiptID(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"description": "No receipt found for that ID.",
			"details":     validator.Errors(),
		})
		return
	}

	points, exists := h.store.GetPoints(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"description": "No receipt found for that ID."})
		return
	}

	c.JSON(http.StatusOK, models.PointsResponse{Points: points})
}
