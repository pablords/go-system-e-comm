package usecase

import (
	"errors"
	"orders/internal/domain/entity"
	"orders/internal/usecase"
	"orders/tests/mocks"
	"testing"
)

// Mock Order Repository
type mockOrderRepository struct {
	orders map[string]*entity.Order
}

func newMockOrderRepository() *mockOrderRepository {
	return &mockOrderRepository{
		orders: make(map[string]*entity.Order),
	}
}

func (m *mockOrderRepository) Create(order *entity.Order) error {
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) FindByID(id string) (*entity.Order, error) {
	if order, ok := m.orders[id]; ok {
		return order, nil
	}
	return nil, errors.New("order not found")
}

func (m *mockOrderRepository) FindAll() ([]entity.Order, error) {
	orders := make([]entity.Order, 0, len(m.orders))
	for _, o := range m.orders {
		orders = append(orders, *o)
	}
	return orders, nil
}

func (m *mockOrderRepository) Update(order *entity.Order) error {
	if _, ok := m.orders[order.ID]; !ok {
		return errors.New("order not found")
	}
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) Delete(id string) error {
	if _, ok := m.orders[id]; !ok {
		return errors.New("order not found")
	}
	delete(m.orders, id)
	return nil
}

func TestCartUseCase_CreateOrder(t *testing.T) {
	orderRepo := newMockOrderRepository()
	productRepo := newMockProductRepository()
	logger := mocks.NewMockLogger()
	uc := usecase.NewCartUseCase(orderRepo, productRepo, logger)

	order, err := uc.CreateOrder()
	if err != nil {
		t.Errorf("CreateOrder() unexpected error = %v", err)
	}

	if order.ID == "" {
		t.Error("CreateOrder() ID should not be empty")
	}
	if order.Status != entity.OrderStatusPending {
		t.Errorf("CreateOrder() status = %v, want %v", order.Status, entity.OrderStatusPending)
	}
	if len(order.Items) != 0 {
		t.Errorf("CreateOrder() items length = %v, want 0", len(order.Items))
	}
}

func TestCartUseCase_AddItemToCart(t *testing.T) {
	orderRepo := newMockOrderRepository()
	productRepo := newMockProductRepository()
	logger := mocks.NewMockLogger()
	uc := usecase.NewCartUseCase(orderRepo, productRepo, logger)

	// Create order and product
	order, _ := uc.CreateOrder()
	product, _ := entity.NewProduct("Laptop", "Dell", 1500.00, 10)
	productRepo.Create(product)

	// Add item
	updatedOrder, err := uc.AddItemToCart(order.ID, product.ID, 2)
	if err != nil {
		t.Errorf("AddItemToCart() unexpected error = %v", err)
	}

	if len(updatedOrder.Items) != 1 {
		t.Errorf("AddItemToCart() items length = %v, want 1", len(updatedOrder.Items))
	}
	if updatedOrder.Items[0].Quantity != 2 {
		t.Errorf("AddItemToCart() quantity = %v, want 2", updatedOrder.Items[0].Quantity)
	}
	if updatedOrder.Total != 3000.00 {
		t.Errorf("AddItemToCart() total = %v, want 3000.00", updatedOrder.Total)
	}
}

func TestCartUseCase_AddItemToCart_ProductNotFound(t *testing.T) {
	orderRepo := newMockOrderRepository()
	productRepo := newMockProductRepository()
	logger := mocks.NewMockLogger()
	uc := usecase.NewCartUseCase(orderRepo, productRepo, logger)

	order, _ := uc.CreateOrder()

	_, err := uc.AddItemToCart(order.ID, "non-existent-product", 2)
	if err != usecase.ErrProductNotFound {
		t.Errorf("AddItemToCart() error = %v, want %v", err, usecase.ErrProductNotFound)
	}
}

