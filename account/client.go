package account

import (
	"context"
	"log"

	"github.com/TejasThombare20/go-microservice/account/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {

	log.Println("url ", url)
	url1 := "account:8080"
	conn, err := grpc.Dial(url1, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	c := pb.NewAccountServiceClient(conn)

	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.PostAccount(ctx, &pb.PostAccountRequest{Name: name})

	if err != nil {
		log.Println("error posting account", err)
		return nil, err
	}

	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	log.Println("inside GetAccount with id ", id)

	response, err := c.service.GetAccount(ctx, &pb.GetAccountRequest{Id: id})

	if err != nil {
		return nil, err
	}

	log.Println("after error check ")

	return &Account{
		ID:   response.Account.Id,
		Name: response.Account.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {

	log.Println("inside GetAccounts with skip abd take ", skip, take)

	response, err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{Skip: skip, Take: take})

	if err != nil {
		return nil, err
	}

	accounts := []Account{}

	for _, account := range response.Accounts {
		accounts = append(accounts, Account{
			ID:   account.Id,
			Name: account.Name,
		})
	}

	return accounts, nil
}