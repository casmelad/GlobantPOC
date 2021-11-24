package users

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type grpcProxyMock struct {
	mock.Mock
}

func (up grpcProxyMock) GetAll(ctx context.Context) ([]User, error) {
	args := up.Called(ctx)
	return args.Get(0).([]User), args.Error(1)
}

func (up grpcProxyMock) Create(ctx context.Context, u User) (User, error) {
	args := up.Called(ctx, u)
	return args.Get(0).(User), args.Error(1)
}

func (up grpcProxyMock) Update(ctx context.Context, u User) (User, error) {
	args := up.Called(ctx, u)
	return args.Get(0).(User), args.Error(1)
}

func (up grpcProxyMock) Delete(ctx context.Context, id int) (bool, error) {
	args := up.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (up grpcProxyMock) GetByEmail(ctx context.Context, email string) (User, error) {
	args := up.Called(ctx, email)
	return args.Get(0).(User), args.Error(1)
}

func TestCases_Create(t *testing.T) {
	ctx := context.Background()

	for _, useCase := range createUserTestCases {
		proxyMock := grpcProxyMock{}
		usrRequest, is := useCase.requestData.(postUserRequest)

		if is {
			proxyMock.On("Create", ctx, usrRequest.User).Return(useCase.responseData, useCase.err)
		}

		endpoints := MakeServerEndpoints(proxyMock)
		result, err := endpoints.PostUserEndpoint(ctx, useCase.requestData)

		if appError, is := result.(AppError); is {
			assert.Equal(t, appError.error(), useCase.err)
		}

		if err != nil {
			assert.EqualError(t, err, "invalid input data")
		}

		if usr, is := result.(postUserResponse); is {
			assert.NotEmpty(t, usr.Href)
		}
	}
}

var createUserTestCases []struct {
	name         string
	requestData  interface{}
	responseData User
	err          error
} = []struct {
	name         string
	requestData  interface{}
	responseData User
	err          error
}{
	{"ValidData_ValidResponse", postUserRequest{User: larryPage}, larryPage, nil},
	{"IncompleteData_BadRequest", postUserRequest{User: User{Email: ""}}, User{}, ErrInvalidInput},
	{"ExistingUser_ReturnsConflict", postUserRequest{User: larryPage}, User{}, ErrUserAlreadyExists},
	{"InvalidRequestUri_InternalError", User{}, User{}, errors.New("invalid input data")},
}

func TestCases_GetAll(t *testing.T) {

	ctx := context.Background()

	for _, useCase := range getAllUsersTestCases {
		proxyMock := grpcProxyMock{}
		proxyMock.On("GetAll", ctx).Return(useCase.users, useCase.err)
		endpoints := MakeServerEndpoints(proxyMock)
		result, err := endpoints.GetAllUsersEndpoint(ctx, getAllUsersRequest{})

		if appError, is := result.(AppError); is {
			assert.Equal(t, appError.error(), ErrInternalFailure)
		}

		if err != nil {
			assert.EqualError(t, err, "invalid input data")
		}

		if users, is := result.(getAllUsersResponse); is {
			assert.Equal(t, useCase.users, users.Users)
			assert.ElementsMatch(t, useCase.users, users.Users)
		}
	}

}

var getAllUsersTestCases []struct {
	testCaseName string
	users        []User
	err          error
} = []struct {
	testCaseName string
	users        []User
	err          error
}{
	{"Ok_ReturnsData", []User{larryPage}, nil},
	{"Ok_InternalError", []User{}, ErrInternalFailure},
}

var larryPage User = User{Id: 1, Name: "Larry", LastName: "Page", Email: "larry.page@gmail.com"}
