package usecase

import (
	"errors"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	"mecanica_xpto/internal/usecase/interfaces"
	"mecanica_xpto/pkg"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type IUserUseCase interface {
	GetByID(id uint) (*entities.User, error)
	Create(user *entities.User) (*entities.User, error)
	Update(id uint, user *entities.User) error
	Delete(id uint) error
	List() ([]entities.User, error)
}

type UserUseCase struct {
	repo interfaces.IUserRepository
}

func NewUserUseCase(repo interfaces.IUserRepository) IUserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) GetByID(id uint) (*entities.User, error) {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil || user.ID == 0 {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (u *UserUseCase) Create(user *entities.User) (*entities.User, error) {
	existing, err := u.repo.GetByEmail(user.Email)
	if err == nil && existing != nil && existing.ID != 0 {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, _ := pkg.HashPassword(user.Password)
	user.Password = hashedPassword

	if user.UserType == "" {
		user.UserType = valueobject.ParseUserType("customer")
	}

	if err := u.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) Update(id uint, user *entities.User) error {
	existing, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil || existing.ID == 0 {
		return ErrUserNotFound
	}

	if user.Email != "" {
		existing.Email = user.Email
	}

	if user.Password != "" {
		hashedPassword, _ := pkg.HashPassword(user.Password)
		existing.Password = hashedPassword
	}

	if user.UserType != "" {
		existing.UserType = user.UserType
	}

	return u.repo.Update(existing)
}

func (u *UserUseCase) Delete(id uint) error {
	existing, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil || existing.ID == 0 {
		return ErrUserNotFound
	}
	return u.repo.Delete(id)
}

func (u *UserUseCase) List() ([]entities.User, error) {
	return u.repo.List()
}
