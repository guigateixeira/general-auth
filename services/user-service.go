package services

import (
	"context"
	"log"
	"net/http"

	"github.com/guigateixeira/general-auth/errors"
	"github.com/guigateixeira/general-auth/model"
	"github.com/guigateixeira/general-auth/repositories"
	"github.com/guigateixeira/general-auth/util"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func New(userRepository *repositories.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) CreateUser(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("Service layer error getting user by email: %v", err)
		return "", err
	}
	if user != nil {
		return "", errors.NewBaseError("Email is already taken", 400)
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		log.Printf("Service layer error hashing password: %v", err)
		return "", err
	}
	userID, err := s.userRepository.CreateUser(ctx, email, hashedPassword)
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

func (s *UserService) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.NewBaseError("Error fetching user", http.StatusInternalServerError)
	}
	if user == nil {
		return "", errors.NewBaseError("Invalid email or password", http.StatusUnauthorized)
	}

	if !util.VerifyPassword(user.Password, password) {
		return "", errors.NewBaseError("Invalid email or password", http.StatusUnauthorized)
	}

	token, err := util.GenerateJWTToken(user.Id.String())
	if err != nil {
		return "", errors.NewBaseError("Error generating token", http.StatusInternalServerError)
	}

	return token, nil
}
