package order

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type Order struct {
	ID         string
	CreatedAt  time.Time
	TotalPrice float64
	AccountID  string
	Products   []OrderedProduct
}

type OrderedProduct struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    uint32
}

type orderService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &orderService{r}
}

func (s *orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	order := &Order{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		AccountID: accountID,
		Products:  products,
	}

	log.Println("products", products)

	order.TotalPrice = 0.0

	for _, product := range products {
		log.Println("product price  and Quantity", product.Price, product.Quantity, float64(product.Quantity))
		order.TotalPrice += product.Price * float64(product.Quantity)
		log.Println("order TotalPrice", order.TotalPrice)
	}

	log.Println(" final Total Price ", order.TotalPrice)

	err := s.repository.PutOrder(ctx, *order)

	if err != nil {
		return nil, err
	}

	return order, nil

}

func (s *orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {

	return s.repository.GetOrdersForAccount(ctx, accountID)
}
