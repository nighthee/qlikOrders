package customer

import (
	"net/http"
	"qlikOrders/internal/collections"

	"github.com/gin-gonic/gin"
)

// GetItemsByCustomerHandler
// Retrieves list of items for a specific customer
func GetItemsByCustomerHandler(collections *collections.OrderCollection) gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID := c.Param("customerId")
		items, err := collections.GetItemsByCustomer(customerID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": items})
	}
}
