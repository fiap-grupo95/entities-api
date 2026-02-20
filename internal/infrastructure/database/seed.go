package database

import (
	"fmt"
	dto "mecanica_xpto/internal/infrastructure/database/model"

	"gorm.io/gorm"
)

const (
	defaultPassword = "jKHN1SmGKuGKyiXhbnaOZg==.0/rdilUJyR5raIXVdOCaX8szZCEUzIpIhYTQIMaLwc8="
	defaultUserType = "admin"
)

func Seed(db *gorm.DB) {
	// Seed user_dtos
	var countUsers int64
	db.Model(&dto.UserModel{}).Count(&countUsers)
	if countUsers == 0 {
		users := []dto.UserModel{
			{
				Email:    "admin@xpto.com",
				Password: defaultPassword,
				UserType: defaultUserType,
			},
			{
				Email:    "joao@xpto.com",
				Password: defaultPassword,
				UserType: defaultUserType,
			},
			{
				Email:    "joana@xpto.com",
				Password: defaultPassword,
				UserType: defaultUserType,
			},
		}

		if err := db.Create(&users).Error; err != nil {
			fmt.Println("Erro ao criar usuários:", err)
			return
		}
		fmt.Println("Seeded users successfully")
	} else {
		fmt.Println("Users already seeded")
	}

}
