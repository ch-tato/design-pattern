package main

import (
	"errors"
	"fmt"
)

// 1. The Product
type DeliveryOrder struct {
	OrderID       string
	Items         []string
	DeliveryNote  string
	PromoCode     string
	IsContactless bool
}

// 2. The Builder
type OrderBuilder struct {
	order DeliveryOrder
}

func NewOrderBuilder(orderID string) *OrderBuilder {
	return &OrderBuilder{
		order: DeliveryOrder{
			OrderID: orderID,
			Items:   []string{}, // Initialize empty slice
		},
	}
}

// Construction steps
func (b *OrderBuilder) AddItem(item string) *OrderBuilder {
	b.order.Items = append(b.order.Items, item)
	return b
}

func (b *OrderBuilder) WithDeliveryNote(note string) *OrderBuilder {
	b.order.DeliveryNote = note
	return b
}

func (b *OrderBuilder) ApplyPromoCode(code string) *OrderBuilder {
	b.order.PromoCode = code
	return b
}

func (b *OrderBuilder) SetContactless() *OrderBuilder {
	b.order.IsContactless = true
	return b
}

// 3. The Build Method with Validation
func (b *OrderBuilder) Build() (*DeliveryOrder, error) {
	// Business Logic: An order must have at least one item
	if len(b.order.Items) == 0 {
		return nil, errors.New("cannot build order: cart is empty")
	}

	// Return a copy or pointer to the finalized order
	return &b.order, nil
}

// Client Code
func main() {
	// Scenario A: A complex order with lots of customizations
	complexOrder, err := NewOrderBuilder("ORD-991").
		AddItem("Spicy Chicken Sandwich").
		AddItem("Large Fries").
		WithDeliveryNote("Leave at the side door, please don't knock.").
		ApplyPromoCode("SAVE20").
		SetContactless().
		Build()

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Complex Order Created: %+v\n", complexOrder)
	}

	// Scenario B: A simple order (just one item, no notes or promos)
	simpleOrder, _ := NewOrderBuilder("ORD-992").
		AddItem("Pepperoni Pizza").
		Build()

	fmt.Printf("Simple Order Created: %+v\n", simpleOrder)
}
