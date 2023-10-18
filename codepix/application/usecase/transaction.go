package usecase

import (
	"fmt"

	"github.com/gui-laranjeira/codepix/codepix/domain/model"
)

type TransactionUseCase struct {
	TransactionRepository model.ITransactionRepository
	PixRepository         model.IPixKeyRepository
}

func (u *TransactionUseCase) Register(accountId string, amount float64, pixKeyTo string, pixKeyKindTo string, description string) (*model.Transaction, error) {
	account, err := u.PixRepository.FindAccount(accountId)
	if err != nil {
		return nil, err
	}

	pixKey, err := u.PixRepository.FindKeyByKind(pixKeyTo, pixKeyKindTo)
	if err != nil {
		return nil, err
	}

	transaction, err := model.NewTransaction(account, amount, pixKey, description, "")
	if err != nil {
		return nil, err
	}

	u.TransactionRepository.Save(transaction)
	if transaction.ID == "" {
		return nil, fmt.Errorf("unable to process this transaction")
	}

	return transaction, nil
}

func (u *TransactionUseCase) Confirm(transactionId string) (*model.Transaction, error) {
	transaction, err := u.TransactionRepository.Find(transactionId)
	if err != nil {
		return nil, err
	}

	transaction.Status = model.TransactionConfirmed
	err = u.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (u *TransactionUseCase) Complete(transactionId string) (*model.Transaction, error) {
	transaction, err := u.TransactionRepository.Find(transactionId)
	if err != nil {
		return nil, err
	}

	transaction.Status = model.TransactionCompleted
	err = u.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (u *TransactionUseCase) Error(transactionId string, reason string) (*model.Transaction, error) {
	transaction, err := u.TransactionRepository.Find(transactionId)
	if err != nil {
		return nil, err
	}

	transaction.CancelDescription = reason
	transaction.Status = model.TransactionError
	err = u.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
