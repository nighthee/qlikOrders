package server

import (
	"qlikOrders/internal/collections"
	"qlikOrders/internal/service/customer"
	"qlikOrders/internal/service/order"
	"qlikOrders/internal/service/summary"

	"github.com/gin-gonic/gin"
)

// NewServer creates a new HTTP server with the defined routes
func NewServer(collections *collections.OrderCollection) *gin.Engine {
	router := gin.Default()

	// Routes
	router.POST("/orders", order.AddOrdersHandler(collections))
	router.GET("/customer/:customerId/items", customer.GetItemsByCustomerHandler(collections))
	router.GET("/summary", summary.GetSummariesHandler(collections))

	return router
}
