package repository

import (
	"context"
	"github.com/Inspirate789/SOK-golang-test-task/internal/transactions/consumer"
	"github.com/Inspirate789/SOK-golang-test-task/pkg/dbutils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type pgRepository struct {
	db *sqlx.DB
}

func NewPgRepository(db *sqlx.DB) Repository {
	return &pgRepository{db: db}
}

func applyTransactionTx(ctx context.Context, tx sqlx.ExtContext, transaction consumer.TransactionDTO) error {
	args := map[string]interface{}{
		"user_id":      transaction.UserId,
		"balance_diff": transaction.BalanceDiff,
	}

	_, err := dbutils.NamedExec(context.Background(), tx, updateUserQuery, args)
	if err != nil {
		return errors.Wrap(err, "database error (table: transactions, method: applyTransactionTx)")
	}

	_, err = dbutils.NamedExec(context.Background(), tx, insertTransactionQuery, args)
	if err != nil {
		return errors.Wrap(err, "database error (table: transactions, method: applyTransactionTx)")
	}

	return nil
}

func (tr *pgRepository) PerformTransaction(transaction consumer.TransactionDTO) error {
	err := dbutils.RunTx(context.Background(), tr.db, func(tx *sqlx.Tx) error {
		err := applyTransactionTx(context.Background(), tx, transaction)
		return err
	})

	return err
}
