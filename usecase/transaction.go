package usecase

import (
	"time"

	"github.com/renatospaka/code-bank/domain"
	"github.com/renatospaka/code-bank/dto"
)

type UseCaseTransaction struct {
  TransactionRepository domain.TransactionRepository
}

func NewUseCaseTransaction(transactionRepository domain.TransactionRepository) UseCaseTransaction {
  return UseCaseTransaction{TransactionRepository: transactionRepository}
}

func (u UseCaseTransaction) ProcessTransaction(transactionDto dto.Transaction) (domain.Transaction, error) {
  creditCard := u.hydrateCreditCard(transactionDto)
  ccBalanceAndLimit, err := u.TransactionRepository.GetCreditCard(creditCard)
  if err != nil {
    return domain.Transaction{}, err
  }

  creditCard.ID = ccBalanceAndLimit.ID
  creditCard.Limit = ccBalanceAndLimit.Limit
  creditCard.Balance = ccBalanceAndLimit.Balance

  t := u.newTransaction(transactionDto, creditCard)
  t.ProcessAndValidate(&creditCard)

  err = u.TransactionRepository.SaveTransaction(*t, creditCard)
  if err != nil {
    return domain.Transaction{}, err
  }
  return *t, nil
}

func (u UseCaseTransaction) hydrateCreditCard(transactionDto dto.Transaction) domain.CreditCard {
  creditCard := domain.NewCreditCard()
  creditCard.Name = transactionDto.Name
  creditCard.Number = transactionDto.Number
  creditCard.ExpirationMonth = transactionDto.ExpirationMonth
  creditCard.ExpirationYear = transactionDto.ExpirationYear
  creditCard.CVV = transactionDto.CVV
  return *creditCard
}

func (u UseCaseTransaction) newTransaction(transactionDto dto.Transaction, cc domain.CreditCard) *domain.Transaction {
  t := domain.NewTransaction()
  t.ID = cc.ID
  t.Amount = transactionDto.Amount
  t.Store = transactionDto.Store
  t.Description = transactionDto.Description
  t.CreatedAt = time.Now()
  return t
}
