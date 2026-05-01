package service

import (
	"context"
	"errors"

	"pgxpostgress/domain"
	"pgxpostgress/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	List(ctx context.Context) ([]domain.User, int, error)
	Create(ctx context.Context, user domain.User) error
	Get(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(
	repo repository.UserRepository,
) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) List(
	ctx context.Context,
) ([]domain.User, int, error) {

	return s.repo.List(ctx)
}

func (s *userService) Create(
	ctx context.Context,
	user domain.User,
) error {

	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Name == "" {
		return errors.New("name is required")
	}
	

	if user.Password == "" {
		return errors.New("password is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return errors.New("could not hash password")
	}

	user.Password = string(hashedPassword)

	return s.repo.Create(ctx, user)
}

func (s *userService) Get(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *userService) Update(
	ctx context.Context,
	user domain.User,
) error {
	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.repo.Delete(ctx, id)
}
