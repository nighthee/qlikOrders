package collections

import (
	"qlikOrders/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to reset the orders slice for each test case
func resetOrders(orderData *OrderCollection) {
	orderData.Orders = []models.Order{}
}

func TestAddOrders(t *testing.T) {
	orderCollection := &OrderCollection{}
	defer resetOrders(orderCollection) // Clean up after test

	t.Run("Add valid orders", func(t *testing.T) {
		orders := []models.Order{
			{
				CustomerID: "01",
				OrderID:    "100",
				Timestamp:  "1637245070513",
				Items: []models.Item{
					{ItemID: "item1", CostEur: 10},
				},
			},
		}

		err := orderCollection.AddOrders(orders)
		assert.Nil(t, err, "Expected no error")
	})

	t.Run("Invalid order missing Customer ID", func(t *testing.T) {
		orders := []models.Order{
			{
				CustomerID: "",
				OrderID:    "200",
				Timestamp:  "1637245070513",
				Items: []models.Item{
					{ItemID: "item2", CostEur: 10},
				},
			},
		}

		err := orderCollection.AddOrders(orders)
		assert.Error(t, err, "Expected error for invalid order")

	})
	t.Run("Invalid order missing Order ID", func(t *testing.T) {
		orders := []models.Order{
			{
				CustomerID: "01",
				OrderID:    "",
				Timestamp:  "1637245070513",
				Items: []models.Item{
					{ItemID: "item2", CostEur: 10},
				},
			},
		}
		err := orderCollection.AddOrders(orders)
		assert.Error(t, err, "Expected error for invalid order")
	})
	t.Run("Invalid order missing timestamp", func(t *testing.T) {
		orders := []models.Order{
			{
				CustomerID: "01",
				OrderID:    "200",
				Timestamp:  "",
				Items: []models.Item{
					{ItemID: "item2", CostEur: 10},
				},
			},
		}
		err := orderCollection.AddOrders(orders)
		assert.Error(t, err, "Expected error for invalid order")
	})
	t.Run("Invalid order missing item ID", func(t *testing.T) {
		orders := []models.Order{
			{
				CustomerID: "01",
				OrderID:    "200",
				Timestamp:  "1637245070513",
				Items: []models.Item{
					{ItemID: "", CostEur: 10},
				},
			},
		}
		err := orderCollection.AddOrders(orders)
		assert.Error(t, err, "Expected error for invalid order")
	})
	t.Run("Invalid order missing invalid currency", func(t *testing.T) {
		orders := []models.Order{
			{
				CustomerID: "01",
				OrderID:    "200",
				Timestamp:  "1637245070513",
				Items: []models.Item{
					{ItemID: "item2", CostEur: -1},
				},
			},
		}
		err := orderCollection.AddOrders(orders)
		assert.Error(t, err, "Expected error for invalid order")
	})
}

func TestGetItemsByCustomer(t *testing.T) {

	orderCollection := &OrderCollection{}
	defer resetOrders(orderCollection) // Clean up after test

	// Seed with sample orders
	orders := []models.Order{
		{
			CustomerID: "01",
			OrderID:    "100",
			Timestamp:  "1637245070513",
			Items: []models.Item{
				{ItemID: "item1", CostEur: 10},
			},
		},
		{
			CustomerID: "01",
			OrderID:    "101",
			Timestamp:  "1637245070523",
			Items: []models.Item{
				{ItemID: "item2", CostEur: 5},
			},
		},
	}
	orderCollection.AddOrders(orders)

	t.Run("Get items by existing customer", func(t *testing.T) {
		items, err := orderCollection.GetItemsByCustomer("01")

		assert.NoError(t, err, "Expected no error")
		assert.Len(t, items, 2, "Expected 2 items")
	})

	t.Run("Get items by non-existing customer", func(t *testing.T) {
		resetOrders(orderCollection)
		_, err := orderCollection.GetItemsByCustomer("01")

		assert.Error(t, err, "Expected error retrieving non existing customer")
	})
}

func TestGetAllCustomerSummaries(t *testing.T) {

	orderCollection := &OrderCollection{}
	defer resetOrders(orderCollection) // Clean up after test

	// Seed with sample orders
	orders := []models.Order{
		{
			CustomerID: "01",
			OrderID:    "100",
			Timestamp:  "1637245070513",
			Items: []models.Item{
				{ItemID: "item1", CostEur: 10},
				{ItemID: "item2", CostEur: 5},
			},
		},
		{
			CustomerID: "02",
			OrderID:    "200",
			Timestamp:  "1637245070533",
			Items: []models.Item{
				{ItemID: "item3", CostEur: 20},
			},
		},
	}

	// Using the map makes comparisons easier
	expectedSummaryDataMap := map[string]models.Summary{
		"01": {
			CustomerID:          "01",
			NbrOfPurchasedItems: 2,
			TotalAmountEur:      15,
		},
		"02": {
			CustomerID:          "02",
			NbrOfPurchasedItems: 1,
			TotalAmountEur:      20,
		},
	}

	orderCollection.AddOrders(orders)

	t.Run("Get summaries of all customers", func(t *testing.T) {
		summaries, err := orderCollection.GetAllCustomerSummaries()

		assert.NoError(t, err)
		assert.Len(t, summaries, 2, "expected length of summaries is 2")
		for _, summary := range summaries {
			assert.EqualValues(t, expectedSummaryDataMap[summary.CustomerID], summary, "expected returned summary to be equal")
		}
	})

	t.Run("Get summaries of all customers empty case", func(t *testing.T) {
		resetOrders(orderCollection)
		summaries, err := orderCollection.GetAllCustomerSummaries()

		assert.NoError(t, err)
		assert.Len(t, summaries, 0, "length should be 0")
	})
}
