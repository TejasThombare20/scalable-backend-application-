package order

import (
	"context"
	"log"
	"time"

	"github.com/TejasThombare20/go-microservice/order/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	c := pb.NewOrderServiceClient(conn)
	return &Client{conn, c}, nil

}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	protoProducts := []*pb.PostOrderRequest_OrderProduct{}

	log.Println("log account ID: ", accountID)

	log.Println("Inside post order client")

	for _, product := range products {
		protoProducts = append(protoProducts, &pb.PostOrderRequest_OrderProduct{
			ProductId: product.ID,
			Quantity:  product.Quantity,
		})
	}

	res, err := c.service.PostOrder(ctx,
		&pb.PostOrderRequest{
			AccountId: accountID,
			Product:   protoProducts,
		})

	if err != nil {
		return nil, err
	}

	newOrder := res.Order
	newOrderCreatedAt := time.Time{}

	newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)

	return &Order{
		ID:         newOrder.Id,
		CreatedAt:  newOrderCreatedAt,
		TotalPrice: newOrder.TotalPrice,
		AccountID:  newOrder.AccountId,
		Products:   products,
	}, nil

}

func (c *Client) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	response, err := c.service.GetOrdersForAccount(ctx, &pb.GetOrdersForAccountRequest{
		AccountId: accountID,
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	orders := []Order{}

	for _, orderProto := range response.Order {
		newOrder := Order{
			ID:         orderProto.Id,
			TotalPrice: orderProto.TotalPrice,
			AccountID:  orderProto.AccountId,
		}
		newOrder.CreatedAt = time.Time{}
		newOrder.CreatedAt.UnmarshalBinary(orderProto.CreatedAt)

		products := []OrderedProduct{}

		for _, product := range orderProto.Products {
			products = append(products, OrderedProduct{
				ID:          product.Id,
				Quantity:    product.Quantity,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
		newOrder.Products = products
		orders = append(orders, newOrder)

	}
	return orders, nil

}
