package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/renatospaka/code-bank/domain"
	"github.com/renatospaka/code-bank/infrastructure/grpc/server"
	"github.com/renatospaka/code-bank/infrastructure/kafka"
	"github.com/renatospaka/code-bank/infrastructure/repository"
	"github.com/renatospaka/code-bank/usecase"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Environment variables not set")
		return
	}
	
	log.Println("Environment are set")
}

func main() {
	fmt.Println("Hello Code Bank")

	db := setupDb()
	defer db.Close()

	//setupCreditCard(db)

	producer := setubKafkaProducer()
	processTransactionUseCase := setupTransactionUseCase(db, producer)
	serveGrpc(processTransactionUseCase)
}

func setupDb() *sql.DB {
	//connection := fmt.Sprintf("port=%s user=%s password=%s dbname=%s sslmode=disable", 
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		os.Getenv("host"),	
		os.Getenv("port"),
		os.Getenv("user"),
		os.Getenv("password"),
		os.Getenv("dbname"),
	)

	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Postgress Server status: %s", db.Ping())
	return db
}

func setupTransactionUseCase(db *sql.DB, producer kafka.KafkaProducer) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDb(db)
	useCase := usecase.NewUseCaseTransaction(transactionRepository)
	useCase.KafkaProducer = producer
	return useCase
}

func setubKafkaProducer() kafka.KafkaProducer {
	producer := kafka.NewKafkaProducer()
	producer.SetupProducer(os.Getenv("KafkaBootstrapServers"))
	log.Println("Kafka Producer running")
	return producer
}

func serveGrpc(processTransactionUseCase usecase.UseCaseTransaction) {
	grpcServer := server.NewGRPCServer()
	grpcServer.ProcessTransactionUseCase = processTransactionUseCase
	log.Println("gRPC Server running")
	grpcServer.Serve()
}

func setupCreditCard(db *sql.DB) {
	cc := domain.NewCreditCard()
	cc.Name = "Eu"
	cc.Number = "1234"
	cc.ExpirationYear = 2025
	cc.ExpirationMonth = 1
	cc.CVV = 321
	cc.Limit = 10000
	cc.Balance = 0

	repo := repository.NewTransactionRepositoryDb(db)
	err := repo.CreateCreditCard(*cc)
	if err != nil {
		log.Printf("Credit Card creation error: %s", err)
	}
	
	log.Println("Credit Card created")
}
