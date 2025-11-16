package entity_test

import (
	"testing"

	"payments/internal/domain/entity"
)

func TestNewPayment(t *testing.T) {
	tests := []struct {
		name          string
		orderID       string
		amount        float64
		paymentMethod entity.PaymentMethod
		customerEmail string
		customerName  string
		expectError   bool
	}{
		{
			name:          "Valid payment",
			orderID:       "order-123",
			amount:        100.50,
			paymentMethod: entity.PaymentMethodCreditCard,
			customerEmail: "test@example.com",
			customerName:  "Test User",
			expectError:   false,
		},
		{
			name:          "Invalid amount - zero",
			orderID:       "order-123",
			amount:        0,
			paymentMethod: entity.PaymentMethodCreditCard,
			customerEmail: "test@example.com",
			customerName:  "Test User",
			expectError:   true,
		},
		{
			name:          "Invalid amount - negative",
			orderID:       "order-123",
			amount:        -10.0,
			paymentMethod: entity.PaymentMethodCreditCard,
			customerEmail: "test@example.com",
			customerName:  "Test User",
			expectError:   true,
		},
		{
			name:          "Empty order ID",
			orderID:       "",
			amount:        100.0,
			paymentMethod: entity.PaymentMethodCreditCard,
			customerEmail: "test@example.com",
			customerName:  "Test User",
			expectError:   true,
		},
		{
			name:          "Empty customer email",
			orderID:       "order-123",
			amount:        100.0,
			paymentMethod: entity.PaymentMethodCreditCard,
			customerEmail: "",
			customerName:  "Test User",
			expectError:   true,
		},
		{
			name:          "Invalid payment method",
			orderID:       "order-123",
			amount:        100.0,
			paymentMethod: "invalid",
			customerEmail: "test@example.com",
			customerName:  "Test User",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment, err := entity.NewPayment(tt.orderID, tt.amount, tt.paymentMethod, tt.customerEmail, tt.customerName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if payment != nil {
					t.Errorf("Expected nil payment but got %+v", payment)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if payment == nil {
					t.Error("Expected payment but got nil")
					return
				}
				if payment.Status != entity.PaymentStatusPending {
					t.Errorf("Expected status %s but got %s", entity.PaymentStatusPending, payment.Status)
				}
				if payment.OrderID != tt.orderID {
					t.Errorf("Expected order_id %s but got %s", tt.orderID, payment.OrderID)
				}
				if payment.Amount != tt.amount {
					t.Errorf("Expected amount %.2f but got %.2f", tt.amount, payment.Amount)
				}
			}
		})
	}
}

func TestPaymentProcess(t *testing.T) {
	payment, _ := entity.NewPayment("order-123", 100.0, entity.PaymentMethodCreditCard, "test@example.com", "Test User")

	err := payment.Process("txn-123")
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}

	if payment.Status != entity.PaymentStatusProcessing {
		t.Errorf("Expected status %s but got %s", entity.PaymentStatusProcessing, payment.Status)
	}

	if payment.TransactionID != "txn-123" {
		t.Errorf("Expected transaction_id txn-123 but got %s", payment.TransactionID)
	}

	// Try to process again (should fail)
	err = payment.Process("txn-456")
	if err == nil {
		t.Error("Expected error when processing already processing payment")
	}
}

func TestPaymentApprove(t *testing.T) {
	payment, _ := entity.NewPayment("order-123", 100.0, entity.PaymentMethodCreditCard, "test@example.com", "Test User")
	payment.Process("txn-123")

	err := payment.Approve()
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}

	if payment.Status != entity.PaymentStatusApproved {
		t.Errorf("Expected status %s but got %s", entity.PaymentStatusApproved, payment.Status)
	}

	// Try to approve again (should fail)
	err = payment.Approve()
	if err == nil {
		t.Error("Expected error when approving already approved payment")
	}
}

func TestPaymentDecline(t *testing.T) {
	payment, _ := entity.NewPayment("order-123", 100.0, entity.PaymentMethodCreditCard, "test@example.com", "Test User")
	payment.Process("txn-123")

	err := payment.Decline()
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}

	if payment.Status != entity.PaymentStatusDeclined {
		t.Errorf("Expected status %s but got %s", entity.PaymentStatusDeclined, payment.Status)
	}
}

func TestPaymentCancel(t *testing.T) {
	payment, _ := entity.NewPayment("order-123", 100.0, entity.PaymentMethodCreditCard, "test@example.com", "Test User")

	err := payment.Cancel("Customer requested cancellation")
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}

	if payment.Status != entity.PaymentStatusCanceled {
		t.Errorf("Expected status %s but got %s", entity.PaymentStatusCanceled, payment.Status)
	}

	if payment.CancelReason != "Customer requested cancellation" {
		t.Errorf("Expected cancel reason 'Customer requested cancellation' but got %s", payment.CancelReason)
	}

	if payment.CanceledAt == nil {
		t.Error("Expected canceled_at to be set")
	}
}

func TestPaymentCannotBeCanceled(t *testing.T) {
	payment, _ := entity.NewPayment("order-123", 100.0, entity.PaymentMethodCreditCard, "test@example.com", "Test User")
	payment.Process("txn-123")
	payment.Approve()

	err := payment.Cancel("Trying to cancel approved payment")
	if err == nil {
		t.Error("Expected error when canceling approved payment")
	}

	if err != entity.ErrPaymentCannotBeCanceled {
		t.Errorf("Expected ErrPaymentCannotBeCanceled but got: %v", err)
	}
}

func TestPaymentCanBeCanceled(t *testing.T) {
	tests := []struct {
		name           string
		status         entity.PaymentStatus
		expectedResult bool
	}{
		{"Pending can be canceled", entity.PaymentStatusPending, true},
		{"Processing can be canceled", entity.PaymentStatusProcessing, true},
		{"Approved cannot be canceled", entity.PaymentStatusApproved, false},
		{"Declined can be canceled", entity.PaymentStatusDeclined, true},
		{"Canceled cannot be canceled", entity.PaymentStatusCanceled, false},
		{"Refunded cannot be canceled", entity.PaymentStatusRefunded, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment, _ := entity.NewPayment("order-123", 100.0, entity.PaymentMethodCreditCard, "test@example.com", "Test User")
			payment.Status = tt.status

			result := payment.CanBeCanceled()
			if result != tt.expectedResult {
				t.Errorf("Expected CanBeCanceled() to return %v but got %v for status %s", tt.expectedResult, result, tt.status)
			}
		})
	}
}
