package main

import (
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/TejasThombare20/go-microservice/account"
	"github.com/TejasThombare20/go-microservice/catalog"
	"github.com/TejasThombare20/go-microservice/order"
)

type Server struct {
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
}

func NewGraphQLServer(accountUrl, catalogUrl, orderUrl string) (*Server, error) {
	log.Println("Hello world")
	accountClient, err := account.NewClient(accountUrl)

	if err != nil {
		log.Println("Error creating client for account: ", err)
		return nil, err
	}

	catlogClient, err := catalog.NewClient(catalogUrl)

	if err != nil {
		log.Println("Error creating client for catlogClient: ", err)
		return nil, err
	}

	orderClient, err := order.NewClient(orderUrl)

	if err != nil {
		log.Println("Error creating client for orderClient: ", err)
		return nil, err
	}

	return &Server{
		accountClient,
		catlogClient,
		orderClient,
	}, nil

}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
