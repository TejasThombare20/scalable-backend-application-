package catalog

import (
	"context"
	"log"

	"github.com/TejasThombare20/go-microservice/catalog/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	c := pb.NewCatalogServiceClient(conn)

	return &Client{conn, c}, nil

}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {

	response, err := c.service.PostProduct(ctx, &pb.PostProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
	})

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          response.Product.Id,
		Name:        response.Product.Name,
		Description: response.Product.Description,
		Price:       response.Product.Price,
	}, nil

}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {

	response, err := c.service.GetProduct(ctx, &pb.GetProductRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          response.Product.Id,
		Name:        response.Product.Name,
		Description: response.Product.Description,
		Price:       response.Product.Price,
	}, nil

}

func (c *Client) GetProducts(ctx context.Context, skip, take uint64, query string, ids []string) ([]Product, error) {

	log.Println("inside GetProducts clinets with ids", ids)

	response, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Ids:   ids,
		Skip:  skip,
		Take:  take,
		Query: query,
	})

	if err != nil {
		return nil, err
	}

	products := []Product{}

	for _, product := range response.Products {
		products = append(products, Product{
			ID:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return products, nil

}
