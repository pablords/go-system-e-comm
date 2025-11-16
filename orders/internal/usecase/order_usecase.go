package usecase

import (
	"log/slog"
	"orders/internal/domain/entity"
	"orders/internal/domain/repository"
)

type OrderUseCase struct {
	orderRepo repository.OrderRepository
	logger    *slog.Logger
}

func NewOrderUseCase(orderRepo repository.OrderRepository, logger *slog.Logger) *OrderUseCase {
	return &OrderUseCase{
		orderRepo: orderRepo,
		logger:    logger,
	}
}

func (uc *OrderUseCase) GetOrder(id string) (*entity.Order, error) {
	uc.logger.Info("Getting order", "order_id", id)
	return uc.orderRepo.FindByID(id)
}

func (uc *OrderUseCase) ListOrders() ([]entity.Order, error) {
	uc.logger.Info("Listing all orders")
	return uc.orderRepo.FindAll()
}

func (uc *OrderUseCase) UpdateOrderStatus(id string, status entity.OrderStatus) (*entity.Order, error) {
	uc.logger.Info("Updating order status", "order_id", id, "status", status)

	order, err := uc.orderRepo.FindByID(id)
	if err != nil {
		uc.logger.Error("Failed to find order for status update", "order_id", id, "error", err)
		return nil, err
	}

	err = order.UpdateStatus(status)
	if err != nil {
		uc.logger.Error("Failed to validate status update", "order_id", id, "status", status, "error", err)
		return nil, err
	}

	err = uc.orderRepo.Update(order)
	if err != nil {
		uc.logger.Error("Failed to update order status", "order_id", id, "error", err)
		return nil, err
	}

	uc.logger.Info("Order status updated successfully", "order_id", id, "status", status)
	return order, nil
}

func (uc *OrderUseCase) DeleteOrder(id string) error {
	uc.logger.Info("Deleting order", "order_id", id)

	err := uc.orderRepo.Delete(id)
	if err != nil {
		uc.logger.Error("Failed to delete order", "order_id", id, "error", err)
		return err
	}

	uc.logger.Info("Order deleted successfully", "order_id", id)
	return nil
}
