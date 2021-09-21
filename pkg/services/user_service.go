package services

import (
	"errors"

	"github.com/casmelad/GlobantPOC/pkg/domain"
	"github.com/casmelad/GlobantPOC/pkg/domain/entities"
	"gopkg.in/go-playground/validator.v9"
)

//application service
type UserServiceInterface interface {
	Create(entities.User) (int, error)
	GetByEmail(string) (entities.User, error)
	GetAll() ([]entities.User, error)
	Update(entities.User) error
	Delete(int) error
}

type UserService struct {
	repository domain.UsersRepositoryInterface
}

func NewUserService(repo domain.UsersRepositoryInterface) *UserService {
	return &UserService{
		repository: repo,
	}
}

func (us *UserService) Create(usr entities.User) (int, error) {

	v := validator.New()
	err := v.Struct(usr)

	if err != nil {
		errorMessage := err.(validator.ValidationErrors)[0]
		return 0, errors.New(errorMessage.Field() + " is not valid")
	}

	dbUser := us.repository.GetByEmail(usr.Email)

	if dbUser.Id <= 0 {
		newId := us.repository.Add(usr)
		return newId, nil
	}

	return 0, errors.New("user already exists")
}

func (us *UserService) GetByEmail(email string) (entities.User, error) {
	dbUser := us.repository.GetByEmail(email)

	if dbUser.Id == 0 {
		return entities.User{}, errors.New("user not found")
	}

	return dbUser, nil

}

func (us *UserService) GetAll() ([]entities.User, error) {

	return us.repository.GetAll(), nil //not used for the moment
}

func (us *UserService) Update(usr entities.User) error {

	v := validator.New()
	err := v.Struct(usr)

	if err != nil {
		errorMessage := err.(validator.ValidationErrors)[0]
		return errors.New(errorMessage.Field() + " is not valid")
	}

	usrToUpdate := us.repository.GetByEmail(usr.Email)

	if usrToUpdate.Id == 0 {
		return errors.New("user not found")
	}

	if us.repository.Update(usr) == 0 {
		return errors.New("cannot update the user")
	}

	return nil
}

func (us *UserService) Delete(usrId int) error {

	if usrId < 1 {
		return errors.New("invalid id")
	}

	usrToUpdate := us.repository.GetById(usrId)

	if usrToUpdate.Id == 0 {
		return errors.New("user not found")
	}

	if us.repository.Delete(usrId) == 1 {
		return nil
	}

	return errors.New("user was not removed")
}
