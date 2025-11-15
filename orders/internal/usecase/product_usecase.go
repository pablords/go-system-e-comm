package usecase

import (
	"log/slog"
	"orders-go/internal/domain/entity"
	"orders-go/internal/domain/repository"
)

type ProductUseCase struct {
	productRepo repository.ProductRepository
	logger      *slog.Logger
}

func NewProductUseCase(productRepo repository.ProductRepository, logger *slog.Logger) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
		logger:      logger,
	}
}

func (uc *ProductUseCase) CreateProduct(name, description string, price float64, stock int) (*entity.Product, error) {
	uc.logger.Info("Creating product", "name", name, "price", price, "stock", stock)

	product, err := entity.NewProduct(name, description, price, stock)
	if err != nil {
		uc.logger.Error("Failed to create product entity", "name", name, "error", err)
		return nil, err
	}

	err = uc.productRepo.Create(product)
	if err != nil {
		uc.logger.Error("Failed to save product", "product_id", product.ID, "error", err)
		return nil, err
	}

	uc.logger.Info("Product created successfully", "product_id", product.ID, "name", name)
	return product, nil
}

func (uc *ProductUseCase) GetProduct(id string) (*entity.Product, error) {
	uc.logger.Info("Getting product", "product_id", id)
	return uc.productRepo.FindByID(id)
}

func (uc *ProductUseCase) ListProducts() ([]entity.Product, error) {
	uc.logger.Info("Listing all products")
	return uc.productRepo.FindAll()
}

func (uc *ProductUseCase) UpdateProduct(id, name, description string, price float64, stock int) (*entity.Product, error) {
	uc.logger.Info("Updating product", "product_id", id)

	product, err := uc.productRepo.FindByID(id)
	if err != nil {
		uc.logger.Error("Failed to find product for update", "product_id", id, "error", err)
		return nil, err
	}

	product.Name = name
	product.Description = description
	product.Price = price
	product.Stock = stock

	if err := product.Validate(); err != nil {
		uc.logger.Error("Product validation failed", "product_id", id, "error", err)
		return nil, err
	}

	err = uc.productRepo.Update(product)
	if err != nil {
		uc.logger.Error("Failed to update product", "product_id", id, "error", err)
		return nil, err
	}

	uc.logger.Info("Product updated successfully", "product_id", id)
	return product, nil
}

func (uc *ProductUseCase) DeleteProduct(id string) error {
	uc.logger.Info("Deleting product", "product_id", id)

	err := uc.productRepo.Delete(id)
	if err != nil {
		uc.logger.Error("Failed to delete product", "product_id", id, "error", err)
		return err
	}

	uc.logger.Info("Product deleted successfully", "product_id", id)
	return nil
}
