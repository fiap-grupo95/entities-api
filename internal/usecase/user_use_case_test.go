package usecase

import (
	"errors"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeUserRepo struct {
	usersByID    map[uint]*entities.User
	usersByEmail map[string]*entities.User
	list         []entities.User
	createErr    error
	updateErr    error
	deleteErr    error
	getIDErr     error
	getEmailErr  error
}

func (f *fakeUserRepo) GetByID(id uint) (*entities.User, error) {
	if f.getIDErr != nil {
		return nil, f.getIDErr
	}
	if user, ok := f.usersByID[id]; ok {
		clone := *user
		return &clone, nil
	}
	return nil, nil
}

func (f *fakeUserRepo) GetByEmail(email string) (*entities.User, error) {
	if f.getEmailErr != nil {
		return nil, f.getEmailErr
	}
	if user, ok := f.usersByEmail[email]; ok {
		clone := *user
		return &clone, nil
	}
	return nil, nil
}

func (f *fakeUserRepo) Create(user *entities.User) error {
	if f.createErr != nil {
		return f.createErr
	}
	if user.ID == 0 {
		user.ID = 1
	}
	f.usersByID[user.ID] = user
	f.usersByEmail[user.Email] = user
	f.list = append(f.list, *user)
	return nil
}

func (f *fakeUserRepo) Update(user *entities.User) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	f.usersByID[user.ID] = user
	f.usersByEmail[user.Email] = user
	return nil
}

func (f *fakeUserRepo) Delete(id uint) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	delete(f.usersByID, id)
	return nil
}

func (f *fakeUserRepo) List() ([]entities.User, error) {
	return f.list, nil
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{usersByID: map[uint]*entities.User{}, usersByEmail: map[string]*entities.User{}}
}

func TestUserUseCase_GetByID(t *testing.T) {
	repo := newFakeUserRepo()
	repo.usersByID[1] = &entities.User{ID: 1, Email: "a@b.com"}
	uc := NewUserUseCase(repo)

	user, err := uc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), user.ID)

	_, err = uc.GetByID(99)
	assert.ErrorIs(t, err, ErrUserNotFound)

	repo.getIDErr = errors.New("db err")
	_, err = uc.GetByID(1)
	assert.EqualError(t, err, "db err")
}

func TestUserUseCase_Create(t *testing.T) {
	repo := newFakeUserRepo()
	uc := NewUserUseCase(repo)

	created, err := uc.Create(&entities.User{Email: "new@xpto.com", Password: "123456", UserType: valueobject.ParseUserType("admin")})
	assert.NoError(t, err)
	assert.NotEmpty(t, created.Password)
	assert.NotEqual(t, "123456", created.Password)

	repo.usersByEmail["exists@xpto.com"] = &entities.User{ID: 10, Email: "exists@xpto.com"}
	_, err = uc.Create(&entities.User{Email: "exists@xpto.com", Password: "123456"})
	assert.ErrorIs(t, err, ErrUserAlreadyExists)

	repo.createErr = errors.New("insert err")
	_, err = uc.Create(&entities.User{Email: "err@xpto.com", Password: "123456"})
	assert.EqualError(t, err, "insert err")

	repo.createErr = nil
	repo.getEmailErr = errors.New("lookup err")
	createdWithDefaultType, err := uc.Create(&entities.User{Email: "default@xpto.com", Password: "123456"})
	assert.NoError(t, err)
	assert.Equal(t, valueobject.ParseUserType("customer"), createdWithDefaultType.UserType)
}

func TestUserUseCase_Update_Delete_List(t *testing.T) {
	repo := newFakeUserRepo()
	repo.usersByID[1] = &entities.User{ID: 1, Email: "old@xpto.com", Password: "old", UserType: valueobject.ParseUserType("customer")}
	repo.list = []entities.User{{ID: 1, Email: "old@xpto.com"}}
	uc := NewUserUseCase(repo)

	err := uc.Update(1, &entities.User{Email: "new@xpto.com", Password: "123456", UserType: valueobject.ParseUserType("admin")})
	assert.NoError(t, err)
	assert.Equal(t, "new@xpto.com", repo.usersByID[1].Email)
	assert.NotEqual(t, "123456", repo.usersByID[1].Password)

	err = uc.Update(99, &entities.User{Email: "none@xpto.com"})
	assert.ErrorIs(t, err, ErrUserNotFound)

	repo.getIDErr = errors.New("get err")
	err = uc.Update(1, &entities.User{Email: "x@y.com"})
	assert.EqualError(t, err, "get err")
	repo.getIDErr = nil

	repo.updateErr = errors.New("update err")
	err = uc.Update(1, &entities.User{Email: "fail@xpto.com"})
	assert.EqualError(t, err, "update err")
	repo.updateErr = nil

	err = uc.Delete(1)
	assert.NoError(t, err)

	err = uc.Delete(77)
	assert.ErrorIs(t, err, ErrUserNotFound)

	repo.getIDErr = errors.New("get err")
	err = uc.Delete(1)
	assert.EqualError(t, err, "get err")
	repo.getIDErr = nil

	repo.deleteErr = errors.New("delete err")
	repo.usersByID[2] = &entities.User{ID: 2, Email: "x@y.com"}
	err = uc.Delete(2)
	assert.EqualError(t, err, "delete err")

	list, err := uc.List()
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}
