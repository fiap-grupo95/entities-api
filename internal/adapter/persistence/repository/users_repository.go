package repository

import (
	"mecanica_xpto/internal/domain/entities"
	dto "mecanica_xpto/internal/infrastructure/database/model"
	"mecanica_xpto/internal/usecase/interfaces"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(User *entities.User) error {
	userDto := dto.FromDomainUser(User)
	return r.db.Create(userDto).Error
}

func (r *UserRepository) GetByID(id uint) (*entities.User, error) {
	var User dto.UserModel
	err := r.db.Preload("User").First(&User, id).Error
	if err != nil {
		return nil, err
	}
	return (&User).ToDomain(), nil
}

func (r *UserRepository) GetByEmail(email string) (*entities.User, error) {
	var userDto dto.UserModel
	err := r.db.Where("email = ?", email).First(&userDto).Error
	if err != nil {
		return nil, err
	}
	return (&userDto).ToDomain(), nil
}

func (r *UserRepository) Update(User *entities.User) error {
	userDto := dto.FromDomainUser(User)
	return r.db.Save(userDto).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&dto.UserModel{}, id).Error
}

func (r *UserRepository) List() ([]entities.User, error) {
	var Users []dto.UserModel
	err := r.db.Preload("User").Find(&Users).Error
	if err != nil {
		return nil, err
	}
	var result []entities.User
	for _, u := range Users {
		result = append(result, *u.ToDomain())
	}
	return result, err
}
