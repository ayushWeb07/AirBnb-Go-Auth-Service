package repositories

import (
	"database/sql"
	"fmt"
)

type UserRepositoryInterface interface {
	CreateUser() error
}

type UserRepository struct {
	db *sql.DB
}

func (ur *UserRepository) CreateUser() error {
	fmt.Println("Create user repository called...")
	return nil
}

func NewUserRepository() UserRepositoryInterface {
	newUserRepository := &UserRepository{
		db: nil,
	}

	return newUserRepository
}
