package main

import (
	"context"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	log.Println("inside query resolver accounts")

	defer cancel()

	if id != nil {

		log.Println("Inside query resolver with id: ", id)

		response, err := r.server.accountClient.GetAccount(ctx, *id)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		return []*Account{{
			ID:   response.ID,
			Name: response.Name,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)

	if pagination != nil {
		skip, take = pagination.bounds()
	}

	log.Println("After ID check ")

	accountList, err := r.server.accountClient.GetAccounts(ctx, skip, take)

	if err != nil {

		log.Println("error while gettting accounts via gRPC ", err)
		return nil, err
	}

	var accounts []*Account

	for _, a := range accountList {
		account := &Account{
			ID:   a.ID,
			Name: a.Name,
		}
		accounts = append(accounts, account)
	}

	return accounts, nil

}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	if id != nil {
		r, err := r.server.catalogClient.GetProduct(ctx, *id)

		if err != nil {
			log.Println(err)
		}

		return []*Product{{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Price:       r.Price,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	q := ""
	if query != nil {
		q = *query
	}
	productList, err := r.server.catalogClient.GetProducts(ctx, skip, take, q, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []*Product
	for _, a := range productList {
		products = append(products,
			&Product{
				ID:          a.ID,
				Name:        a.Name,
				Description: a.Description,
				Price:       a.Price,
			},
		)
	}

	return products, nil

}

func (p PaginationInput) bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(0)

	if p.Skip != nil {
		skipValue = uint64(*p.Skip)
	}

	if p.Take != nil {
		takeValue = uint64(*p.Take)
	}

	return skipValue, takeValue
}