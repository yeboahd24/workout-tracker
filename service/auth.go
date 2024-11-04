package service

import (
	"context"

	"github.com/yeboahd24/workout-tracker/model"
	"github.com/yeboahd24/workout-tracker/repository"
	"github.com/yeboahd24/workout-tracker/util"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *AuthService) SignUp(ctx context.Context, username, email, password string) error {
	user, err := model.NewUser(username, email, password)
	if err != nil {
		return err
	}
	return s.userRepo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if !user.CheckPassword(password) {
		return "", util.ErrInvalidCredentials
	}

	token, err := util.GenerateJWT(user.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
