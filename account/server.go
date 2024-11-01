package account

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/TejasThombare20/go-microservice/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service Services
}

func ListernGRPC(s Services, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterAccountServiceServer(serv, &grpcServer{
		UnimplementedAccountServiceServer: pb.UnimplementedAccountServiceServer{},
		service:                           s,
	})
	reflection.Register(serv)
	return serv.Serve(lis)

}

func (s *grpcServer) PostAccount(ctx context.Context, r *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {

	log.Println("Inside postAccount server")

	if s.service == nil {
		log.Println("Error: s.service is nil")
		return nil, errors.New("service not initialized")
	}

	a, err := s.service.PostAccount(ctx, r.Name)

	log.Println("after postAccount service call")

	if err != nil {
		log.Println("error in postAccount")
		return nil, err
	}
	return &pb.PostAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		}}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	account, err := s.service.GetAccount(ctx, r.Id)

	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{Account: &pb.Account{
		Id:   account.ID,
		Name: account.Name,
	}}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	accountData, err := s.service.GetAccounts(ctx, r.Skip, r.Take)

	if err != nil {
		return nil, err
	}

	accounts := []*pb.Account{}

	for _, account := range accountData {
		accounts = append(accounts,
			&pb.Account{
				Id:   account.ID,
				Name: account.Name,
			})
	}

	return &pb.GetAccountsResponse{Accounts: accounts}, nil
}
