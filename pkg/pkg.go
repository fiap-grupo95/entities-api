package pkg

import (
	passwordpkg "mecanica_xpto/internal/domain/valueobject/password_utils"
)

func HashPassword(password string) (string, error) {
	return passwordpkg.HashPassword(password)
}
