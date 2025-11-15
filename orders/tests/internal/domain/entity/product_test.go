package entity

import (
	"orders-go/internal/domain/entity"
	"testing"
)

func TestNewProduct(t *testing.T) {
	tests := []struct {
		name        string
		productName string
		description string
		price       float64
		stock       int
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "valid product",
			productName: "Laptop",
			description: "Dell Inspiron",
			price:       1500.00,
			stock:       10,
			wantErr:     false,
		},
		{
			name:        "empty name",
			productName: "",
			description: "Test",
			price:       100.00,
			stock:       5,
			wantErr:     true,
			expectedErr: entity.ErrInvalidProductName,
		},
		{
			name:        "zero price",
			productName: "Test Product",
			description: "Test",
			price:       0,
			stock:       5,
			wantErr:     true,
			expectedErr: entity.ErrInvalidProductPrice,
		},
		{
			name:        "negative price",
			productName: "Test Product",
			description: "Test",
			price:       -10.00,
			stock:       5,
			wantErr:     true,
			expectedErr: entity.ErrInvalidProductPrice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := entity.NewProduct(tt.productName, tt.description, tt.price, tt.stock)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewProduct() expected error, got nil")
				}
				if err != tt.expectedErr {
					t.Errorf("NewProduct() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewProduct() unexpected error = %v", err)
				return
			}

			if product.Name != tt.productName {
				t.Errorf("NewProduct() name = %v, want %v", product.Name, tt.productName)
			}
			if product.Price != tt.price {
				t.Errorf("NewProduct() price = %v, want %v", product.Price, tt.price)
			}
			if product.Stock != tt.stock {
				t.Errorf("NewProduct() stock = %v, want %v", product.Stock, tt.stock)
			}
			if product.ID == "" {
				t.Error("NewProduct() ID should not be empty")
			}
		})
	}
}

func TestProduct_UpdateStock(t *testing.T) {
	product, _ := entity.NewProduct("Test", "Test", 100.00, 10)

	tests := []struct {
		name     string
		quantity int
		wantErr  bool
	}{
		{
			name:     "add stock",
			quantity: 5,
			wantErr:  false,
		},
		{
			name:     "remove stock",
			quantity: -3,
			wantErr:  false,
		},
		{
			name:     "insufficient stock",
			quantity: -100,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialStock := product.Stock
			err := product.UpdateStock(tt.quantity)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateStock() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateStock() unexpected error = %v", err)
			}

			expectedStock := initialStock + tt.quantity
			if product.Stock != expectedStock {
				t.Errorf("UpdateStock() stock = %v, want %v", product.Stock, expectedStock)
			}
		})
	}
}
