package main

import (
	"errors"
	"fmt"
)

// The Product
type DeliveryOrder struct {
	OrderID       string
	Items         []string
	DeliveryNote  string
	PromoCode     string
	IsContactless bool
}

// The Massive Constructor
// We have to include every possible option as a parameter here.
func NewDeliveryOrder(orderID string, items []string, deliveryNote string, promoCode string, isContactless bool) (*DeliveryOrder, error) {
	// Validation
	if len(items) == 0 {
		return nil, errors.New("cannot create order: cart is empty")
	}

	return &DeliveryOrder{
		OrderID:       orderID,
		Items:         items,
		DeliveryNote:  deliveryNote,
		PromoCode:     promoCode,
		IsContactless: isContactless,
	}, nil
}

func main() {
	// Scenario A: A complex order
	// This looks somewhat okay because we are using all the fields.
	complexOrder, err := NewDeliveryOrder(
		"ORD-991",
		[]string{"Spicy Chicken Sandwich", "Large Fries"},
		"Leave at the side door, please don't knock.",
		"SAVE20",
		true,
	)
	if err == nil {
		fmt.Printf("Complex Order Created: %+v\n", complexOrder)
	}

	// Scenario B: A simple order
	// DEMONSTRATION OF THE PROBLEM:
	// The user just wants a pizza. They don't have a note, no promo code,
	// and don't care about contactless delivery.
	// But because of our constructor, we are FORCED to pass empty strings ("") and false.
	simpleOrder, err := NewDeliveryOrder(
		"ORD-992",
		[]string{"Pepperoni Pizza"},
		"",    // What is this empty string for?
		"",    // And this one?
		false, // What does false mean here?
	)
	if err == nil {
		fmt.Printf("Simple Order Created: %+v\n", simpleOrder)
	}
}
