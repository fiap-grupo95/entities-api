package database

import (
	"fmt"
	Model "mecanica_xpto/internal/infrastructure/database/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&Model.PartsSupplyModel{},
		&Model.ServiceModel{},
		&Model.VehicleModel{},
		&Model.CustomerModel{},
		&Model.UserModel{},
		&Model.UserTypeModel{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	fmt.Println("Database migrated successfully")
}
