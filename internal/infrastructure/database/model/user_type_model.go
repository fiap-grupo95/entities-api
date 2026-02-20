package dto

import (
	"mecanica_xpto/internal/domain/valueobject"
)

type UserTypeModel struct {
	ID   uint   `gorm:"primaryKey"`
	Type string `gorm:"size:50;not null"`
	//Users []UserModel `gorm:"foreignKey:UserTypeID"`
}

func (utm *UserTypeModel) TableName() string {
	return "user_type"
}

func (utm *UserTypeModel) ToDomain() valueobject.UserType {
	return valueobject.ParseUserType(utm.Type)
}
