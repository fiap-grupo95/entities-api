package interfaces

import (
	"mecanica_xpto/internal/domain/entities"
)

type ICustomerRepository interface {
	GetByID(id uint) (*entities.Customer, error)
	GetByDocument(CpfCnpj string) (*entities.Customer, error)
	Create(customer *entities.Customer) error
	Update(customer *entities.Customer) error
	Delete(id uint) error
	List() ([]entities.Customer, error)
}
