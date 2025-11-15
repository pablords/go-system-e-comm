package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodDebitCard  PaymentMethod = "debit_card"
	PaymentMethodPix        PaymentMethod = "pix"
	PaymentMethodBoleto     PaymentMethod = "boleto"
	PaymentMethodPayPal     PaymentMethod = "paypal"
)

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusApproved   PaymentStatus = "approved"
	PaymentStatusDeclined   PaymentStatus = "declined"
	PaymentStatusCanceled   PaymentStatus = "canceled"
	PaymentStatusRefunded   PaymentStatus = "refunded"
)

var (
	ErrInvalidAmount           = errors.New("amount must be greater than zero")
	ErrInvalidPaymentMethod    = errors.New("invalid payment method")
	ErrInvalidPaymentStatus    = errors.New("invalid payment status")
	ErrPaymentNotFound         = errors.New("payment not found")
	ErrPaymentCannotBeCanceled = errors.New("payment cannot be canceled in current status")
	ErrEmptyOrderID            = errors.New("order_id cannot be empty")
	ErrEmptyCustomerEmail      = errors.New("customer email cannot be empty")
)

type Payment struct {
	ID            string        `json:"id"`
	OrderID       string        `json:"order_id"`
	Amount        float64       `json:"amount"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Status        PaymentStatus `json:"status"`
	TransactionID string        `json:"transaction_id,omitempty"`
	CustomerEmail string        `json:"customer_email"`
	CustomerName  string        `json:"customer_name"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	CanceledAt    *time.Time    `json:"canceled_at,omitempty"`
	CancelReason  string        `json:"cancel_reason,omitempty"`
}

func NewPayment(orderID string, amount float64, paymentMethod PaymentMethod, customerEmail, customerName string) (*Payment, error) {
	if err := validatePaymentData(orderID, amount, paymentMethod, customerEmail); err != nil {
		return nil, err
	}

	return &Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Status:        PaymentStatusPending,
		CustomerEmail: customerEmail,
		CustomerName:  customerName,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func validatePaymentData(orderID string, amount float64, paymentMethod PaymentMethod, customerEmail string) error {
	if orderID == "" {
		return ErrEmptyOrderID
	}

	if amount <= 0 {
		return ErrInvalidAmount
	}

	if !isValidPaymentMethod(paymentMethod) {
		return ErrInvalidPaymentMethod
	}

	if customerEmail == "" {
		return ErrEmptyCustomerEmail
	}

	return nil
}

func isValidPaymentMethod(method PaymentMethod) bool {
	switch method {
	case PaymentMethodCreditCard, PaymentMethodDebitCard, PaymentMethodPix, PaymentMethodBoleto, PaymentMethodPayPal:
		return true
	default:
		return false
	}
}

func (p *Payment) Process(transactionID string) error {
	if p.Status != PaymentStatusPending {
		return errors.New("payment must be in pending status to be processed")
	}

	p.Status = PaymentStatusProcessing
	p.TransactionID = transactionID
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Payment) Approve() error {
	if p.Status != PaymentStatusProcessing {
		return errors.New("payment must be in processing status to be approved")
	}

	p.Status = PaymentStatusApproved
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Payment) Decline() error {
	if p.Status != PaymentStatusProcessing {
		return errors.New("payment must be in processing status to be declined")
	}

	p.Status = PaymentStatusDeclined
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Payment) Cancel(reason string) error {
	if p.Status == PaymentStatusApproved || p.Status == PaymentStatusCanceled || p.Status == PaymentStatusRefunded {
		return ErrPaymentCannotBeCanceled
	}

	p.Status = PaymentStatusCanceled
	p.CancelReason = reason
	now := time.Now()
	p.CanceledAt = &now
	p.UpdatedAt = now
	return nil
}

func (p *Payment) Refund() error {
	if p.Status != PaymentStatusApproved {
		return errors.New("only approved payments can be refunded")
	}

	p.Status = PaymentStatusRefunded
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Payment) CanBeCanceled() bool {
	return p.Status != PaymentStatusApproved &&
		p.Status != PaymentStatusCanceled &&
		p.Status != PaymentStatusRefunded
}

func (p *Payment) IsFinalized() bool {
	return p.Status == PaymentStatusApproved ||
		p.Status == PaymentStatusDeclined ||
		p.Status == PaymentStatusCanceled ||
		p.Status == PaymentStatusRefunded
}
