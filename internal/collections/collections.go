package collections

import (
	"errors"
	"qlikOrders/internal/models"
	"sync"
)

/*
According to https://pkg.go.dev/sync#Map, sync.Map is slower than regular maps for doing a lot of writes.
Will be using a regular map with locking rather than sync.Map.
*/

type Collections interface {
	AddOrders(newOrders []models.Order) error
	GetItemsByCustomer(customerID string) ([]models.Item, error)
	GetAllCustomerSummaries() ([]models.Summary, error)
}
type OrderCollection struct {
	Orders      []models.Order
	ordersMutex sync.Mutex
}

// AddOrders adds a batch of orders
func (o *OrderCollection) AddOrders(newOrders []models.Order) error {
	o.ordersMutex.Lock()
	defer o.ordersMutex.Unlock()

	for _, order := range newOrders {
		if err := validateOrder(order); err != nil {
			return err
		}
		o.Orders = append(o.Orders, order)
	}
	return nil
}

// GetItemsByCustomer retrieves items for a specific customer
func (o *OrderCollection) GetItemsByCustomer(customerID string) ([]models.CustomerItem, error) {
	o.ordersMutex.Lock()
	defer o.ordersMutex.Unlock()

	customerItems := []models.CustomerItem{}
	for _, order := range o.Orders {
		if order.CustomerID == customerID {
			for _, v := range order.Items {
				// Copy over the data and add with the customer ID
				customerItem := &models.CustomerItem{
					CustomerID: customerID,
					ItemID:     v.ItemID,
					CostEur:    v.CostEur,
				}
				customerItems = append(customerItems, *customerItem)
			}
		}
	}

	if len(customerItems) == 0 {
		return nil, errors.New("customer not found or no items")
	}
	return customerItems, nil
}

// GetAllCustomerSummaries provides summaries of all customers
func (o *OrderCollection) GetAllCustomerSummaries() ([]models.Summary, error) {
	o.ordersMutex.Lock()
	defer o.ordersMutex.Unlock()

	// Summarize orders for each customer using a map
	customerSummary := make(map[string]models.Summary)
	for _, order := range o.Orders { // Iterate over order collection

		// for every item in current order object
		for _, item := range order.Items {
			// Get the summary for a customer ID, if not exist, the map will return default values
			summary := customerSummary[order.CustomerID]

			// Update the summary object values
			summary.CustomerID = order.CustomerID
			summary.NbrOfPurchasedItems++
			summary.TotalAmountEur += item.CostEur

			// Update the value for the customerID
			customerSummary[order.CustomerID] = summary
		}
	}

	summaries := []models.Summary{}

	for _, summary := range customerSummary {
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

// Validates order structure to check for any missing fields
func validateOrder(order models.Order) error {
	if order.CustomerID == "" || order.OrderID == "" || order.Timestamp == "" || len(order.Items) == 0 {
		return errors.New("order is missing required fields or has invalid cost")
	}
	for _, item := range order.Items {
		if item.ItemID == "" || item.CostEur <= 0 {
			return errors.New("item is missing required fields or has invalid cost")
		}
	}
	return nil
}
