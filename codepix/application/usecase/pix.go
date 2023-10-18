package usecase

import (
	"fmt"

	"github.com/gui-laranjeira/codepix/codepix/domain/model"
)

type PixUseCase struct {
	PixKeyRepository model.IPixKeyRepository
}

func (p *PixUseCase) RegisterKey(key string, kind string, accountId string) (*model.PixKey, error) {
	account, err := p.PixKeyRepository.FindAccount(accountId)
	if err != nil {
		return nil, err
	}

	pixKey, err := model.NewPixKey(kind, account, key)
	if err != nil {
		return nil, err
	}

	p.PixKeyRepository.RegisterKey(pixKey)

	if pixKey.ID == "" {
		return nil, fmt.Errorf("unable to create new key at the moment")
	}

	return pixKey, err
}

func (p *PixUseCase) FindKey(key string, kind string) (*model.PixKey, error) {
	pixKey, err := p.PixKeyRepository.FindKeyByKind(key, kind)
	if err != nil {
		return nil, err
	}

	return pixKey, nil
}
