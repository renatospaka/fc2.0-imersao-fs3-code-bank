package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/renatospaka/code-bank/domain"
	"github.com/renatospaka/code-bank/infrastructure/repository"
	"github.com/renatospaka/code-bank/usecase"
)

func main() {
	fmt.Println("Hello Code Bank")

	db := setupDb()
	defer db.Close()

	cc := domain.NewCreditCard()
	cc.Name = "Renato"
	cc.Number = "1234"
	cc.ExpirationYear = 2024
	cc.ExpirationMonth = 7
	cc.CVV = 123
	cc.Limit = 1000
	cc.Balance = 0

	repo := repository.NewTransactionRepositoryDb(db)
	err := repo.CreateCreditCard(*cc)
	if err != nil {
		fmt.Println(err)
	}
}

func setupDb() *sql.DB {
	username := "postgres"
	password := "Postgres2020!"
	database := "codebank"
	port := "15432"
	//host := "postgres-compose"

	//connection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", 
	connection := fmt.Sprintf("port=%s user=%s password=%s dbname=%s sslmode=disable", 
		port,
		username,
		password,
		database)

	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func setupTransactionUseCase(db *sql.DB) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDb(db)
	useCase := usecase.NewUseCaseTransaction(transactionRepository)
	return useCase
}