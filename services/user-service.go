package services

import (
	"context"
	"log"

	"github.com/guigateixeira/general-auth/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func New(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, email string, password string) (string, error) {
	userID, err := s.repo.CreateUser(ctx, email, password)
	if err != nil {
		log.Printf("Service layer error creating user: %v", err)
		return "", err
	}
	return userID, nil
}
