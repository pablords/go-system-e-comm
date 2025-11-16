package usecase

import (
	"errors"
	"log/slog"
	"orders/internal/domain/entity"
	"orders/internal/domain/repository"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type CartUseCase struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	logger      *slog.Logger
}

func NewCartUseCase(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	logger *slog.Logger,
) *CartUseCase {
	return &CartUseCase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		logger:      logger,
	}
}

// CreateOrder creates a new order (cart)
func (uc *CartUseCase) CreateOrder() (*entity.Order, error) {
	uc.logger.Info("Creating new cart/order")

	order := entity.NewOrder()
	err := uc.orderRepo.Create(order)
	if err != nil {
		uc.logger.Error("Failed to create cart/order", "error", err)
		return nil, err
	}

	uc.logger.Info("Cart/order created successfully", "order_id", order.ID)
	return order, nil
}

// AddItemToCart adds an item to the order/cart
func (uc *CartUseCase) AddItemToCart(orderID, productID string, quantity int) (*entity.Order, error) {
	uc.logger.Info("Adding item to cart", "order_id", orderID, "product_id", productID, "quantity", quantity)

	// Get order
	order, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		uc.logger.Error("Failed to find order", "order_id", orderID, "error", err)
		return nil, err
	}

	// Get product
	product, err := uc.productRepo.FindByID(productID)
	if err != nil {
		uc.logger.Error("Product not found", "product_id", productID, "error", err)
		return nil, ErrProductNotFound
	}

	// Create item
	item, err := entity.NewItem(orderID, productID, product, quantity)
	if err != nil {
		uc.logger.Error("Failed to create item", "product_id", productID, "error", err)
		return nil, err
	}

	// Add item to order
	order.AddItem(item)

	// Update order
	err = uc.orderRepo.Update(order)
	if err != nil {
		uc.logger.Error("Failed to update order with new item", "order_id", orderID, "error", err)
		return nil, err
	}

	uc.logger.Info("Item added to cart successfully", "order_id", orderID, "product_id", productID)
	return order, nil
}

// RemoveItemFromCart removes an item from the cart
func (uc *CartUseCase) RemoveItemFromCart(orderID, itemID string) (*entity.Order, error) {
	uc.logger.Info("Removing item from cart", "order_id", orderID, "item_id", itemID)

	order, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		uc.logger.Error("Failed to find order", "order_id", orderID, "error", err)
		return nil, err
	}

	err = order.RemoveItem(itemID)
	if err != nil {
		uc.logger.Error("Failed to remove item from order", "order_id", orderID, "item_id", itemID, "error", err)
		return nil, err
	}

	err = uc.orderRepo.Update(order)
	if err != nil {
		uc.logger.Error("Failed to update order after removing item", "order_id", orderID, "error", err)
		return nil, err
	}

	uc.logger.Info("Item removed from cart successfully", "order_id", orderID, "item_id", itemID)
	return order, nil
}

// UpdateItemQuantity updates the quantity of an item in the cart
func (uc *CartUseCase) UpdateItemQuantity(orderID, itemID string, quantity int) (*entity.Order, error) {
	uc.logger.Info("Updating item quantity", "order_id", orderID, "item_id", itemID, "quantity", quantity)

	order, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		uc.logger.Error("Failed to find order", "order_id", orderID, "error", err)
		return nil, err
	}

	err = order.UpdateItemQuantity(itemID, quantity)
	if err != nil {
		uc.logger.Error("Failed to update item quantity", "order_id", orderID, "item_id", itemID, "error", err)
		return nil, err
	}

	err = uc.orderRepo.Update(order)
	if err != nil {
		uc.logger.Error("Failed to update order after quantity change", "order_id", orderID, "error", err)
		return nil, err
	}

	uc.logger.Info("Item quantity updated successfully", "order_id", orderID, "item_id", itemID, "quantity", quantity)
	return order, nil
}

// CalculateTotal calculates and returns the total for payment
func (uc *CartUseCase) CalculateTotal(orderID string) (*entity.Order, error) {
	uc.logger.Info("Calculating total", "order_id", orderID)

	order, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		uc.logger.Error("Failed to find order for total calculation", "order_id", orderID, "error", err)
		return nil, err
	}

	err = order.PrepareForPayment()
	if err != nil {
		uc.logger.Error("Failed to prepare order for payment", "order_id", orderID, "error", err)
		return nil, err
	}

	uc.logger.Info("Total calculated successfully", "order_id", orderID, "total", order.Total)
	return order, nil
}

// GetCart retrieves the current cart/order
func (uc *CartUseCase) GetCart(orderID string) (*entity.Order, error) {
	return uc.orderRepo.FindByID(orderID)
}
