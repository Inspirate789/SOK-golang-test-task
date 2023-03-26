package repository

import (
	"context"
	"github.com/Inspirate789/SOK-golang-test-task/internal/models"
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

func (ur *pgRepository) CreateUser(name string) (int, error) {
	args := map[string]interface{}{
		"name":    name,
		"balance": 0,
	}
	var id int
	err := dbutils.NamedGet(context.Background(), ur.db, &id, insertUserQuery, args)
	if err != nil {
		return 0, errors.Wrap(err, "database error (table: users, method: CreateUser)")
	}

	return id, nil
}

func (ur *pgRepository) DeleteUser(id int) error {
	args := map[string]interface{}{
		"id": id,
	}
	_, err := dbutils.NamedExec(context.Background(), ur.db, deleteUserQuery, args)
	if err != nil {
		return errors.Wrap(err, "database error (table: users, method: DeleteUser)")
	}

	return nil
}

func (ur *pgRepository) GetUserByID(id int) (*models.User, error) {
	args := map[string]interface{}{
		"id": id,
	}
	var user models.User
	err := dbutils.NamedGet(context.Background(), ur.db, &user, selectUserQuery, args)
	if err != nil {
		return nil, errors.Wrap(err, "database error (table: users, method: GetUserByID)")
	}

	return &user, nil
}
