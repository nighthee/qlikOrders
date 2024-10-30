package order

import (
	"errors"
	"fmt"
	"net/http"
	"qlikOrders/internal/collections"
	"qlikOrders/internal/models"

	"github.com/gin-gonic/gin"
)

// This should be a configuration in a real production environment
// Set intentionally low for testing
const MaxBatchSize = 5

// AddOrdersHandler adds orders in a batch
func AddOrdersHandler(collection *collections.OrderCollection) gin.HandlerFunc {
	return func(c *gin.Context) {

		var newOrders []models.Order

		if err := c.BindJSON(&newOrders); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := validateOrder(newOrders); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if len(newOrders) > MaxBatchSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Batch size exceeds the allowed limit",
				"message": fmt.Sprintf("The maximum allowed number of orders in a single request is %d. Please split your request and try again.", MaxBatchSize),
			})
			return
		}

		if err := collection.AddOrders(newOrders); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add orders"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Orders added successfully"})
	}
}

func validateOrder(orders []models.Order) error {
	for _, order := range orders {
		if order.CustomerID == "" || order.OrderID == "" || order.Timestamp == "" || len(order.Items) == 0 {
			return errors.New("order is missing required fields or has invalid cost")
		}
		for _, item := range order.Items {
			if item.ItemID == "" || item.CostEur <= 0 {
				return errors.New("item is missing required fields or has invalid cost")
			}
		}
	}
	return nil
}
