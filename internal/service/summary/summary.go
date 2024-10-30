package summary

import (
	"net/http"
	"qlikOrders/internal/collections"

	"github.com/gin-gonic/gin"
)

// GetSummariesHandler
// Retrieves a summary total spend and number of items for all customers
func GetSummariesHandler(collections *collections.OrderCollection) gin.HandlerFunc {
	return func(c *gin.Context) {

		summaries, err := collections.GetAllCustomerSummaries()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve summaries"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"summaries": summaries})
	}
}
