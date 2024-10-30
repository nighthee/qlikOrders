package models

type Order struct {
	CustomerID string `json:"customerId" validate:"required"`
	OrderID    string `json:"orderId" validate:"required"`
	Timestamp  string `json:"timestamp" validate:"required"`
	Items      []Item `json:"items" validate:"required,dive,required"` // Validate items array
}

// Item struct to represent an item within an order
type Item struct {
	ItemID  string `json:"itemId"`
	CostEur int    `json:"costEur"`
}

type CustomerItem struct {
	CustomerID string `json:"customerId"`
	ItemID     string `json:"itemId"`
	CostEur    int    `json:"costEur"`
}

// Summary struct for customer summary
type Summary struct {
	CustomerID          string `json:"customerId"`
	NbrOfPurchasedItems int    `json:"nbrOfPurchasedItems"`
	TotalAmountEur      int    `json:"totalAmountEur"`
}
