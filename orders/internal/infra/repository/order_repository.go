package repository

import (
	"database/sql"
	"log/slog"
	"orders/internal/domain/entity"
	"time"
)

type OrderRepositoryMySQL struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewOrderRepository(db *sql.DB, logger *slog.Logger) *OrderRepositoryMySQL {
	return &OrderRepositoryMySQL{
		db:     db,
		logger: logger,
	}
}

func (r *OrderRepositoryMySQL) Create(order *entity.Order) error {
	r.logger.Info("Creating order", "order_id", order.ID, "items_count", len(order.Items))

	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error("Failed to begin transaction", "order_id", order.ID, "error", err)
		return err
	}
	defer tx.Rollback()

	// Insert order
	query := `
		INSERT INTO orders (id, status, total, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err = tx.Exec(query,
		order.ID,
		order.Status,
		order.Total,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		r.logger.Error("Failed to insert order", "order_id", order.ID, "error", err)
		return err
	}

	// Insert items
	for _, item := range order.Items {
		itemQuery := `
			INSERT INTO items (id, order_id, product_id, quantity, unit_price, total)
			VALUES (?, ?, ?, ?, ?, ?)
		`
		_, err = tx.Exec(itemQuery,
			item.ID,
			order.ID,
			item.ProductID,
			item.Quantity,
			item.UnitPrice,
			item.Total,
		)
		if err != nil {
			r.logger.Error("Failed to insert order item", "order_id", order.ID, "item_id", item.ID, "error", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("Failed to commit transaction", "order_id", order.ID, "error", err)
		return err
	}

	r.logger.Info("Order created successfully", "order_id", order.ID, "total", order.Total)
	return nil
}

func (r *OrderRepositoryMySQL) FindByID(id string) (*entity.Order, error) {
	r.logger.Info("Finding order by ID", "order_id", id)

	query := `
		SELECT id, status, total, created_at, updated_at
		FROM orders
		WHERE id = ?
	`
	var order entity.Order
	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.Status,
		&order.Total,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("Order not found", "order_id", id)
		} else {
			r.logger.Error("Failed to find order", "order_id", id, "error", err)
		}
		return nil, err
	}

	// Load items
	itemsQuery := `
		SELECT i.id, i.order_id, i.product_id, i.quantity, i.unit_price, i.total,
		       p.id, p.name, p.description, p.price, p.stock, p.created_at, p.updated_at
		FROM items i
		INNER JOIN products p ON i.product_id = p.id
		WHERE i.order_id = ?
	`
	rows, err := r.db.Query(itemsQuery, id)
	if err != nil {
		r.logger.Error("Failed to query order items", "order_id", id, "error", err)
		return nil, err
	}
	defer rows.Close()

	var items []entity.Item
	for rows.Next() {
		var item entity.Item
		var product entity.Product
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&item.Total,
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan order item", "order_id", id, "error", err)
			return nil, err
		}
		item.Product = &product
		items = append(items, item)
	}
	order.Items = items

	r.logger.Info("Order found", "order_id", id, "items_count", len(items))
	return &order, nil
}

func (r *OrderRepositoryMySQL) FindAll() ([]entity.Order, error) {
	r.logger.Info("Finding all orders")

	query := `
		SELECT id, status, total, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error("Failed to query orders", "error", err)
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		err := rows.Scan(
			&order.ID,
			&order.Status,
			&order.Total,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan order row", "error", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	r.logger.Info("Orders found", "count", len(orders))
	return orders, nil
}

func (r *OrderRepositoryMySQL) Update(order *entity.Order) error {
	r.logger.Info("Updating order", "order_id", order.ID)

	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error("Failed to begin transaction", "order_id", order.ID, "error", err)
		return err
	}
	defer tx.Rollback()

	order.UpdatedAt = time.Now()

	// Update order
	query := `
		UPDATE orders
		SET status = ?, total = ?, updated_at = ?
		WHERE id = ?
	`
	_, err = tx.Exec(query,
		order.Status,
		order.Total,
		order.UpdatedAt,
		order.ID,
	)
	if err != nil {
		r.logger.Error("Failed to update order", "order_id", order.ID, "error", err)
		return err
	}

	// Delete existing items
	_, err = tx.Exec(`DELETE FROM items WHERE order_id = ?`, order.ID)
	if err != nil {
		r.logger.Error("Failed to delete existing items", "order_id", order.ID, "error", err)
		return err
	}

	// Insert updated items
	for _, item := range order.Items {
		itemQuery := `
			INSERT INTO items (id, order_id, product_id, quantity, unit_price, total)
			VALUES (?, ?, ?, ?, ?, ?)
		`
		_, err = tx.Exec(itemQuery,
			item.ID,
			order.ID,
			item.ProductID,
			item.Quantity,
			item.UnitPrice,
			item.Total,
		)
		if err != nil {
			r.logger.Error("Failed to insert updated item", "order_id", order.ID, "item_id", item.ID, "error", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("Failed to commit transaction", "order_id", order.ID, "error", err)
		return err
	}

	r.logger.Info("Order updated successfully", "order_id", order.ID)
	return nil
}

func (r *OrderRepositoryMySQL) Delete(id string) error {
	r.logger.Info("Deleting order", "order_id", id)

	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error("Failed to begin transaction", "order_id", id, "error", err)
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM items WHERE order_id = ?`, id)
	if err != nil {
		r.logger.Error("Failed to delete order items", "order_id", id, "error", err)
		return err
	}

	_, err = tx.Exec(`DELETE FROM orders WHERE id = ?`, id)
	if err != nil {
		r.logger.Error("Failed to delete order", "order_id", id, "error", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("Failed to commit transaction", "order_id", id, "error", err)
		return err
	}

	r.logger.Info("Order deleted successfully", "order_id", id)
	return nil
}
