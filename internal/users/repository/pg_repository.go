package repository

import (
	"context"
	"github.com/Inspirate789/SOK-golang-test-task/internal/models"
	"github.com/Inspirate789/SOK-golang-test-task/pkg/dbutils"
	"github.com/jmoiron/sqlx"
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
	res, err := dbutils.NamedExec(context.Background(), ur.db, insertUserQuery, args)
	if err != nil {
		return 0, errors.Wrap(err, "database error (table: users, method: CreateUser)")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "database error (table: users, method: CreateUser)")
	}

	return int(id), nil
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
