package users

import (
	"context"
	"errors"
	"fmt"
)

type UsersService struct {
	repo map[string]UserResponse
}

func NewUserService() *UsersService {
	return &UsersService{
		repo: map[string]UserResponse{},
	}
}

func (u *UsersService) GetUser(c context.Context, uid *UserEmailRequest) (*UserResponse, error) {

	if user, exists := u.repo[uid.EMail]; exists {
		return &user, nil
	}

	return &UserResponse{}, errors.New("user not found")
}

func (u *UsersService) Create(ctx context.Context, user *UserRequest) (*CreateUserResponse, error) {

	if user == nil {
		return &CreateUserResponse{Code: CreateUserResponse_INVALIDINPUT}, nil
	}

	user.Id = int32(len(u.repo)) + 1
	_, exists := u.repo[user.Email]

	if !exists {
		u.repo[user.Email] = UserResponse{Id: user.Id, Email: user.Email, Name: user.Name, LastName: user.LastName}
		userCreated, _ := u.GetUser(ctx, &UserEmailRequest{EMail: user.Email})
		fmt.Println(userCreated)
		return &CreateUserResponse{Code: CreateUserResponse_OK, User: userCreated}, nil
	}

	return &CreateUserResponse{Code: CreateUserResponse_FAILED}, nil

}

func (u *UsersService) GetAllUsers(ctx context.Context, v *Void) (*UserCollectionResponse, error) {

	response := &UserCollectionResponse{Users: []*UserResponse{}}
	ch1 := make(chan UserResponse, len(u.repo))

	for _, user := range u.repo {
		go func(usr UserResponse) {
			ch1 <- usr
		}(user)
	}

	for i := 0; i < len(u.repo); i++ {
		userT := <-ch1
		response.Users = append(response.Users, &userT)
	}

	return response, nil
}

func (u *UsersService) Update(ctx context.Context, user *UserRequest) (*UpdateUserResponse, error) {

	if user == nil {
		return &UpdateUserResponse{Code: UpdateUserResponse_INVALIDINPUT}, errors.New("invalid data")
	}

	userToUpdate, err := u.GetUser(ctx, &UserEmailRequest{EMail: user.Email})

	if err != nil {
		return &UpdateUserResponse{Code: UpdateUserResponse_FAILED}, err
	}

	userToUpdate.Name = user.Name
	userToUpdate.LastName = user.LastName
	u.repo[user.Email] = *userToUpdate

	return &UpdateUserResponse{Code: UpdateUserResponse_FAILED}, nil
}

func (u *UsersService) Delete(ctx context.Context, userId *UserId) (*DeleteUserResponse, error) {

	ch := make(chan UserResponse)

	go func() {
		for _, userFromRepo := range u.repo {
			if userFromRepo.Id == userId.Id {
				ch <- userFromRepo
			}
		}
		ch <- UserResponse{}
	}()

	user := <-ch
	userToRemove, err := u.GetUser(ctx, &UserEmailRequest{EMail: user.Email})

	if err != nil && userToRemove.Id == 0 {
		return &DeleteUserResponse{Code: DeleteUserResponse_FAILED}, err
	}

	delete(u.repo, userToRemove.Email)
	return &DeleteUserResponse{Code: DeleteUserResponse_OK}, nil
}
