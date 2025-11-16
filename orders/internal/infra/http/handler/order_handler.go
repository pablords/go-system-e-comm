package handler

import (
	"log/slog"
	"net/http"
	"orders/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	orderUseCase *usecase.OrderUseCase
	logger       *slog.Logger
}

func NewOrderHandler(orderUseCase *usecase.OrderUseCase, logger *slog.Logger) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
		logger:       logger,
	}
}

// GetByID godoc
// @Summary Get order by ID
// @Description Get a single order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} entity.Order
// @Failure 404 {object} ErrorResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.logger.Info("Getting order by ID", "order_id", id)

	order, err := h.orderUseCase.GetOrder(id)
	if err != nil {
		h.logger.Error("Order not found", "order_id", id, "error", err)
		respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	respondWithJSON(w, http.StatusOK, order)
}

// List godoc
// @Summary List all orders
// @Description Get a list of all orders
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {array} entity.Order
// @Failure 500 {object} ErrorResponse
// @Router /orders [get]
func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Listing all orders")

	orders, err := h.orderUseCase.ListOrders()
	if err != nil {
		h.logger.Error("Failed to list orders", "error", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Info("Orders listed successfully", "count", len(orders))
	respondWithJSON(w, http.StatusOK, orders)
}

// Delete godoc
// @Summary Delete an order
// @Description Delete an order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id} [delete]
func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.logger.Info("Deleting order", "order_id", id)

	err := h.orderUseCase.DeleteOrder(id)
	if err != nil {
		h.logger.Error("Failed to delete order", "order_id", id, "error", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Info("Order deleted via API", "order_id", id)
	respondWithJSON(w, http.StatusNoContent, nil)
}