func TestCartUseCase_RemoveItemFromCart(t *testing.T) {
	orderRepo := newMockOrderRepository()
	productRepo := newMockProductRepository()
	logger := mocks.NewMockLogger()
	uc := usecase.NewCartUseCase(orderRepo, productRepo, logger)

	// Create order and add item
	order, _ := uc.CreateOrder()
	product, _ := entity.NewProduct("Laptop", "Dell", 1500.00, 10)
	productRepo.Create(product)
	order, _ = uc.AddItemToCart(order.ID, product.ID, 2)

	// Remove item
	itemID := order.Items[0].ID
	updatedOrder, err := uc.RemoveItemFromCart(order.ID, itemID)
	if err != nil {
		t.Errorf("RemoveItemFromCart() unexpected error = %v", err)
	}

	if len(updatedOrder.Items) != 0 {
		t.Errorf("RemoveItemFromCart() items length = %v, want 0", len(updatedOrder.Items))
	}
	if updatedOrder.Total != 0 {
		t.Errorf("RemoveItemFromCart() total = %v, want 0", updatedOrder.Total)
	}
}

func TestCartUseCase_UpdateItemQuantity(t *testing.T) {
	orderRepo := newMockOrderRepository()
	productRepo := newMockProductRepository()
	logger := mocks.NewMockLogger()
	uc := usecase.NewCartUseCase(orderRepo, productRepo, logger)

	// Create order and add item
	order, _ := uc.CreateOrder()
	product, _ := entity.NewProduct("Laptop", "Dell", 1500.00, 10)
	productRepo.Create(product)
	order, _ = uc.AddItemToCart(order.ID, product.ID, 2)

	// Update quantity
	itemID := order.Items[0].ID
	updatedOrder, err := uc.UpdateItemQuantity(order.ID, itemID, 5)
	if err != nil {
		t.Errorf("UpdateItemQuantity() unexpected error = %v", err)
	}

	if updatedOrder.Items[0].Quantity != 5 {
		t.Errorf("UpdateItemQuantity() quantity = %v, want 5", updatedOrder.Items[0].Quantity)
	}
	if updatedOrder.Total != 7500.00 {
		t.Errorf("UpdateItemQuantity() total = %v, want 7500.00", updatedOrder.Total)
	}
}

func TestCartUseCase_CalculateTotal(t *testing.T) {
	orderRepo := newMockOrderRepository()
	productRepo := newMockProductRepository()
	logger := mocks.NewMockLogger()
	uc := usecase.NewCartUseCase(orderRepo, productRepo, logger)

	// Create order and add items
	order, _ := uc.CreateOrder()
	product1, _ := entity.NewProduct("Laptop", "Dell", 1500.00, 10)
	product2, _ := entity.NewProduct("Mouse", "Logitech", 50.00, 20)
	productRepo.Create(product1)
	productRepo.Create(product2)

	uc.AddItemToCart(order.ID, product1.ID, 2)
	uc.AddItemToCart(order.ID, product2.ID, 3)

	// Calculate total
	result, err := uc.CalculateTotal(order.ID)
	if err != nil {
		t.Errorf("CalculateTotal() unexpected error = %v", err)
	}

	expectedTotal := (1500.00 * 2) + (50.00 * 3)
	if result.Total != expectedTotal {
		t.Errorf("CalculateTotal() total = %v, want %v", result.Total, expectedTotal)
	}
}

func TestCartUseCase_CalculateTotal_EmptyCart(t *testing.T) {
	orderRepo := newMockOrderRepository()
	productRepo := newMockProductRepository()
	logger := mocks.NewMockLogger()
	uc := usecase.NewCartUseCase(orderRepo, productRepo, logger)

	order, _ := uc.CreateOrder()

	_, err := uc.CalculateTotal(order.ID)
	if err != entity.ErrEmptyOrder {
		t.Errorf("CalculateTotal() error = %v, want %v", err, entity.ErrEmptyOrder)
	}
}
