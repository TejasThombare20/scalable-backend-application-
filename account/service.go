package account

import (
	"context"
	"log"

	"github.com/segmentio/ksuid"
)

type Services interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type accountService struct {
	repository Repository
}

func NewService(r Repository) Services {
	return &accountService{r}
}

func (s *accountService) PostAccount(ctx context.Context, name string) (*Account, error) {

	log.Println("Inside post account service")

	account := &Account{
		Name: name,
		ID:   ksuid.New().String(),
	}

	err := s.repository.PutAccount(ctx, *account)

	if err != nil {
		log.Println("error while storing account details:", err)
		return nil, err
	}

	return account, nil
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	return s.repository.GetAccountByID(ctx, id)

}

func (s *accountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {

	if (take > 100) || (skip == 0 && take == 0) {
		take = 100
	}

	return s.repository.ListAccounts(ctx, skip, take)
}
