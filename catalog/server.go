package catalog

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/TejasThombare20/go-microservice/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
	pb.UnimplementedCatalogServiceServer
}

func ListerGRPC(s Service, port int) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}

	serv := grpc.NewServer()

	pb.RegisterCatalogServiceServer(serv, &grpcServer{
		UnimplementedCatalogServiceServer: pb.UnimplementedCatalogServiceServer{},
		service:                           s,
	})

	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	product, err := s.service.PostProduct(ctx, r.Name, r.Description, r.Price)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.PostProductResponse{Product: &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}}, nil

}

func (s *grpcServer) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {

	product, err := s.service.GetProduct(ctx, r.Id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		},
	}, nil

}

func (s *grpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	var response []Product
	var err error
	log.Println("inside GetProducts server with ids", r.Ids)
	if r.Query != "" {
		response, err = s.service.SearchProducts(ctx, r.Query, r.Skip, r.Take)
	} else if len(r.Ids) != 0 {
		response, err = s.service.GetProductsByID(ctx, r.Ids)
	} else {
		response, err = s.service.GetProducts(ctx, r.Skip, r.Take)
	}

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	products := []*pb.Product{}

	for _, product := range response {
		products = append(products,
			&pb.Product{
				Id:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
	}

	return &pb.GetProductsResponse{Products: products}, nil

}