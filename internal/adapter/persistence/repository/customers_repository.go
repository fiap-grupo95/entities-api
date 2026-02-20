package repository

import (
	"mecanica_xpto/internal/domain/entities"
	dto "mecanica_xpto/internal/infrastructure/database/model"
	"mecanica_xpto/internal/usecase/interfaces"
	"strings"

	"gorm.io/gorm"
)

// CustomerRepository implements ICustomerRepository interface
type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) interfaces.ICustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(customer *entities.Customer) error {
	customerDto := dto.FromDomainCustomer(customer)
	if err := r.db.Create(customerDto).Error; err != nil {
		return err
	}

	created := customerDto.ToDomain()
	customer.ID = created.ID
	customer.UserID = created.UserID

	if created.User != nil {
		if customer.User == nil {
			customer.User = created.User
		} else {
			customer.User.ID = created.User.ID
			if customer.User.Email == "" {
				customer.User.Email = created.User.Email
			}
			if customer.User.Password == "" {
				customer.User.Password = created.User.Password
			}
			if customer.User.UserType == "" {
				customer.User.UserType = created.User.UserType
			}
		}
	}

	return nil
}

func (r *CustomerRepository) GetByID(id uint) (*entities.Customer, error) {
	var customer dto.CustomerModel
	err := r.db.Preload("User").Preload("Vehicles").First(&customer, id).Error
	if err != nil {
		if strings.EqualFold(err.Error(), gorm.ErrRecordNotFound.Error()) {
			return nil, nil
		}
	}
	return (&customer).ToDomain(), nil
}

func (r *CustomerRepository) GetByDocument(CpfCnpj string) (*entities.Customer, error) {
	var customer dto.CustomerModel
	err := r.db.Preload("User").First(&customer, map[string]interface{}{"cpf_cnpj": CpfCnpj}).Error
	if err != nil {
		if strings.EqualFold(err.Error(), gorm.ErrRecordNotFound.Error()) {
			return nil, nil
		}
	}
	return (&customer).ToDomain(), nil
}

func (r *CustomerRepository) Update(customer *entities.Customer) error {
	customerDto := dto.FromDomainCustomer(customer)
	return r.db.Save(customerDto).Error
}

func (r *CustomerRepository) Delete(id uint) error {
	var customer dto.CustomerModel
	err := r.db.Preload("User").First(&customer, id).Error
	if err != nil {
		return err
	}
	if err := r.db.Delete(&customer.User).Error; err != nil {
		return err
	}

	return r.db.Delete(&customer).Error
}

func (r *CustomerRepository) List() ([]entities.Customer, error) {
	var customers []dto.CustomerModel
	err := r.db.Preload("User").Find(&customers).Error
	if err != nil {
		return nil, err
	}
	var result []entities.Customer
	for _, c := range customers {
		result = append(result, *c.ToDomain())
	}
	return result, err
}
