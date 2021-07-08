package repository

import (
	"database/sql"
	"errors"

	"github.com/renatospaka/code-bank/domain"
)

type TransactionRepositoryDb struct {
	db *sql.DB
}

func NewTransactionRepositoryDb(db *sql.DB) *TransactionRepositoryDb {
	return &TransactionRepositoryDb{db: db}
}

func (t *TransactionRepositoryDb) SaveTransaction(transaction domain.Transaction, creditcard domain.CreditCard) error {
	stmt, err := t.db.Prepare(`insert into transactions (id, credit_card_id, amount, status, description, store, created_at)
														values ($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		transaction.ID,
		transaction.CreditCardId,
		transaction.Amount,
		transaction.Status,
		transaction.Description,
		transaction.Store,
		transaction.CreatedAt,
	)
	if err != nil {
		return err
	}

	if transaction.Status == "approved" {
		err = t.updateBalance(creditcard)
		if err != nil {
			return err
		}
	}

	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionRepositoryDb) CreateCreditCard(creditcard domain.CreditCard) error {
	stmt, err := t.db.Prepare(`insert into credit_cards (id, name, number, expiration_month, expiration_year, cvv, balance, balance_limit)
														values ($1, $2, $3, $4, $5, $6, $7, $8)`)
	if err != nil {
		return err
	}
	
	_, err = stmt.Exec(
		creditcard.ID,
		creditcard.Name,
		creditcard.Number,
		creditcard.ExpirationMonth,
		creditcard.ExpirationYear,
		creditcard.CVV,
		creditcard.Balance,
		creditcard.Limit,
	)
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionRepositoryDb) GetCreditCard(creditcard domain.CreditCard) (domain.CreditCard, error) {
	var c domain.CreditCard
	stmt, err := t.db.Prepare(`select id, balance, balance_limit from credit_cards where number = $1`)
	if err != nil {
		return c, err
	}

	if err = stmt.QueryRow(creditcard.Number).Scan(&c.ID, &c.Balance, &c.Limit); err != nil {
		return c, errors.New("credit card number does not exist")
	}

	return c, nil
}

func (t *TransactionRepositoryDb) updateBalance(creditcard domain.CreditCard) error {
	_, err := t.db.Exec("update credit_cards set balance = $1 where id = $2",
							creditcard.Balance, creditcard.ID)
	if err != nil {
		return err
	}
	return nil
}
