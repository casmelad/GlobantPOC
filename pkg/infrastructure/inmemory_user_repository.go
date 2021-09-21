package infrastructure

import (
	"github.com/casmelad/GlobantPOC/pkg/domain/entities"
)

type InMemoryUserRepository struct {
	dict   map[string]entities.User
	regist []int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		dict:   map[string]entities.User{},
		regist: []int{},
	}
}

func (repo *InMemoryUserRepository) Add(u entities.User) int {

	if _, ok := repo.dict[u.Email]; ok {
		return 0
	}

	u.Id = len(repo.regist) + 1
	repo.regist = append(repo.regist, u.Id)
	repo.dict[u.Email] = u

	return u.Id
}

func (repo *InMemoryUserRepository) GetById(userId int) entities.User {

	for _, usr := range repo.dict {
		if usr.Id == userId {
			delete(repo.dict, usr.Email)
			return usr
		}
	}

	return entities.User{}
}

func (repo *InMemoryUserRepository) GetByEmail(id string) entities.User {
	return repo.dict[id]
}

func (repo *InMemoryUserRepository) GetAll() []entities.User {

	result := []entities.User{}

	for _, usr := range repo.dict {
		result = append(result, usr)
	}

	return result
}

func (repo *InMemoryUserRepository) Update(u entities.User) int {

	result := 0
	userToUpdate := repo.GetByEmail(u.Email)

	if userToUpdate.Id > 0 {
		userToUpdate.Name = u.Name
		userToUpdate.LastName = u.LastName
		repo.dict[userToUpdate.Email] = userToUpdate
		result = 1
	}

	return result
}

func (repo *InMemoryUserRepository) Delete(userId int) int {

	result := 0

	for _, usr := range repo.dict {
		if usr.Id == userId {
			delete(repo.dict, usr.Email)
			result = 1
			break
		}
	}

	return result
}
