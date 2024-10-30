package customer

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

func TestGetItemsByCustomerHandler(t *testing.T) {
	// Set up the Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	testCollection := &collections.OrderCollection{}
	testCollection.Orders = append(testCollection.Orders, models.Order{
		CustomerID: "01",
		OrderID:    "50",
		Timestamp:  "1637245070513",
		Items: []models.Item{
			{ItemID: "20201", CostEur: 2},
		},
	})

	// Define the route for the test
	router.GET("/customer/:customerId/items", GetItemsByCustomerHandler(testCollection))

	// Test case for a valid customer
	t.Run("Valid Customer", func(t *testing.T) {

		// Create a GET request for a valid customer
		req, _ := http.NewRequest("GET", "/customer/01/items", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Check the response status
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify the response structure
		var response struct {
			Items []models.Item `json:"items"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response.Items))
		assert.Equal(t, "20201", response.Items[0].ItemID)
		assert.Equal(t, 2, response.Items[0].CostEur)
	})

	// Test case for an invalid customer
	t.Run("Invalid Customer", func(t *testing.T) {

		// Create a GET request for an invalid customer
		req, _ := http.NewRequest("GET", "/customer/99/items", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Check the response status
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Verify the error message
		var response struct {
			Error string `json:"error"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "customer not found or no items", response.Error)
	})
}
