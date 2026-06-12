package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) UpdateProfile(ctx context.Context, id int64, req model.UpdateProfileRequest) error {
	if req.Name == "" || req.Email == "" {
		return errors.New("ad ve e-posta zorunlu")
	}
	return s.repo.UpdateProfile(ctx, id, req)
}

func (s *UserService) UpdatePassword(ctx context.Context, id int64, req model.UpdatePasswordRequest) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Current)); err != nil {
		return errors.New("mevcut şifre yanlış")
	}
	if len(req.New) < 6 {
		return errors.New("yeni şifre en az 6 karakter olmalı")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.New), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, id, string(hash))
}
