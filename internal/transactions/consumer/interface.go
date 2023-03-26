package consumer

import (
	"context"
)

type TransactionDTO struct {
	UserId      int `form:"user_id" json:"user_id"`
	BalanceDiff int `form:"balance_diff" json:"balance_diff"`
}

type TransactionConsumer interface {
	Consume(ctx context.Context) (TransactionDTO, error)
	Close() error
}
