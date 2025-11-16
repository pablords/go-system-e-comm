package entity

import (
	"orders/internal/domain/entity"
	"testing"
)

func TestNewOrder(t *testing.T) {
	order := entity.NewOrder()

	if order.ID == "" {
		t.Error("NewOrder() ID should not be empty")
	}
	if order.Status != entity.OrderStatusPending {
		t.Errorf("NewOrder() status = %v, want %v", order.Status, entity.OrderStatusPending)
	}
	if order.Total != 0 {
		t.Errorf("NewOrder() total = %v, want 0", order.Total)
	}
	if len(order.Items) != 0 {
		t.Errorf("NewOrder() items length = %v, want 0", len(order.Items))
	}
}

func TestOrder_AddItem(t *testing.T) {
	order := entity.NewOrder()
	product, _ := entity.NewProduct("Test Product", "Test", 100.00, 10)
	item, _ := entity.NewItem(order.ID, product.ID, product, 2)

	order.AddItem(item)

	if len(order.Items) != 1 {
		t.Errorf("AddItem() items length = %v, want 1", len(order.Items))
	}
	if order.Total != 200.00 {
		t.Errorf("AddItem() total = %v, want 200.00", order.Total)
	}

	// Add same product again
	item2, _ := entity.NewItem(order.ID, product.ID, product, 3)
	order.AddItem(item2)

	if len(order.Items) != 1 {
		t.Errorf("AddItem() items length = %v, want 1 (should merge)", len(order.Items))
	}
	if order.Items[0].Quantity != 5 {
		t.Errorf("AddItem() quantity = %v, want 5", order.Items[0].Quantity)
	}
	if order.Total != 500.00 {
		t.Errorf("AddItem() total = %v, want 500.00", order.Total)
	}
}

func TestOrder_RemoveItem(t *testing.T) {
	order := entity.NewOrder()
	product, _ := entity.NewProduct("Test Product", "Test", 100.00, 10)
	item, _ := entity.NewItem(order.ID, product.ID, product, 2)
	order.AddItem(item)

	// Remove existing item
	err := order.RemoveItem(order.Items[0].ID)
	if err != nil {
		t.Errorf("RemoveItem() unexpected error = %v", err)
	}
	if len(order.Items) != 0 {
		t.Errorf("RemoveItem() items length = %v, want 0", len(order.Items))
	}
	if order.Total != 0 {
		t.Errorf("RemoveItem() total = %v, want 0", order.Total)
	}

	// Remove non-existent item
	err = order.RemoveItem("non-existent-id")
	if err != entity.ErrItemNotFound {
		t.Errorf("RemoveItem() error = %v, want %v", err, entity.ErrItemNotFound)
	}
}

func TestOrder_UpdateItemQuantity(t *testing.T) {
	order := entity.NewOrder()
	product, _ := entity.NewProduct("Test Product", "Test", 100.00, 10)
	item, _ := entity.NewItem(order.ID, product.ID, product, 2)
	order.AddItem(item)

	// Update quantity
	err := order.UpdateItemQuantity(order.Items[0].ID, 5)
	if err != nil {
		t.Errorf("UpdateItemQuantity() unexpected error = %v", err)
	}
	if order.Items[0].Quantity != 5 {
		t.Errorf("UpdateItemQuantity() quantity = %v, want 5", order.Items[0].Quantity)
	}
	if order.Total != 500.00 {
		t.Errorf("UpdateItemQuantity() total = %v, want 500.00", order.Total)
	}

	// Invalid quantity
	err = order.UpdateItemQuantity(order.Items[0].ID, 0)
	if err != entity.ErrInvalidQuantity {
		t.Errorf("UpdateItemQuantity() error = %v, want %v", err, entity.ErrInvalidQuantity)
	}
}

func TestOrder_PrepareForPayment(t *testing.T) {
	// Empty order
	order := entity.NewOrder()
	err := order.PrepareForPayment()
	if err != entity.ErrEmptyOrder {
		t.Errorf("PrepareForPayment() error = %v, want %v", err, entity.ErrEmptyOrder)
	}

	// Order with items
	product, _ := entity.NewProduct("Test Product", "Test", 100.00, 10)
	item, _ := entity.NewItem(order.ID, product.ID, product, 2)
	order.AddItem(item)

	err = order.PrepareForPayment()
	if err != nil {
		t.Errorf("PrepareForPayment() unexpected error = %v", err)
	}
	if order.Total != 200.00 {
		t.Errorf("PrepareForPayment() total = %v, want 200.00", order.Total)
	}
}

func TestOrder_UpdateStatus(t *testing.T) {
	order := entity.NewOrder()

	tests := []struct {
		name    string
		status  entity.OrderStatus
		wantErr bool
	}{
		{
			name:    "valid status - paid",
			status:  entity.OrderStatusPaid,
			wantErr: false,
		},
		{
			name:    "valid status - completed",
			status:  entity.OrderStatusCompleted,
			wantErr: false,
		},
		{
			name:    "invalid status",
			status:  entity.OrderStatus("invalid"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := order.UpdateStatus(tt.status)

			if tt.wantErr {
				if err != entity.ErrInvalidOrderStatus {
					t.Errorf("UpdateStatus() error = %v, want %v", err, entity.ErrInvalidOrderStatus)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateStatus() unexpected error = %v", err)
			}
			if order.Status != tt.status {
				t.Errorf("UpdateStatus() status = %v, want %v", order.Status, tt.status)
			}
		})
	}
}
