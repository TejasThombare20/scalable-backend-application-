package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostProduct(ctx context.Context, name, description string, price float64) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	GetProductsByID(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
}

type catalogService struct {
	repository Repository
}

func NewService(repository Repository) *catalogService {
	return &catalogService{repository}
}

func (s *catalogService) PostProduct(ctx context.Context, name string, description string, price float64) (*Product, error) {

	product := &Product{
		Name:        name,
		Description: description,
		Price:       price,
		ID:          ksuid.New().String()}

	if err := s.repository.PutProduct(ctx, *product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *catalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repository.GetProductByID(ctx, id)
}

func (s *catalogService) GetProducts(ctx context.Context, skip, take uint64) ([]Product, error) {

	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListProducts(ctx, skip, take)
}

func (s *catalogService) GetProductsByID(ctx context.Context, ids []string) ([]Product, error) {

	return s.repository.ListProductWithIDs(ctx, ids)
}

func (s *catalogService) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {

	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return s.repository.SearchProducts(ctx, query, skip, take)
}
