package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"qlikOrders/internal/collections"
	"qlikOrders/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type summariesResponse struct {
	Summaries []models.Summary `json:"summaries"`
}

type BadOrder struct {
	OrderID   string        `json:"orderId"`
	Timestamp string        `json:"timestamp"`
	Items     []models.Item `json:"items"`
}

func TestNewServer(t *testing.T) {
	testCollection := &collections.OrderCollection{}
	server := NewServer(testCollection)

	t.Run("Test AddOrdersHandler", func(t *testing.T) {
		order := models.Order{
			CustomerID: "01",
			OrderID:    "50",
			Timestamp:  "1637245070513",
			Items: []models.Item{
				{ItemID: "20201", CostEur: 2},
			},
		}
		payload, _ := json.Marshal([]models.Order{order})

		// Create a POST request to add orders
		req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		// Check the response
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Orders added successfully")
	})

	t.Run("Test AddOrdersHandler Bad Request", func(t *testing.T) {
		order := BadOrder{
			OrderID:   "50",
			Timestamp: "1637245070513",
			Items: []models.Item{
				{ItemID: "20201", CostEur: 2},
			},
		}
		payload, _ := json.Marshal([]BadOrder{order})

		// Create a POST request to add orders
		req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		// Check the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test GetItemsByCustomerHandler", func(t *testing.T) {
		// Create a GET request to retrieve items for a specific customer
		req, _ := http.NewRequest("GET", "/customer/01/items", nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		// Define a struct matching the expected response format
		type response struct {
			Items []models.Item `json:"items"`
		}
		var resp response

		// Parse the JSON response
		err := json.Unmarshal(w.Body.Bytes(), &resp)

		// Asserts
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, len(resp.Items) > 0)
		for _, item := range resp.Items {
			assert.Equal(t, "20201", item.ItemID)
			assert.Equal(t, 2, item.CostEur)
		}
	})

	t.Run("Test GetSummariesHandler", func(t *testing.T) {
		// Create a GET request to retrieve summaries
		req, _ := http.NewRequest("GET", "/summary", nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		// Parse the JSON response
		var resp summariesResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)

		// Asserts
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, len(resp.Summaries) > 0)

		// Check the contents
		for _, summary := range resp.Summaries {
			assert.NotEmpty(t, summary.CustomerID)
			assert.Equal(t, summary.NbrOfPurchasedItems, 1)
			assert.Equal(t, summary.TotalAmountEur, 2)
		}
	})
}
