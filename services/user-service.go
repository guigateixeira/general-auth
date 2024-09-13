package services

import (
	"context"
	"errors"
	"log"

	"github.com/guigateixeira/general-auth/model"
	"github.com/guigateixeira/general-auth/repositories"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func New(userRepository *repositories.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) CreateUser(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if user != nil {
		return "", errors.New("Email is already taken")
	}
	userID, err := s.userRepository.CreateUser(ctx, email, password)
	if err != nil {
		log.Printf("Service layer error creating user: %v", err)
		return "", err
	}
	return userID, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("Service layer error getting user by email: %v", err)
		return nil, err
	}
	return user, nil
}
