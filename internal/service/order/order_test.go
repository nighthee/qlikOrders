package order

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"qlikOrders/internal/collections"
	"qlikOrders/internal/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(collection *collections.OrderCollection) *gin.Engine {
	router := gin.Default()
	router.POST("/orders", AddOrdersHandler(collection))
	return router
}

func TestAddOrdersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		input        []models.Order
		expectedCode int
		expectedBody string
	}{
		{
			name: "Valid Input",
			input: []models.Order{
				{
					CustomerID: "01",
					OrderID:    "50",
					Timestamp:  "1637245070513",
					Items: []models.Item{
						{ItemID: "20201", CostEur: 2},
					},
				},
			},
			expectedCode: http.StatusCreated,
			expectedBody: `{"message":"Orders added successfully"}`,
		},
		{
			name: "Invalid Input - Missing CustomerID",
			input: []models.Order{
				{
					OrderID:   "50",
					Timestamp: "1637245070513",
					Items: []models.Item{
						{ItemID: "20201", CostEur: 2},
					},
				},
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Invalid input"}`,
		},
		{
			name: "Invalid Input - Empty Items",
			input: []models.Order{
				{
					CustomerID: "01",
					OrderID:    "50",
					Timestamp:  "1637245070513",
					Items:      []models.Item{},
				},
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Invalid input"}`,
		},
		{
			name: "Invalid Input - Batch Size Exceeds Limit",
			input: []models.Order{
				{CustomerID: "01", OrderID: "50", Timestamp: "1637245070513", Items: []models.Item{{ItemID: "20201", CostEur: 2}}},
				{CustomerID: "02", OrderID: "51", Timestamp: "1637245070514", Items: []models.Item{{ItemID: "20202", CostEur: 3}}},
				{CustomerID: "03", OrderID: "52", Timestamp: "1637245070515", Items: []models.Item{{ItemID: "20203", CostEur: 4}}},
				{CustomerID: "04", OrderID: "53", Timestamp: "1637245070516", Items: []models.Item{{ItemID: "20204", CostEur: 5}}},
				{CustomerID: "05", OrderID: "54", Timestamp: "1637245070517", Items: []models.Item{{ItemID: "20205", CostEur: 6}}},
				{CustomerID: "06", OrderID: "55", Timestamp: "1637245070518", Items: []models.Item{{ItemID: "20206", CostEur: 7}}},
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Batch size exceeds the allowed limit","message":"The maximum allowed number of orders in a single request is 5. Please split your request and try again."}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up a mock OrderCollection
			collection := &collections.OrderCollection{
				Orders: []models.Order{},
			}

			router := setupRouter(collection)

			// Convert the input to JSON
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Asserts the response code
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
