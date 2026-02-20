package dto

import (
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	ID        uint                 `gorm:"primaryKey"`
	Email     string               `gorm:"size:100;not null;unique" json:"email"`
	Password  string               `gorm:"size:255;not null" json:"password"`
	UserType  valueobject.UserType `gorm:"not null"`
	CreatedAt time.Time            `gorm:"autoCreateTime"`
	UpdatedAt time.Time            `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt       `gorm:"index"`
	Customer  *CustomerModel       `gorm:"foreignKey:UserID;references:ID"`
}

func (m *UserModel) TableName() string {
	return "user"
}

func (m *UserModel) ToDomain() *entities.User {
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}
	var customer *entities.Customer
	if m.Customer != nil {
		c := m.Customer.ToDomain()
		customer = c
	}
	return &entities.User{
		ID:        m.ID,
		Email:     m.Email,
		Password:  m.Password,
		UserType:  m.UserType,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
		Customer:  customer,
	}
}

func FromDomainUser(user *entities.User) *UserModel {
	var deletedAt gorm.DeletedAt
	if user.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *user.DeletedAt, Valid: true}
	}
	var customer *CustomerModel
	if user.Customer != nil {
		c := FromDomainCustomer(user.Customer)
		customer = c
	}
	return &UserModel{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: deletedAt,
		Customer:  customer,
	}
}
