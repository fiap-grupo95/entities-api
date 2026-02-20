package usecase

import (
	"context"
	"errors"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	"mecanica_xpto/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

type covCustomerRepo struct {
	getByID       func(id uint) (*entities.Customer, error)
	getByDocument func(doc string) (*entities.Customer, error)
	create        func(c *entities.Customer) error
	update        func(c *entities.Customer) error
	delete        func(id uint) error
	list          func() ([]entities.Customer, error)
}

func (r covCustomerRepo) GetByID(id uint) (*entities.Customer, error) { return r.getByID(id) }
func (r covCustomerRepo) GetByDocument(doc string) (*entities.Customer, error) {
	return r.getByDocument(doc)
}
func (r covCustomerRepo) Create(c *entities.Customer) error  { return r.create(c) }
func (r covCustomerRepo) Update(c *entities.Customer) error  { return r.update(c) }
func (r covCustomerRepo) Delete(id uint) error               { return r.delete(id) }
func (r covCustomerRepo) List() ([]entities.Customer, error) { return r.list() }

type covVehicleRepo struct {
	findByID func(id uint) (*entities.Vehicle, error)
	update   func(v entities.Vehicle) error
	delete   func(id uint) error
}

func (r covVehicleRepo) FindAll() ([]entities.Vehicle, error)        { return []entities.Vehicle{}, nil }
func (r covVehicleRepo) FindByID(id uint) (*entities.Vehicle, error) { return r.findByID(id) }
func (r covVehicleRepo) FindByPlate(plate valueobject.Plate) (*entities.Vehicle, error) {
	if plate.String() == "ABC1D23" {
		return &entities.Vehicle{ID: 1, Plate: plate}, nil
	}
	if plate.String() == "ERR1234" {
		return nil, errors.New("plate err")
	}
	if plate.String() == "ZZZ9999" {
		return &entities.Vehicle{}, nil
	}
	return nil, nil
}
func (r covVehicleRepo) FindByCustomerID(customerID uint) ([]entities.Vehicle, error) {
	if customerID == 1 {
		return nil, nil
	}
	if customerID == 2 {
		return nil, errors.New("repo err")
	}
	return []entities.Vehicle{{ID: 1}}, nil
}
func (r covVehicleRepo) Create(vehicle entities.Vehicle) (*entities.Vehicle, error) {
	if vehicle.Plate.String() == "AAA0001" {
		return nil, errors.New("create err")
	}
	vehicle.ID = 9
	return &vehicle, nil
}
func (r covVehicleRepo) Update(v entities.Vehicle) error { return r.update(v) }
func (r covVehicleRepo) Delete(id uint) error            { return r.delete(id) }

type covPartsRepo struct{}

func (covPartsRepo) Create(ctx context.Context, ps *entities.PartsSupply) (entities.PartsSupply, error) {
	return *ps, nil
}
func (covPartsRepo) GetByID(ctx context.Context, id uint) (entities.PartsSupply, error) {
	if id == 1 {
		return entities.PartsSupply{ID: 1}, nil
	}
	return entities.PartsSupply{}, nil
}
func (covPartsRepo) GetByName(ctx context.Context, name string) (entities.PartsSupply, error) {
	return entities.PartsSupply{}, errors.New("lookup err")
}
func (covPartsRepo) Update(ctx context.Context, ps *entities.PartsSupply) error { return nil }
func (covPartsRepo) Delete(ctx context.Context, id uint) error                  { return nil }
func (covPartsRepo) List(ctx context.Context) ([]entities.PartsSupply, error) {
	return []entities.PartsSupply{}, nil
}
func (covPartsRepo) GetByServiceOrderID(ctx context.Context, serviceOrderID uint) ([]entities.PartsSupply, error) {
	if serviceOrderID == 1 {
		return []entities.PartsSupply{{ID: 1}}, nil
	}
	if serviceOrderID == 2 {
		return nil, errors.New("so err")
	}
	return []entities.PartsSupply{}, nil
}

type covToken struct {
	value string
	err   error
}

func (t covToken) GenerateToken(subject string) (string, error) {
	return t.value, t.err
}

func TestCoverage_Customer_Branches(t *testing.T) {
	validDoc, _ := valueobject.NewCpfCnpj("52998224725")
	uc := NewCustomerUseCase(covCustomerRepo{
		getByID: func(id uint) (*entities.Customer, error) {
			if id == 1 {
				return &entities.Customer{ID: 1}, nil
			}
			if id == 2 {
				return nil, nil
			}
			return nil, errors.New("db")
		},
		getByDocument: func(doc string) (*entities.Customer, error) {
			if doc == "ok" {
				return &entities.Customer{ID: 1}, nil
			}
			if doc == "none" {
				return nil, nil
			}
			return nil, errors.New("db")
		},
		create: func(c *entities.Customer) error { return nil },
		update: func(c *entities.Customer) error { return nil },
		delete: func(id uint) error { return nil },
		list:   func() ([]entities.Customer, error) { return []entities.Customer{{ID: 1}}, nil },
	}, nil)

	_, err := uc.GetById(1)
	assert.NoError(t, err)
	_, err = uc.GetById(2)
	assert.ErrorIs(t, err, ErrCustomerNotFound)
	_, err = uc.GetById(3)
	assert.ErrorIs(t, err, ErrGeneric)

	_, err = uc.GetByDocument("ok")
	assert.NoError(t, err)
	_, err = uc.GetByDocument("none")
	assert.ErrorIs(t, err, ErrCustomerNotFound)
	_, err = uc.GetByDocument("err")
	assert.ErrorIs(t, err, ErrGeneric)

	_, err = uc.CreateCustomer(&entities.Customer{CpfCnpj: valueobject.CpfCnpj("123"), Email: "x@x.com"})
	assert.ErrorIs(t, err, ErrInvalidDocumentFormat)

	created, err := uc.CreateCustomer(&entities.Customer{CpfCnpj: validDoc, Email: "x@x.com"})
	assert.NoError(t, err)
	assert.NotNil(t, created.User)
	assert.NotEmpty(t, created.User.Password)

	list, err := uc.ListCustomer()
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	errUC := NewCustomerUseCase(covCustomerRepo{
		getByID:       func(id uint) (*entities.Customer, error) { return &entities.Customer{ID: 1}, nil },
		getByDocument: func(doc string) (*entities.Customer, error) { return &entities.Customer{ID: 1}, nil },
		create:        func(c *entities.Customer) error { return errors.New("create customer err") },
		update:        func(c *entities.Customer) error { return nil },
		delete:        func(id uint) error { return nil },
		list:          func() ([]entities.Customer, error) { return nil, errors.New("list err") },
	}, nil)

	_, err = errUC.CreateCustomer(&entities.Customer{CpfCnpj: validDoc, Email: "err@x.com"})
	assert.EqualError(t, err, "create customer err")
	_, err = errUC.ListCustomer()
	assert.EqualError(t, err, "list err")
}

func TestCoverage_Vehicle_Branches(t *testing.T) {
	repo := covVehicleRepo{
		findByID: func(id uint) (*entities.Vehicle, error) {
			if id == 1 {
				return &entities.Vehicle{ID: 1, Plate: valueobject.ParsePlate("ABC1D23")}, nil
			}
			if id == 2 {
				return &entities.Vehicle{}, nil
			}
			if id == 4 {
				return &entities.Vehicle{ID: 4}, nil
			}
			if id == 3 {
				return nil, errors.New("find err")
			}
			return nil, nil
		},
		update: func(v entities.Vehicle) error {
			if v.Model == "fail" {
				return errors.New("upd err")
			}
			return nil
		},
		delete: func(id uint) error {
			if id == 4 {
				return errors.New("del err")
			}
			return nil
		},
	}
	svc := NewVehicleService(repo)

	_, err := svc.GetVehicleByPlate("INVALID")
	assert.ErrorIs(t, err, ErrInvalidPlateFormat)
	_, err = svc.GetVehicleByPlate("ERR1234")
	assert.EqualError(t, err, "plate err")
	_, err = svc.GetVehicleByPlate("ZZZ9999")
	assert.ErrorIs(t, err, ErrVehicleNotFound)

	_, err = svc.GetVehiclesByCustomerID(1)
	assert.NoError(t, err)
	_, err = svc.GetVehiclesByCustomerID(2)
	assert.EqualError(t, err, "repo err")

	_, err = svc.CreateVehicle(entities.Vehicle{Plate: valueobject.ParsePlate("ABC1D23")})
	assert.ErrorIs(t, err, ErrVehicleAlreadyExists)
	_, err = svc.CreateVehicle(entities.Vehicle{Plate: valueobject.ParsePlate("AAA0001")})
	assert.EqualError(t, err, "create err")

	msg, err := svc.UpdateVehicle(entities.Vehicle{ID: 3, Plate: valueobject.ParsePlate("ABC1D23")})
	assert.Equal(t, MessageErrorSearch, msg)
	assert.EqualError(t, err, "find err")
	msg, err = svc.UpdateVehicle(entities.Vehicle{ID: 2, Plate: valueobject.ParsePlate("ABC1D23")})
	assert.ErrorIs(t, err, ErrVehicleNotFound)
	assert.Equal(t, MessageVehicleNotFound, msg)
	msg, err = svc.UpdateVehicle(entities.Vehicle{ID: 1, Plate: valueobject.ParsePlate("ABC1D23"), Model: "fail"})
	assert.Equal(t, MessageErrorUpdatingVehicle, msg)
	assert.EqualError(t, err, "upd err")

	_, err = svc.UpdateVehiclePartial(3, map[string]interface{}{})
	assert.EqualError(t, err, "find err")
	_, err = svc.UpdateVehiclePartial(5, map[string]interface{}{})
	assert.ErrorIs(t, err, ErrVehicleNotFound)
	_, err = svc.UpdateVehiclePartial(1, map[string]interface{}{"plate": "BAD"})
	assert.ErrorIs(t, err, ErrInvalidPlateFormat)
	_, err = svc.UpdateVehiclePartial(1, map[string]interface{}{"model": "fail"})
	assert.EqualError(t, err, "upd err")
	_, err = svc.UpdateVehiclePartial(1, map[string]interface{}{"model": "ok", "customer_id": float64(7)})
	assert.NoError(t, err)
	_, err = svc.UpdateVehiclePartial(1, map[string]interface{}{"brand": "B", "year": "2026", "plate": "ABC1D23"})
	assert.NoError(t, err)

	err = svc.DeleteVehicle(0)
	assert.ErrorIs(t, err, ErrInvalidID)
	err = svc.DeleteVehicle(3)
	assert.EqualError(t, err, "find err")
	err = svc.DeleteVehicle(2)
	assert.ErrorIs(t, err, ErrVehicleNotFound)
	err = svc.DeleteVehicle(4)
	assert.EqualError(t, err, "del err")
}

func TestCoverage_PartsSupply_And_Auth(t *testing.T) {
	ps := NewPartsSupplyUseCase(covPartsRepo{})
	_, err := ps.GetPartsSupplyByServiceOrderID(context.Background(), 1)
	assert.NoError(t, err)
	_, err = ps.GetPartsSupplyByServiceOrderID(context.Background(), 2)
	assert.EqualError(t, err, "so err")
	_, err = ps.GetPartsSupplyByServiceOrderID(context.Background(), 3)
	assert.ErrorIs(t, err, ErrPartsSupplyNotFound)

	hashed, _ := pkg.HashPassword("123456")
	userRepo := &fakeUserRepo{usersByID: map[uint]*entities.User{}, usersByEmail: map[string]*entities.User{"ok@xpto.com": {ID: 1, Email: "ok@xpto.com", Password: hashed}}}
	auth := NewAuthUseCase(covToken{value: "token"}, userRepo)

	token, appErr := auth.Login("ok@xpto.com", "123456")
	assert.Nil(t, appErr)
	assert.Equal(t, "token", token)

	_, appErr = auth.Login("ok@xpto.com", "wrong")
	assert.NotNil(t, appErr)
	notFoundRepo := &fakeUserRepo{usersByID: map[uint]*entities.User{}, usersByEmail: map[string]*entities.User{}, getEmailErr: errors.New("not found")}
	authNotFound := NewAuthUseCase(covToken{value: "token"}, notFoundRepo)
	_, appErr = authNotFound.Login("notfound@xpto.com", "123456")
	assert.NotNil(t, appErr)

	authFail := NewAuthUseCase(covToken{err: errors.New("token err")}, userRepo)
	_, appErr = authFail.Login("ok@xpto.com", "123456")
	assert.NotNil(t, appErr)
}
