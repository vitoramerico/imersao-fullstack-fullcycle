package models

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending   string = "pending"
	TransactionCompleted string = "completed"
	TransactionError     string = "error"
	TransactionConfirmed string = "confirmed"
)

type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transaction []Transaction
}

type Transaction struct {
	Base              `valid:"required"`
	AccountFrom       *Account `valid:"-"`
	Amount            float64  `json: "amount" valid:"notnull"`
	PixKeyTo          *PixKey  `valid:"-"`
	Status            string   `json: "status" valid:"notnull"`
	Description       string   `json: "description" valid:"notnull"`
	CancelDescription string   `json: "cancelDescription" valid:"-"`
}

func (transaction *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(transaction)

	if transaction.Amount <= 0 {
		return errors.New("the amound must be greater than 0")
	}

	if transaction.Status != TransactionPending && transaction.Status != TransactionCompleted && transaction.Status != TransactionError {
		return errors.New("invalid status for the transaction")
	}

	if transaction.PixKeyTo.AccountID == transaction.AccountFrom.ID {
		return errors.New("the source and destination account cannot be the same")
	}

	if err != nil {
		return err
	}

	return nil
}

func NewTransaction(accountFrom *Account, amount float64, pixKey *PixKey, description string) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom: accountFrom,
		Amount:      amount,
		PixKeyTo:    pixKey,
		Status:      TransactionPending,
		Description: description,
	}
	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	err := transaction.isValid()

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (transaction *Transaction) Completed() error {
	transaction.Status = TransactionCompleted
	transaction.UpdateddAt = time.Now()
	err := transaction.isValid()
	return err
}

func (transaction *Transaction) Confirmed() error {
	transaction.Status = TransactionConfirmed
	transaction.UpdateddAt = time.Now()
	err := transaction.isValid()
	return err
}

func (transaction *Transaction) Cancel(description string) error {
	transaction.Status = TransactionError
	transaction.Description = description
	transaction.UpdateddAt = time.Now()
	err := transaction.isValid()
	return err
}
