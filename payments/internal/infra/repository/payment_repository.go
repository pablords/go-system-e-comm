package repository

import (
	"context"
	"database/sql"
	"fmt"
	"payments/internal/domain/entity"
	"time"
)

type PaymentRepositoryMySQL struct {
	db *sql.DB
}

func NewPaymentRepositoryMySQL(db *sql.DB) *PaymentRepositoryMySQL {
	return &PaymentRepositoryMySQL{db: db}
}

func (r *PaymentRepositoryMySQL) Create(ctx context.Context, payment *entity.Payment) error {
	query := `
		INSERT INTO payments (id, order_id, amount, payment_method, status, transaction_id, 
		                     customer_email, customer_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		payment.ID,
		payment.OrderID,
		payment.Amount,
		payment.PaymentMethod,
		payment.Status,
		payment.TransactionID,
		payment.CustomerEmail,
		payment.CustomerName,
		payment.CreatedAt,
		payment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

func (r *PaymentRepositoryMySQL) FindByID(ctx context.Context, id string) (*entity.Payment, error) {
	query := `
		SELECT id, order_id, amount, payment_method, status, transaction_id,
		       customer_email, customer_name, created_at, updated_at, canceled_at, cancel_reason
		FROM payments
		WHERE id = ?
	`

	payment := &entity.Payment{}
	var canceledAt sql.NullTime
	var cancelReason sql.NullString
	var transactionID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Amount,
		&payment.PaymentMethod,
		&payment.Status,
		&transactionID,
		&payment.CustomerEmail,
		&payment.CustomerName,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&canceledAt,
		&cancelReason,
	)

	if err == sql.ErrNoRows {
		return nil, entity.ErrPaymentNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find payment: %w", err)
	}

	if transactionID.Valid {
		payment.TransactionID = transactionID.String
	}

	if canceledAt.Valid {
		payment.CanceledAt = &canceledAt.Time
	}

	if cancelReason.Valid {
		payment.CancelReason = cancelReason.String
	}

	return payment, nil
}

func (r *PaymentRepositoryMySQL) FindByOrderID(ctx context.Context, orderID string) ([]*entity.Payment, error) {
	query := `
		SELECT id, order_id, amount, payment_method, status, transaction_id,
		       customer_email, customer_name, created_at, updated_at, canceled_at, cancel_reason
		FROM payments
		WHERE order_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find payments by order_id: %w", err)
	}
	defer rows.Close()

	var payments []*entity.Payment

	for rows.Next() {
		payment := &entity.Payment{}
		var canceledAt sql.NullTime
		var cancelReason sql.NullString
		var transactionID sql.NullString

		err := rows.Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.Amount,
			&payment.PaymentMethod,
			&payment.Status,
			&transactionID,
			&payment.CustomerEmail,
			&payment.CustomerName,
			&payment.CreatedAt,
			&payment.UpdatedAt,
			&canceledAt,
			&cancelReason,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}

		if transactionID.Valid {
			payment.TransactionID = transactionID.String
		}

		if canceledAt.Valid {
			payment.CanceledAt = &canceledAt.Time
		}

		if cancelReason.Valid {
			payment.CancelReason = cancelReason.String
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PaymentRepositoryMySQL) Update(ctx context.Context, payment *entity.Payment) error {
	query := `
		UPDATE payments
		SET status = ?, transaction_id = ?, updated_at = ?, canceled_at = ?, cancel_reason = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		payment.Status,
		payment.TransactionID,
		time.Now(),
		payment.CanceledAt,
		payment.CancelReason,
		payment.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}

func (r *PaymentRepositoryMySQL) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM payments WHERE id = ?"

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}

	return nil
}

func (r *PaymentRepositoryMySQL) List(ctx context.Context) ([]*entity.Payment, error) {
	query := `
		SELECT id, order_id, amount, payment_method, status, transaction_id,
		       customer_email, customer_name, created_at, updated_at, canceled_at, cancel_reason
		FROM payments
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var payments []*entity.Payment

	for rows.Next() {
		payment := &entity.Payment{}
		var canceledAt sql.NullTime
		var cancelReason sql.NullString
		var transactionID sql.NullString

		err := rows.Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.Amount,
			&payment.PaymentMethod,
			&payment.Status,
			&transactionID,
			&payment.CustomerEmail,
			&payment.CustomerName,
			&payment.CreatedAt,
			&payment.UpdatedAt,
			&canceledAt,
			&cancelReason,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}

		if transactionID.Valid {
			payment.TransactionID = transactionID.String
		}

		if canceledAt.Valid {
			payment.CanceledAt = &canceledAt.Time
		}

		if cancelReason.Valid {
			payment.CancelReason = cancelReason.String
		}

		payments = append(payments, payment)
	}

	return payments, nil
}
