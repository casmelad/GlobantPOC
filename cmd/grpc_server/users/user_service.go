package users

import (
	"context"
	"errors"
	"fmt"
)

type UsersService struct {
	repo map[string]User
}

func NewUserService() *UsersService {

	return &UsersService{
		repo: map[string]User{},
	}

}

func (u *UsersService) GetUser(c context.Context, uid *UserEmail) (*User, error) {

	if user, ok := u.repo[uid.EMail]; ok {
		return &user, nil
	}

	return &User{}, errors.New("user not found")
}

func (u *UsersService) Create(ctx context.Context, user *User) (*TaskResult, error) {

	if user == nil {
		return &TaskResult{HasBody: true, Result: 0, Code: TaskResult_InvalidInput}, nil
	}

	user.Id = int32(len(u.repo)) + 1

	_, ok := u.repo[user.EMail]

	if !ok {
		u.repo[user.EMail] = User{Id: user.Id, EMail: user.EMail, Name: user.Name, LastName: user.LastName}
		return &TaskResult{HasBody: true, Result: user.Id, Code: TaskResult_Ok}, nil
	}

	return &TaskResult{HasBody: true, Result: 0, Code: TaskResult_Failed}, nil

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

	if user == nil {
		return &TaskResult{Code: TaskResult_InvalidInput}, errors.New("invalid data")
	}

	userToUpdate, err := u.GetUser(ctx, &UserEmail{EMail: user.EMail})

	fmt.Println(err)

	if err != nil {
		return &TaskResult{Code: TaskResult_Failed}, err
	}

	userToUpdate.Name = user.Name
	userToUpdate.LastName = user.LastName

	u.repo[user.EMail] = *userToUpdate

	return &TaskResult{Code: TaskResult_Ok}, nil
}

func (u *UsersService) Delete(ctx context.Context, userId *UserId) (*TaskResult, error) {

	ch := make(chan User)

	go func() {
		for _, userFromRepo := range u.repo {
			if userFromRepo.Id == userId.Id {
				ch <- userFromRepo
			}
		}
		ch <- User{}
	}()

	user := <-ch

	userToRemove, err := u.GetUser(ctx, &UserEmail{EMail: user.EMail})

	if err != nil && userToRemove.Id == 0 {
		return &TaskResult{Code: TaskResult_Failed}, err
	}

	delete(u.repo, userToRemove.EMail)

	return &TaskResult{Code: TaskResult_Ok}, nil
}
