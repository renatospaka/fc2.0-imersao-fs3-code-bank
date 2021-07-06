package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type TransactionRepository interface {
	SaveTransaction(transaction Transaction, creditcard CreditCard) error
	GetCreditCard(creditcard CreditCard) (CreditCard, error)
	CreateCreditCard(creditcard CreditCard) error
}

type Transaction struct {
	ID           string
	Amount       float64
	Status       string
	Description  string
	Store        string
	CreditCardId string
	CreatedAt    time.Time
}

// * => é um ponteiro para o endereço de memória
func NewTransaction() *Transaction {
	t := &Transaction{}
	t.ID = uuid.NewV4().String()
	t.CreatedAt = time.Now()

	return t
}

//como creditcard foi passado como ponteiro, a variável é alterada diretamente, 
//sem necessidade de um retorno
func (t *Transaction) ProcessAndValidate(creditcard *CreditCard) {
	if t.Amount + creditcard.Balance > creditcard.Limit {
		t.Status = "rejected"
	} else {
		t.Status = "approved"
		creditcard.Balance = creditcard.Balance + t.Amount
	}
}
