package repositories

import (
	"context"
	"database/sql"
	"log"

	"github.com/guigateixeira/general-auth/internal/database"
	"github.com/guigateixeira/general-auth/model"
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

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error getting user by email: %v", err)
		return nil, err
	}

	userModel := model.New(user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	return userModel, nil
}
