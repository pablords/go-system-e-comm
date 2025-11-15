package repository

import (
	"database/sql"
	"orders-go/internal/domain/entity"
)

type ItemRepositoryMySQL struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepositoryMySQL {
	return &ItemRepositoryMySQL{db: db}
}

func (r *ItemRepositoryMySQL) Create(item *entity.Item) error {
	query := `
		INSERT INTO items (id, order_id, product_id, quantity, unit_price, total)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		item.ID,
		item.OrderID,
		item.ProductID,
		item.Quantity,
		item.UnitPrice,
		item.Total,
	)
	return err
}

func (r *ItemRepositoryMySQL) FindByID(id string) (*entity.Item, error) {
	query := `
		SELECT i.id, i.order_id, i.product_id, i.quantity, i.unit_price, i.total,
		       p.id, p.name, p.description, p.price, p.stock, p.created_at, p.updated_at
		FROM items i
		INNER JOIN products p ON i.product_id = p.id
		WHERE i.id = ?
	`
	var item entity.Item
	var product entity.Product
	err := r.db.QueryRow(query, id).Scan(
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
		return nil, err
	}
	item.Product = &product
	return &item, nil
}

func (r *ItemRepositoryMySQL) FindByOrderID(orderID string) ([]entity.Item, error) {
	query := `
		SELECT i.id, i.order_id, i.product_id, i.quantity, i.unit_price, i.total,
		       p.id, p.name, p.description, p.price, p.stock, p.created_at, p.updated_at
		FROM items i
		INNER JOIN products p ON i.product_id = p.id
		WHERE i.order_id = ?
	`
	rows, err := r.db.Query(query, orderID)
	if err != nil {
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
			return nil, err
		}
		item.Product = &product
		items = append(items, item)
	}
	return items, nil
}

func (r *ItemRepositoryMySQL) Update(item *entity.Item) error {
	query := `
		UPDATE items
		SET quantity = ?, unit_price = ?, total = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query,
		item.Quantity,
		item.UnitPrice,
		item.Total,
		item.ID,
	)
	return err
}

func (r *ItemRepositoryMySQL) Delete(id string) error {
	query := `DELETE FROM items WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
