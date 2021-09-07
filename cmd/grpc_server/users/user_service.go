package users

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type UsersService struct {
	mtx  sync.RWMutex
	repo map[string]User
}

func NewUserService() *UsersService {

	return &UsersService{
		repo: map[string]User{},
	}

}

func (u *UsersService) GetUserById(c context.Context, uid *UserId) (*User, error) {

	fmt.Println("entra")

	if uid.Id <= 0 {
		fmt.Println("Invalid user Id")
		return nil, errors.New("Invalid user Id")
	}

	if user, ok := u.repo[""]; ok {
		return &user, nil
	}

	return &User{}, nil
}

func (u *UsersService) Create(ctx context.Context, user *User) (*TaskResult, error) {

	u.mtx.Lock()
	defer u.mtx.Unlock()

	if user == nil {
		fmt.Println("entra crea es nil")
	}

	user.Id = int32(len(u.repo)) + 1

	if _, ok := u.repo[user.Name]; !ok {
		fmt.Println(" no existe ")
		u.repo[user.Name] = User{Id: user.Id, Name: user.Name, LastName: user.LastName}
	} else {
		fmt.Println("existe")
	}

	fmt.Println(u.repo)

	return &TaskResult{Code: 1}, nil

}
func (u *UsersService) GetAllUsers(ctx context.Context, v *Void) (*UserCollection, error) {

	response := &UserCollection{Users: []*User{}}

	ch1 := make(chan User, len(u.repo))

	for _, user := range u.repo {
		go func(usr User) {
			ch1 <- usr
		}(user)
	}

	for i := 0; i < len(u.repo); i++ {
		userT := <-ch1
		response.Users = append(response.Users, &userT)
	}

	return response, nil
}

func (u *UsersService) Update(ctx context.Context, user *User) (*TaskResult, error) {
	return nil, nil
}
func (u *UsersService) Delete(ctx context.Context, user *UserId) (*TaskResult, error) {
	return nil, nil
}
