package repository

import (
	"database/sql"
	"log/slog"
	"orders-go/internal/domain/entity"
	"time"
)

type ProductRepositoryMySQL struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewProductRepository(db *sql.DB, logger *slog.Logger) *ProductRepositoryMySQL {
	return &ProductRepositoryMySQL{
		db:     db,
		logger: logger,
	}
}

func (r *ProductRepositoryMySQL) Create(product *entity.Product) error {
	r.logger.Info("Creating product", "product_id", product.ID, "name", product.Name)

	query := `
		INSERT INTO products (id, name, description, price, stock, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CreatedAt,
		product.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create product", "product_id", product.ID, "error", err)
		return err
	}

	r.logger.Info("Product created successfully", "product_id", product.ID)
	return nil
}

func (r *ProductRepositoryMySQL) FindByID(id string) (*entity.Product, error) {
	r.logger.Info("Finding product by ID", "product_id", id)

	query := `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products
		WHERE id = ?
	`
	var product entity.Product
	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("Product not found", "product_id", id)
		} else {
			r.logger.Error("Failed to find product", "product_id", id, "error", err)
		}
		return nil, err
	}

	r.logger.Info("Product found", "product_id", id)
	return &product, nil
}

func (r *ProductRepositoryMySQL) FindAll() ([]entity.Product, error) {
	r.logger.Info("Finding all products")

	query := `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error("Failed to query products", "error", err)
		return nil, err
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan product row", "error", err)
			return nil, err
		}
		products = append(products, product)
	}

	r.logger.Info("Products found", "count", len(products))
	return products, nil
}

func (r *ProductRepositoryMySQL) Update(product *entity.Product) error {
	r.logger.Info("Updating product", "product_id", product.ID)

	query := `
		UPDATE products
		SET name = ?, description = ?, price = ?, stock = ?, updated_at = ?
		WHERE id = ?
	`
	product.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.UpdatedAt,
		product.ID,
	)

	if err != nil {
		r.logger.Error("Failed to update product", "product_id", product.ID, "error", err)
		return err
	}

	r.logger.Info("Product updated successfully", "product_id", product.ID)
	return nil
}

func (r *ProductRepositoryMySQL) Delete(id string) error {
	r.logger.Info("Deleting product", "product_id", id)

	query := `DELETE FROM products WHERE id = ?`
	_, err := r.db.Exec(query, id)

	if err != nil {
		r.logger.Error("Failed to delete product", "product_id", id, "error", err)
		return err
	}

	r.logger.Info("Product deleted successfully", "product_id", id)
	return nil
}
