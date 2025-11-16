package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"orders/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
	logger         *slog.Logger
}

func NewProductHandler(productUseCase *usecase.ProductUseCase, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		logger:         logger,
	}
}

type CreateProductRequest struct {
	Name        string  `json:"name" example:"Laptop Dell Inspiron"`
	Description string  `json:"description" example:"Laptop com 16GB RAM e SSD 512GB"`
	Price       float64 `json:"price" example:"3500.00"`
	Stock       int     `json:"stock" example:"10"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" example:"Laptop Dell Inspiron Pro"`
	Description string  `json:"description" example:"Laptop com 32GB RAM e SSD 1TB"`
	Price       float64 `json:"price" example:"5000.00"`
	Stock       int     `json:"stock" example:"5"`
}

// Create godoc
// @Summary Create a new product
// @Description Create a new product with name, description, price and stock
// @Tags products
// @Accept json
// @Produce json
// @Param product body CreateProductRequest true "Product data"
// @Success 201 {object} entity.Product
// @Failure 400 {object} ErrorResponse
// @Router /products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	product, err := h.productUseCase.CreateProduct(req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		h.logger.Error("Failed to create product", "error", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info("Product created via API", "product_id", product.ID)
	respondWithJSON(w, http.StatusCreated, product)
}

// GetByID godoc
// @Summary Get product by ID
// @Description Get a single product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} entity.Product
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.logger.Info("Getting product by ID", "product_id", id)

	product, err := h.productUseCase.GetProduct(id)
	if err != nil {
		h.logger.Error("Product not found", "product_id", id, "error", err)
		respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	respondWithJSON(w, http.StatusOK, product)
}

// List godoc
// @Summary List all products
// @Description Get a list of all products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} entity.Product
// @Failure 500 {object} ErrorResponse
// @Router /products [get]
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Listing all products")

	products, err := h.productUseCase.ListProducts()
	if err != nil {
		h.logger.Error("Failed to list products", "error", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(products) == 0 {
		respondWithJSON(w, http.StatusOK, []interface{}{})
		return
	}

	h.logger.Info("Products listed successfully", "count", len(products))
	respondWithJSON(w, http.StatusOK, products)
}

// Update godoc
// @Summary Update a product
// @Description Update an existing product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body UpdateProductRequest true "Updated product data"
// @Success 200 {object} entity.Product
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.logger.Info("Updating product", "product_id", id)

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", "product_id", id, "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	product, err := h.productUseCase.UpdateProduct(id, req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		h.logger.Error("Failed to update product", "product_id", id, "error", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info("Product updated via API", "product_id", id)
	respondWithJSON(w, http.StatusOK, product)
}

// Delete godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.logger.Info("Deleting product", "product_id", id)

	err := h.productUseCase.DeleteProduct(id)
	if err != nil {
		h.logger.Error("Failed to delete product", "product_id", id, "error", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Info("Product deleted via API", "product_id", id)

	respondWithJSON(w, http.StatusNoContent, nil)
}
