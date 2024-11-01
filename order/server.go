package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/TejasThombare20/go-microservice/account"
	"github.com/TejasThombare20/go-microservice/catalog"
	"github.com/TejasThombare20/go-microservice/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
}

func ListenGRPC(s Service, accountURL, catalogURL string, port int) error {

	log.Println("accountURL , catalogURL ", accountURL, catalogURL)

	accountURL = "account:8080"
	catalogURL = "catalog:8080"

	accountClient, err := account.NewClient(accountURL)

	if err != nil {
		log.Println("error creating account client", err)
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)

	if err != nil {
		accountClient.Close()
		log.Println("error creating catalog client", err)
		return err
	}

	log.Println("port", port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		log.Println("error while listening ", err)
		return err
	}

	serv := grpc.NewServer()

	pb.RegisterOrderServiceServer(serv, &grpcServer{
		UnimplementedOrderServiceServer: pb.UnimplementedOrderServiceServer{},
		service:                         s,
		accountClient:                   accountClient,
		catalogClient:                   catalogClient,
	})
	reflection.Register(serv)

	return serv.Serve(lis)

}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {

	log.Println("Inside post order server")
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)

	if err != nil {
		log.Println("error getting account", err)
		return nil, errors.New("account not found")
	}

	log.Println("Products", r.Product)
	productIDs := []string{}

	for _, p := range r.Product {
		productIDs = append(productIDs, p.ProductId)
	}

	log.Println("productsIDs", productIDs)

	orderedProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, "", productIDs)

	log.Println("orderedProducts", orderedProducts)

	if err != nil {
		log.Println("error getting products", err)
		return nil, errors.New("products not found")
	}

	log.Println("products found")

	products := []OrderedProduct{}

	for _, product := range orderedProducts {

		newProduct := OrderedProduct{
			ID:          product.ID,
			Quantity:    0,
			Price:       product.Price,
			Name:        product.Name,
			Description: product.Description,
		}

		for _, rp := range r.Product {
			if rp.ProductId == product.ID {
				newProduct.Quantity = rp.Quantity
				break
			}
		}

		if newProduct.Quantity != 0 {
			products = append(products, newProduct)
		}
	}

	log.Println("products", products)

	log.Println("before posting order")

	order, err := s.service.PostOrder(ctx, r.AccountId, products)

	if err != nil {
		return nil, errors.New("could not post order")
	}

	log.Println("After posting order")

	orderProto := &pb.Order{
		Id:         order.ID,
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_OrderProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()

	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Quantity:    p.Quantity,
		})
	}
	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil
}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, r *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {

	accountOrders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	productIDMap := map[string]bool{}

	for _, order := range accountOrders {
		for _, product := range order.Products {
			productIDMap[product.ID] = true
		}
	}

	productIDs := []string{}

	for id := range productIDMap {
		productIDs = append(productIDs, id)
	}

	products, err := s.catalogClient.GetProducts(ctx, 0, 0, "", productIDs)

	if err != nil {
		log.Println("Error getting products", err)
		return nil, err
	}

	orders := []*pb.Order{}

	for _, order := range accountOrders {
		op := &pb.Order{
			AccountId:  order.AccountID,
			Id:         order.ID,
			TotalPrice: order.TotalPrice,
			Products:   []*pb.Order_OrderProduct{},
		}

		op.CreatedAt, _ = order.CreatedAt.MarshalBinary()

		for _, product := range order.Products {
			for _, p := range products {
				if p.ID == product.ID {
					product.Name = p.Name
					product.Description = p.Description
					product.Price = p.Price

					break
				}
			}

			op.Products = append(op.Products, &pb.Order_OrderProduct{
				Id:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    product.Quantity,
			})
		}

		orders = append(orders, op)
	}
	return &pb.GetOrdersForAccountResponse{Order: orders}, nil

}
