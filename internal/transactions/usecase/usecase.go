package usecase

import (
	"github.com/Inspirate789/SOK-golang-test-task/internal/models"
	"github.com/Inspirate789/SOK-golang-test-task/internal/transactions/consumer"
	transactionRepo "github.com/Inspirate789/SOK-golang-test-task/internal/transactions/repository"
	userRepo "github.com/Inspirate789/SOK-golang-test-task/internal/users/repository"
)

type transactionUseCaseImpl struct {
	userRepo        userRepo.Repository
	transactionRepo transactionRepo.Repository
}

func NewUseCase(ur userRepo.Repository, tr transactionRepo.Repository) UseCase {
	return &transactionUseCaseImpl{
		userRepo:        ur,
		transactionRepo: tr,
	}
}

func (tuc *transactionUseCaseImpl) PerformTransaction(transaction consumer.TransactionDTO) error {
	user, err := tuc.userRepo.GetUserByID(transaction.UserId)
	if err != nil {
		return models.ErrDatabaseWrap(err)
	}

	if user.Balance-transaction.BalanceDiff < 0 {
		return models.ErrNotEnoughMoney()
	}

	return tuc.transactionRepo.PerformTransaction(transaction)
}
