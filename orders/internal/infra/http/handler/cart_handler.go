package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"orders/internal/domain/entity"
	"orders/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type CartHandler struct {
	cartUseCase *usecase.CartUseCase
	logger      *slog.Logger
}

func NewCartHandler(cartUseCase *usecase.CartUseCase, logger *slog.Logger) *CartHandler {
	return &CartHandler{
		cartUseCase: cartUseCase,
		logger:      logger,
	}
}

type AddItemRequest struct {
	ProductID string `json:"product_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Quantity  int    `json:"quantity" example:"2"`
}

type UpdateItemRequest struct {
	Quantity int `json:"quantity" example:"5"`
}

// CreateCart godoc
// @Summary Create a new cart
// @Description Create a new shopping cart (order)
// @Tags cart
// @Accept json
// @Produce json
// @Success 201 {object} entity.Order
// @Failure 500 {object} ErrorResponse
// @Router /cart [post]
func (h *CartHandler) CreateCart(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Creating new cart")

	order, err := h.cartUseCase.CreateOrder()
	if err != nil {
		h.logger.Error("Failed to create cart", "error", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Info("Cart created via API", "order_id", order.ID)
	respondWithJSON(w, http.StatusCreated, order)
}

// GetCart godoc
// @Summary Get cart by ID
// @Description Get a shopping cart by its ID
// @Tags cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Success 200 {object} entity.Order
// @Failure 404 {object} ErrorResponse
// @Router /cart/{id} [get]
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	h.logger.Info("Getting cart", "order_id", orderID)

	order, err := h.cartUseCase.GetCart(orderID)
	if err != nil {
		h.logger.Error("Cart not found", "order_id", orderID, "error", err)
		respondWithError(w, http.StatusNotFound, "Cart not found")
		return
	}

	respondWithJSON(w, http.StatusOK, order)
}

// AddItem godoc
// @Summary Add item to cart
// @Description Add a product item to the shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Param item body AddItemRequest true "Item to add"
// @Success 200 {object} entity.Order
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /cart/{id}/items [post]
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	h.logger.Info("Adding item to cart", "order_id", orderID)

	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", "order_id", orderID, "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	order, err := h.cartUseCase.AddItemToCart(orderID, req.ProductID, req.Quantity)
	if err != nil {
		h.logger.Error("Failed to add item to cart", "order_id", orderID, "product_id", req.ProductID, "error", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info("Item added to cart via API", "order_id", orderID, "product_id", req.ProductID)
	respondWithJSON(w, http.StatusOK, order)
}

// RemoveItem godoc
// @Summary Remove item from cart
// @Description Remove an item from the shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Param itemId path string true "Item ID"
// @Success 200 {object} entity.Order
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /cart/{id}/items/{itemId} [delete]
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")
	h.logger.Info("Removing item from cart", "order_id", orderID, "item_id", itemID)

	order, err := h.cartUseCase.RemoveItemFromCart(orderID, itemID)
	if err != nil {
		h.logger.Error("Failed to remove item", "order_id", orderID, "item_id", itemID, "error", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info("Item removed via API", "order_id", orderID, "item_id", itemID)
	respondWithJSON(w, http.StatusOK, order)
}

// UpdateItemQuantity godoc
// @Summary Update item quantity
// @Description Update the quantity of an item in the cart
// @Tags cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Param itemId path string true "Item ID"
// @Param item body UpdateItemRequest true "New quantity"
// @Success 200 {object} entity.Order
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /cart/{id}/items/{itemId} [put]
func (h *CartHandler) UpdateItemQuantity(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")
	h.logger.Info("Updating item quantity", "order_id", orderID, "item_id", itemID)

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", "order_id", orderID, "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	order, err := h.cartUseCase.UpdateItemQuantity(orderID, itemID, req.Quantity)
	if err != nil {
		h.logger.Error("Failed to update quantity", "order_id", orderID, "item_id", itemID, "error", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info("Item quantity updated via API", "order_id", orderID, "item_id", itemID, "quantity", req.Quantity)
	respondWithJSON(w, http.StatusOK, order)
}

// CalculateTotal godoc
// @Summary Calculate cart total
// @Description Calculate the total amount for payment
// @Tags cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /cart/{id}/calculate [get]
func (h *CartHandler) CalculateTotal(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	h.logger.Info("Calculating total", "order_id", orderID)

	order, err := h.cartUseCase.CalculateTotal(orderID)
	if err != nil {
		h.logger.Error("Failed to calculate total", "order_id", orderID, "error", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := map[string]interface{}{
		"order_id": order.ID,
		"items":    order.Items,
		"total":    order.Total,
		"status":   order.Status,
	}

	h.logger.Info("Total calculated via API", "order_id", orderID, "total", order.Total)
	respondWithJSON(w, http.StatusOK, response)
}

type UpdateOrderStatusRequest struct {
	Status entity.OrderStatus `json:"status" example:"paid"`
}

// UpdateStatus godoc
// @Summary Update order status
// @Description Update the status of an order (pending, paid, canceled, completed)
// @Tags cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Param status body UpdateOrderStatusRequest true "New status"
// @Success 200 {object} entity.Order
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /cart/{id}/status [put]
func (h *CartHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	h.logger.Info("Updating order status", "order_id", orderID)

	var req UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", "order_id", orderID, "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Get order through cart use case and update status
	order, err := h.cartUseCase.GetCart(orderID)
	if err != nil {
		h.logger.Error("Order not found", "order_id", orderID, "error", err)
		respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	if err := order.UpdateStatus(req.Status); err != nil {
		h.logger.Error("Failed to update status", "order_id", orderID, "status", req.Status, "error", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info("Order status updated via API", "order_id", orderID, "status", req.Status)
	respondWithJSON(w, http.StatusOK, order)
}
