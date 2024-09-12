package repositories

import (
	"context"
	"log"

	"github.com/guigateixeira/general-auth/internal/database"
)

type UserRepository struct {
	db *database.Queries
}

func New(db *database.Queries) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(context context.Context, email string, password string) (string, error) {
	user, err := r.db.CreateUser(context, database.CreateUserParams{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return "", err
	}
	return user.ID.String(), nil
}
