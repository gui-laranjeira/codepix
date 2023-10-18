package model

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

type ITransactionRepository interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transactions []Transaction
}

type Transaction struct {
	Base              `valid:"required"`
	AccountFrom       *Account `valid:"-"`
	AccountFromID     string   `gorm:"column:account_from_id;type:uuid;" valid:"notnull"`
	Amount            float64  `json:"amount" gorm:"type:float" valid:"notnull"`
	PixKeyTo          *PixKey  `valid:"-"`
	PixKeyIdTo        string   `gorm:"column:pix_key_id_to;type:uuid;" valid:"notnull"`
	Status            string   `json:"status" gorm:"type:varchar(20)" valid:"notnull"`
	Description       string   `json:"description" gorm:"type:varchar(255)" valid:"notnull"`
	CancelDescription string   `json:"cancel_description" gorm:"type:varchar(255)" valid:"notnull"`
}

func (t *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(t)
	if err != nil {
		return err
	}

	if t.Amount <= 0 {
		return errors.New("amount must be grater than 0")
	}

	if t.Status != TransactionPending && t.Status != TransactionCompleted && t.Status != TransactionError {
		return errors.New("invalid status for the transaction")
	}

	if t.PixKeyTo.AccountID == t.AccountFrom.ID {
		return errors.New("the source and destination account cannot be the same")
	}

	return nil
}

func NewTransaction(accountFrom *Account, amount float64, pixKeyTo *PixKey, description string, id string) (*Transaction, error) {
	t := Transaction{
		AccountFrom: accountFrom,
		Amount:      amount,
		PixKeyTo:    pixKeyTo,
		PixKeyIdTo:  pixKeyTo.ID,
		Status:      TransactionPending,
		Description: description,
	}

	if id == "" {
		t.ID = uuid.NewV4().String()
	} else {
		t.ID = id
	}

	t.CreatedAt = time.Now()

	err := t.isValid()
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (t *Transaction) Complete() error {
	t.Status = TransactionCompleted
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}

func (t *Transaction) Confirm() error {
	t.Status = TransactionConfirmed
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}

func (t *Transaction) Cancel(description string) error {
	t.Status = TransactionError
	t.CancelDescription = description
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}
