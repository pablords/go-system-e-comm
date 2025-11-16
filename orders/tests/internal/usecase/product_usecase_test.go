package usecase

import (
	"errors"
	"orders/internal/domain/entity"
	"orders/internal/usecase"
	"orders/tests/mocks"
	"testing"
)

// Mock Repository
type mockProductRepository struct {
	products map[string]*entity.Product
}

func newMockProductRepository() *mockProductRepository {
	return &mockProductRepository{
		products: make(map[string]*entity.Product),
	}
}

func (m *mockProductRepository) Create(product *entity.Product) error {
	m.products[product.ID] = product
	return nil
}

func (m *mockProductRepository) FindByID(id string) (*entity.Product, error) {
	if product, ok := m.products[id]; ok {
		return product, nil
	}
	return nil, errors.New("product not found")
}

func (m *mockProductRepository) FindAll() ([]entity.Product, error) {
	products := make([]entity.Product, 0, len(m.products))
	for _, p := range m.products {
		products = append(products, *p)
	}
	return products, nil
}

func (m *mockProductRepository) Update(product *entity.Product) error {
	if _, ok := m.products[product.ID]; !ok {
		return errors.New("product not found")
	}
	m.products[product.ID] = product
	return nil
}

func (m *mockProductRepository) Delete(id string) error {
	if _, ok := m.products[id]; !ok {
		return errors.New("product not found")
	}
	delete(m.products, id)
	return nil
}

func TestProductUseCase_CreateProduct(t *testing.T) {
	repo := newMockProductRepository()
	logger := mocks.NewMockLogger()
	uc := usecase.NewProductUseCase(repo, logger)

	product, err := uc.CreateProduct("Laptop", "Dell Inspiron", 1500.00, 10)
	if err != nil {
		t.Errorf("CreateProduct() unexpected error = %v", err)
	}

	if product.Name != "Laptop" {
		t.Errorf("CreateProduct() name = %v, want Laptop", product.Name)
	}
	if product.Price != 1500.00 {
		t.Errorf("CreateProduct() price = %v, want 1500.00", product.Price)
	}
	if product.Stock != 10 {
		t.Errorf("CreateProduct() stock = %v, want 10", product.Stock)
	}

	// Verify it's in repository
	saved, err := repo.FindByID(product.ID)
	if err != nil {
		t.Errorf("Product not saved in repository")
	}
	if saved.Name != "Laptop" {
		t.Errorf("Saved product name = %v, want Laptop", saved.Name)
	}
}

func TestProductUseCase_CreateProduct_InvalidData(t *testing.T) {
	repo := newMockProductRepository()
	logger := mocks.NewMockLogger()
	uc := usecase.NewProductUseCase(repo, logger)

	// Empty name
	_, err := uc.CreateProduct("", "Test", 100.00, 10)
	if err == nil {
		t.Error("CreateProduct() expected error for empty name")
	}

	// Invalid price
	_, err = uc.CreateProduct("Test", "Test", 0, 10)
	if err == nil {
		t.Error("CreateProduct() expected error for zero price")
	}
}

func TestProductUseCase_GetProduct(t *testing.T) {
	repo := newMockProductRepository()
	logger := mocks.NewMockLogger()
	uc := usecase.NewProductUseCase(repo, logger)

	// Create product
	created, _ := uc.CreateProduct("Laptop", "Dell Inspiron", 1500.00, 10)

	// Get product
	product, err := uc.GetProduct(created.ID)
	if err != nil {
		t.Errorf("GetProduct() unexpected error = %v", err)
	}
	if product.ID != created.ID {
		t.Errorf("GetProduct() id = %v, want %v", product.ID, created.ID)
	}

	// Get non-existent product
	_, err = uc.GetProduct("non-existent-id")
	if err == nil {
		t.Error("GetProduct() expected error for non-existent product")
	}
}

func TestProductUseCase_ListProducts(t *testing.T) {
	repo := newMockProductRepository()
	logger := mocks.NewMockLogger()
	uc := usecase.NewProductUseCase(repo, logger)

	// Create products
	uc.CreateProduct("Laptop", "Dell", 1500.00, 10)
	uc.CreateProduct("Mouse", "Logitech", 50.00, 20)

	products, err := uc.ListProducts()
	if err != nil {
		t.Errorf("ListProducts() unexpected error = %v", err)
	}
	if len(products) != 2 {
		t.Errorf("ListProducts() length = %v, want 2", len(products))
	}
}

func TestProductUseCase_UpdateProduct(t *testing.T) {
	repo := newMockProductRepository()
	logger := mocks.NewMockLogger()
	uc := usecase.NewProductUseCase(repo, logger)

	// Create product
	created, _ := uc.CreateProduct("Laptop", "Dell", 1500.00, 10)

	// Update product
	updated, err := uc.UpdateProduct(created.ID, "Laptop Pro", "Dell Inspiron", 2000.00, 15)
	if err != nil {
		t.Errorf("UpdateProduct() unexpected error = %v", err)
	}
	if updated.Name != "Laptop Pro" {
		t.Errorf("UpdateProduct() name = %v, want Laptop Pro", updated.Name)
	}
	if updated.Price != 2000.00 {
		t.Errorf("UpdateProduct() price = %v, want 2000.00", updated.Price)
	}

	// Update non-existent product
	_, err = uc.UpdateProduct("non-existent", "Test", "Test", 100.00, 10)
	if err == nil {
		t.Error("UpdateProduct() expected error for non-existent product")
	}
}

func TestProductUseCase_DeleteProduct(t *testing.T) {
	repo := newMockProductRepository()
	logger := mocks.NewMockLogger()
	uc := usecase.NewProductUseCase(repo, logger)

	// Create product
	created, _ := uc.CreateProduct("Laptop", "Dell", 1500.00, 10)

	// Delete product
	err := uc.DeleteProduct(created.ID)
	if err != nil {
		t.Errorf("DeleteProduct() unexpected error = %v", err)
	}

	// Verify it's deleted
	_, err = repo.FindByID(created.ID)
	if err == nil {
		t.Error("Product should be deleted")
	}

	// Delete non-existent product
	err = uc.DeleteProduct("non-existent")
	if err == nil {
		t.Error("DeleteProduct() expected error for non-existent product")
	}
}
