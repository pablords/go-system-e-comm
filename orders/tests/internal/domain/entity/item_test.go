package entity

import (
	"orders/internal/domain/entity"
	"testing"
)

func TestNewItem(t *testing.T) {
	product, _ := entity.NewProduct("Test Product", "Test", 100.00, 10)
	tests := []struct {
		name        string
		orderID     string
		productID   string
		product     *entity.Product
		quantity    int
		wantErr     bool
		expectedErr error
	}{
		{
			name:      "valid item",
			orderID:   "order-123",
			productID: product.ID,
			product:   product,
			quantity:  2,
			wantErr:   false,
		},
		{
			name:        "nil product",
			orderID:     "order-123",
			productID:   "prod-123",
			product:     nil,
			quantity:    2,
			wantErr:     true,
			expectedErr: entity.ErrInvalidProduct,
		},
		{
			name:        "zero quantity",
			orderID:     "order-123",
			productID:   product.ID,
			product:     product,
			quantity:    0,
			wantErr:     true,
			expectedErr: entity.ErrInvalidQuantity,
		},
		{
			name:        "negative quantity",
			orderID:     "order-123",
			productID:   product.ID,
			product:     product,
			quantity:    -1,
			wantErr:     true,
			expectedErr: entity.ErrInvalidQuantity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := entity.NewItem(tt.orderID, tt.productID, tt.product, tt.quantity)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewItem() expected error, got nil")
				}
				if err != tt.expectedErr {
					t.Errorf("NewItem() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewItem() unexpected error = %v", err)
				return
			}

			if item.ID == "" {
				t.Error("NewItem() ID should not be empty")
			}
			if item.OrderID != tt.orderID {
				t.Errorf("NewItem() orderID = %v, want %v", item.OrderID, tt.orderID)
			}
			if item.ProductID != tt.productID {
				t.Errorf("NewItem() productID = %v, want %v", item.ProductID, tt.productID)
			}
			if item.Quantity != tt.quantity {
				t.Errorf("NewItem() quantity = %v, want %v", item.Quantity, tt.quantity)
			}
			if item.UnitPrice != product.Price {
				t.Errorf("NewItem() unitPrice = %v, want %v", item.UnitPrice, product.Price)
			}
			expectedTotal := product.Price * float64(tt.quantity)
			if item.Total != expectedTotal {
				t.Errorf("NewItem() total = %v, want %v", item.Total, expectedTotal)
			}
		})
	}
}

func TestItem_CalculateTotal(t *testing.T) {
	product, _ := entity.NewProduct("Test Product", "Test", 50.00, 10)
	item, _ := entity.NewItem("order-123", product.ID, product, 3)

	if item.Total != 150.00 {
		t.Errorf("CalculateTotal() total = %v, want 150.00", item.Total)
	}

	// Change quantity and recalculate
	item.Quantity = 5
	item.CalculateTotal()

	if item.Total != 250.00 {
		t.Errorf("CalculateTotal() total = %v, want 250.00", item.Total)
	}
}

func TestItem_UpdateQuantity(t *testing.T) {
	product, _ := entity.NewProduct("Test Product", "Test", 50.00, 10)
	item, _ := entity.NewItem("order-123", product.ID, product, 3)

	// Valid update
	err := item.UpdateQuantity(5)
	if err != nil {
		t.Errorf("UpdateQuantity() unexpected error = %v", err)
	}
	if item.Quantity != 5 {
		t.Errorf("UpdateQuantity() quantity = %v, want 5", item.Quantity)
	}
	if item.Total != 250.00 {
		t.Errorf("UpdateQuantity() total = %v, want 250.00", item.Total)
	}

	// Invalid quantity
	err = item.UpdateQuantity(0)
	if err != entity.ErrInvalidQuantity {
		t.Errorf("UpdateQuantity() error = %v, want %v", err, entity.ErrInvalidQuantity)
	}

	err = item.UpdateQuantity(-1)
	if err != entity.ErrInvalidQuantity {
		t.Errorf("UpdateQuantity() error = %v, want %v", err, entity.ErrInvalidQuantity)
	}
}
