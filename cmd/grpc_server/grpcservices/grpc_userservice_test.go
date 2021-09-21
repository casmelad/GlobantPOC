package grpcservices

import (
	"context"
	"errors"
	"testing"

	"github.com/casmelad/GlobantPOC/cmd/grpc_server/users"
	"github.com/casmelad/GlobantPOC/pkg/domain/entities"
	"github.com/casmelad/GlobantPOC/pkg/mappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/* Create(entities.User) (int, error)
GetById(string) (entities.User, error)
GetAll() ([]entities.User, error)
Update(entities.User) error
Delete(int) error */

type applicationServiceMock struct {
	mock.Mock
}

func (r *applicationServiceMock) Create(u entities.User) (int, error) {
	args := r.Called(u)
	return args.Int(0), args.Error(1)
}

func (r *applicationServiceMock) GetByEmail(email string) (entities.User, error) {
	args := r.Called(email)
	return args.Get(0).(entities.User), args.Error(1)
}

func (r *applicationServiceMock) GetAll() ([]entities.User, error) {
	args := r.Called()
	return args.Get(0).([]entities.User), args.Error(1)
}

func (r *applicationServiceMock) Update(u entities.User) error {
	args := r.Called(u)
	return args.Error(0)
}

func (r *applicationServiceMock) Delete(id int) error {
	args := r.Called(id)
	return args.Error(0)
}

var emailAddress string = "test@gmail.com"
var ctx context.Context = context.Background()
var applicationService applicationServiceMock = applicationServiceMock{}
var grpcService users.UsersServer = NewGrpcUserService(&applicationService)

func Test_GetUser_ValidEmail_ReturnsUser(t *testing.T) {
	//Arrange
	expectedValue := users.User{Email: emailAddress}
	applicationService.On("GetByEmail", expectedValue.Email).Return(entities.User{Email: expectedValue.Email}, nil).Once()
	//Act
	result, err := grpcService.GetUser(ctx, &users.EmailAddress{Value: emailAddress})
	//Assert
	assert.Nil(t, err)
	assert.Equal(t, &users.GetUserResponse{User: &expectedValue}, result)
}

func Test_GetUser_InvalidEmail_ReturnsNotFoundError(t *testing.T) {
	//Arrange
	applicationService.On("GetByEmail", emailAddress).Return(entities.User{}, errors.New("user not found")).Once()
	//Act
	_, err := grpcService.GetUser(ctx, &users.EmailAddress{Value: emailAddress})
	//Assert
	assert.NotNil(t, err)
}

func Test_GetAll_ReturnsNoError(t *testing.T) {
	//Arrange
	applicationService.On("GetAll").Return([]entities.User{{}}, nil).Once()
	//Act
	users, err := grpcService.GetAllUsers(ctx, &users.Filters{})
	//Assert
	assert.Nil(t, err)
	assert.NotEmpty(t, users.Users)
}

func Test_Create_ValidData_ReturnsNoError(t *testing.T) {
	//Arrange
	userToCreate := users.User{}
	mappedUser, _ := mappers.ToDomainUser(userToCreate)
	applicationService.On("Create", mappedUser).Return(1, nil).Once()
	//Act
	result, err := grpcService.Create(ctx, &users.CreateUserRequest{User: &userToCreate})
	//Arrange
	assert.Nil(t, err)
	assert.Equal(t, &users.CreateUserResponse{UserId: 1, Code: users.CodeResult_OK}, result)
}

func Test_Create_InvalidData_ReturnsInvalidDataError(t *testing.T) {
	//Arrange
	userToCreate := users.User{}
	mappedUser, _ := mappers.ToDomainUser(userToCreate)
	applicationService.On("Create", mappedUser).Return(0, errors.New("some data is missing")).Once()
	//Act
	result, err := grpcService.Create(ctx, &users.CreateUserRequest{User: &userToCreate})
	//Arrange
	assert.NotNil(t, err)
	assert.Equal(t, &users.CreateUserResponse{UserId: 0, Code: users.CodeResult_INVALIDINPUT}, result)
}

func Test_Create_DuplicatedData_ReturnsAlreadyExistsError(t *testing.T) {
	//Arrange
	userToCreate := users.User{}
	mappedUser, _ := mappers.ToDomainUser(userToCreate)
	applicationService.On("Create", mappedUser).Return(0, errors.New("user already exists")).Once()
	//Act
	result, err := grpcService.Create(ctx, &users.CreateUserRequest{User: &userToCreate})
	//Arrange
	assert.NotNil(t, err)
	assert.Equal(t, &users.CreateUserResponse{UserId: 0, Code: users.CodeResult_FAILED}, result)
}

func Test_Update_ValidData_ReturnsNoError(t *testing.T) {
	//Arrange
	userToCreate := users.User{}
	mappedUser, _ := mappers.ToDomainUser(userToCreate)
	applicationService.On("Update", mappedUser).Return(nil).Once()
	//Act
	result, err := grpcService.Update(ctx, &users.UpdateUserRequest{User: &userToCreate})
	//Arrange
	assert.Nil(t, err)
	assert.Equal(t, &users.UpdateUserResponse{Code: users.CodeResult_OK}, result)
}

func Test_Update_InvalidUserData_ReturnsNotFoundError(t *testing.T) {
	//Arrange
	userToCreate := users.User{}
	mappedUser, _ := mappers.ToDomainUser(userToCreate)
	applicationService.On("Update", mappedUser).Return(errors.New("user not found")).Once()
	//Act
	result, err := grpcService.Update(ctx, &users.UpdateUserRequest{User: &userToCreate})
	//Arrange
	assert.NotNil(t, err)
	assert.Equal(t, users.CodeResult_NOTFOUND, result.Code)
}

func Test_Update_InvalidUserData_ReturnsInvalidInputError(t *testing.T) {
	//Arrange
	userToCreate := users.User{}
	mappedUser, _ := mappers.ToDomainUser(userToCreate)
	applicationService.On("Update", mappedUser).Return(errors.New("invalid input")).Once()
	//Act
	result, err := grpcService.Update(ctx, &users.UpdateUserRequest{User: &userToCreate})
	//Arrange
	assert.NotNil(t, err)
	assert.Equal(t, users.CodeResult_INVALIDINPUT, result.Code)
}

func Test_Update_InvalidUserData_ReturnsFailedError(t *testing.T) {
	//Arrange
	userToCreate := users.User{}
	mappedUser, _ := mappers.ToDomainUser(userToCreate)
	applicationService.On("Update", mappedUser).Return(errors.New("cannot update the user")).Once()
	//Act
	result, err := grpcService.Update(ctx, &users.UpdateUserRequest{User: &userToCreate})
	//Arrange
	assert.NotNil(t, err)
	assert.Equal(t, users.CodeResult_FAILED, result.Code)
}

func Test_Delete_InvalidId_ReturnsNotFoundError(t *testing.T) {
	//Arrange
	applicationService.On("Delete", 1).Return(errors.New("user not found")).Once()
	//Act
	result, err := grpcService.Delete(ctx, &users.Id{Value: 1})
	//Assert
	assert.NotNil(t, err)
	assert.Equal(t, users.CodeResult_NOTFOUND, result.Code)

}

func Test_Delete_IdZero_ReturnsInvalidInputError(t *testing.T) {
	//Arrange
	applicationService.On("Delete", 0).Return(errors.New("invalid id")).Once()
	//Act
	result, err := grpcService.Delete(ctx, &users.Id{Value: 0})
	//Assert
	assert.NotNil(t, err)
	assert.Equal(t, users.CodeResult_INVALIDINPUT, result.Code)

}

func Test_Delete_InternalError_ReturnsError(t *testing.T) {
	//Arrange
	applicationService.On("Delete", 1).Return(errors.New("user was not removed")).Once()
	//Act
	result, err := grpcService.Delete(ctx, &users.Id{Value: 1})
	//Assert
	assert.NotNil(t, err)
	assert.Equal(t, users.CodeResult_FAILED, result.Code)

}
