package server

import (
	"log"
	"net"

	"github.com/renatospaka/code-bank/infrastructure/grpc/pb"
	"github.com/renatospaka/code-bank/infrastructure/grpc/service"
	"github.com/renatospaka/code-bank/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	ProcessTransactionUseCase usecase.UseCaseTransaction
}

func NewGRPCServer() GRPCServer {
	return GRPCServer{}
}

func (g GRPCServer) Serve() {
	listener, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("could not listen to port %d", 50052)
	}

	transactionService := service.NewTransactionService()
	transactionService.ProcessTransactionUseCase = g.ProcessTransactionUseCase
	
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer) //modo relfexion por conta do evans (client gRPC)
	pb.RegisterPaymentServiceServer(grpcServer, transactionService)
	grpcServer.Serve(listener)
}