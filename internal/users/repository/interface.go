package repository

import "github.com/Inspirate789/SOK-golang-test-task/internal/models"

type Repository interface {
	CreateUser(name string) (int, error)
	DeleteUser(id int) error
	GetUserByID(id int) (*models.User, error)
}
