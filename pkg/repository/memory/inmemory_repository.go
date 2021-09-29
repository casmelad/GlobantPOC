package infrastructure

import (
	"context"
	"fmt"

	"github.com/casmelad/GlobantPOC/pkg/users"
)

//InMemoryUserRepository is an in memory implementation of user Repository
type InMemoryUserRepository struct {
	dict   map[string]users.User
	regist []int
}

//NewInMemoryUserRepository returns an InMemoryUserRepository type pointer
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		dict:   map[string]users.User{},
		regist: []int{},
	}
}

//Add - adds a user to the repository
func (repo *InMemoryUserRepository) Add(ctx context.Context, u users.User) (int, error) {

	fmt.Println(u)
	if _, ok := repo.dict[u.Email]; ok {
		return 0, nil
	}

	u.ID = len(repo.regist) + 1
	repo.regist = append(repo.regist, u.ID)
	repo.dict[u.Email] = u

	return u.ID, nil
}

//GetByID - retrieves a user from the repository based on the integer id
func (repo *InMemoryUserRepository) GetByID(ctx context.Context, userID int) (users.User, error) {

	for _, usr := range repo.dict {
		if usr.ID == userID {
			return usr, nil
		}
	}

	return users.User{}, nil
}

//GetByEmail - retrieves a user from the repository based on the email address
func (repo *InMemoryUserRepository) GetByEmail(ctx context.Context, id string) (users.User, error) {
	return repo.dict[id], nil
}

//GetAll - retrieves all the users from the repository
func (repo *InMemoryUserRepository) GetAll(ctx context.Context) ([]users.User, error) {

	result := []users.User{}

	for _, usr := range repo.dict {
		result = append(result, usr)
	}

	return result, nil
}

//Update -  updates the information of a user
func (repo *InMemoryUserRepository) Update(ctx context.Context, u users.User) error {

	userToUpdate, err := repo.GetByEmail(ctx, u.Email)

	if err != nil {
		return err
	}

	if userToUpdate.ID > 0 {
		userToUpdate.Name = u.Name
		userToUpdate.LastName = u.LastName
		repo.dict[userToUpdate.Email] = userToUpdate
	}

	return nil
}

//Delete - deletes a user from the repository
func (repo *InMemoryUserRepository) Delete(ctx context.Context, userID int) error {

	for _, usr := range repo.dict {
		if usr.ID == userID {
			delete(repo.dict, usr.Email)
			return nil
		}
	}

	return nil
}
