package repository

import (
	"fmt"

	"github.com/gui-laranjeira/codepix/codepix/domain/model"
	"gorm.io/gorm"
)

type TransactionRepositoryDb struct {
	Db *gorm.DB
}

func (r *TransactionRepositoryDb) Register(transaction *model.Transaction) error {
	err := r.Db.Create(transaction).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *TransactionRepositoryDb) Save(transaction *model.Transaction) error {
	err := r.Db.Save(transaction).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepositoryDb) Find(id string) (*model.Transaction, error) {
	var t model.Transaction
	r.Db.Preload("AccountFrom.Bank").First(&t, "id = ?", id)
	if t.ID == "" {
		return nil, fmt.Errorf("no key was found")
	}
	return &t, nil
}
