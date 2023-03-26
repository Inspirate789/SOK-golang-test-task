package usecase

import (
	"github.com/Inspirate789/SOK-golang-test-task/internal/models"
	userRepo "github.com/Inspirate789/SOK-golang-test-task/internal/users/repository"
)

type userUseCaseImpl struct {
	userRepo userRepo.Repository
}

func NewUseCase(ur userRepo.Repository) UseCase {
	return &userUseCaseImpl{
		userRepo: ur,
	}
}

func (uuc *userUseCaseImpl) CreateUser(name string) (int, error) {
	return uuc.userRepo.CreateUser(name)
}

func (uuc *userUseCaseImpl) DeleteUser(id int) error {
	return uuc.userRepo.DeleteUser(id)
}

func (uuc *userUseCaseImpl) GetUserBalance(id int) (int, error) {
	user, err := uuc.userRepo.GetUserByID(id)
	if err != nil {
		return 0, models.ErrDatabaseWrap(err)
	}

	return user.Balance, nil
}
