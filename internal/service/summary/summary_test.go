package summary

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"qlikOrders/internal/collections"
	"qlikOrders/internal/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetSummariesHandler(t *testing.T) {
	// Set up the Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Create a test collection with sample summaries
	testCollection := &collections.OrderCollection{
		Orders: []models.Order{
			{
				CustomerID: "01",
				OrderID:    "50",
				Timestamp:  "1637245070513",
				Items: []models.Item{
					{ItemID: "20201", CostEur: 2},
					{ItemID: "20202", CostEur: 3},
				},
			},
			{
				CustomerID: "02",
				OrderID:    "51",
				Timestamp:  "1637245070514",
				Items: []models.Item{
					{ItemID: "20203", CostEur: 5},
				},
			},
		},
	}

	// Define the route for the test
	router.GET("/summary", GetSummariesHandler(testCollection))

	// Test case for a successful summary retrieval
	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/summary", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Check the response status
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify the response structure
		var response struct {
			Summaries []models.Summary `json:"summaries"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Summaries))

		// Validate the contents of the summaries
		assert.Equal(t, "01", response.Summaries[0].CustomerID)
		assert.Equal(t, 2, response.Summaries[0].NbrOfPurchasedItems)
		assert.Equal(t, 5, response.Summaries[0].TotalAmountEur)

		assert.Equal(t, "02", response.Summaries[1].CustomerID)
		assert.Equal(t, 1, response.Summaries[1].NbrOfPurchasedItems)
		assert.Equal(t, 5, response.Summaries[1].TotalAmountEur)
	})

	// Test case for error retrieving summaries
	t.Run("Summaries empty", func(t *testing.T) {
		// Create a new router that simulates an error
		router := gin.Default()
		testCollection = &collections.OrderCollection{}

		// Define the route for the error case
		router.GET("/summary", GetSummariesHandler(testCollection))

		req, _ := http.NewRequest("GET", "/summary", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Check the response status
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify the error message
		var response struct {
			Summaries []models.Summary `json:"summaries"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(response.Summaries))
	})
}
