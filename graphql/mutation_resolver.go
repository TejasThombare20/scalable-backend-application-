package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/TejasThombare20/go-microservice/order"
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, in AccountInput) (*Account, error) {

	log.Println("Inside CreateAccount")
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	a, err := r.server.accountClient.PostAccount(ctx, in.Name)

	if err != nil {
		log.Println("error creating account via CreateAccount gRPC", err)
		return nil, err
	}

	return &Account{
		ID:   a.ID,
		Name: a.Name,
	}, nil

}

func (r *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	p, err := r.server.catalogClient.PostProduct(ctx, in.Name, in.Description, in.Price)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil

}

func (r *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	log.Println("Inside CreateOrder mutation resolver")

	var products []order.OrderedProduct

	for _, product := range in.Products {
		if product.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}

		products = append(products, order.OrderedProduct{
			ID:       product.ID,
			Quantity: uint32(product.Quantity),
		})
	}

	log.Println("product", products)

	order, err := r.server.orderClient.PostOrder(ctx, in.AccountID, products)

	if err != nil {
		log.Println("error while sending the psot order request via gRPC", err)
		return nil, err
	}

	return &Order{
		ID:         order.ID,
		CreatedAt:  order.CreatedAt,
		TotalPrice: order.TotalPrice,
	}, nil

}
