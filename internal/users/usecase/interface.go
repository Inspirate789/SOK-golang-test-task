package usecase

type UseCase interface {
	CreateUser(name string) (int, error)
	DeleteUser(id int) error
	GetUserBalance(id int) (int, error)
}
