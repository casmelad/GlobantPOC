package grpc

import (
	"context"
	"errors"
	"os"
	"testing"

	mappers "github.com/casmelad/GlobantPOC/cmd/grpcService/users/mappers"
	proto "github.com/casmelad/GlobantPOC/cmd/grpcService/users/proto"
	entities "github.com/casmelad/GlobantPOC/pkg/users"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
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

func (r applicationServiceMock) Create(ctx context.Context, u entities.User) (int, error) {
	args := r.Called(ctx, u)
	return args.Int(0), args.Error(1)
}

func (r applicationServiceMock) GetByEmail(ctx context.Context, email string) (entities.User, error) {
	args := r.Called(ctx, email)
	return args.Get(0).(entities.User), args.Error(1)
}

func (r applicationServiceMock) GetAll(ctx context.Context) ([]entities.User, error) {
	args := r.Called(ctx)
	return args.Get(0).([]entities.User), args.Error(1)
}

func (r applicationServiceMock) Update(ctx context.Context, u entities.User) error {
	args := r.Called(ctx, u)
	return args.Error(0)
}

func (r applicationServiceMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

var emailAddress string = "test@gmail.com"
var ctx context.Context = context.Background()
var applicationService applicationServiceMock = applicationServiceMock{}
var logger log.Logger = log.NewLogfmtLogger(os.Stderr)

var zipkinTracer *zipkin.Tracer

// Determine which OpenTracing tracer to use. We'll pass the tracer to all the
// components that use it, as a dependency.
var tracer stdopentracing.Tracer
var endpoints = NewGrpcUsersServer(applicationService)
var grpcService proto.UsersServer = NewGrpcUserServer(*endpoints, tracer, zipkinTracer, logger)

func Test_GetUser_ValidEmail_ReturnsUser(t *testing.T) {
	//Arrange
	expectedValue := proto.User{Email: emailAddress}
	applicationService.On("GetByEmail", expectedValue.Email).Return(entities.User{Email: expectedValue.Email}, nil).Once()
	//Act
	result, err := grpcService.GetUser(ctx, &proto.EmailAddress{Value: emailAddress})
	//Assert
	assert.Nil(t, err)
	assert.Equal(t, &proto.GetUserResponse{User: &expectedValue}, result)
	applicationService.AssertExpectations(t)
}

func Test_GetUser_InvalidEmail_ReturnsNotFoundError(t *testing.T) {
	//Arrange
	applicationService.On("GetByEmail", emailAddress).Return(entities.User{}, errors.New("user not found")).Once()
	//Act
	_, err := grpcService.GetUser(ctx, &proto.EmailAddress{Value: emailAddress})
	//Assert
	assert.NotNil(t, err)
	applicationService.AssertExpectations(t)
}

func Test_GetAll_ReturnsNoError(t *testing.T) {
	//Arrange
	applicationService.On("GetAll").Return([]entities.User{{}}, nil).Once()
	//Act
	users, err := grpcService.GetAllUsers(ctx, &proto.Filters{})
	//Assert
	assert.Nil(t, err)
	assert.NotEmpty(t, users)
	applicationService.AssertExpectations(t)
}

func Test_Create_ValidData_ReturnsNoError(t *testing.T) {
	//Arrange
	userToCreate := proto.User{}
	mappedUser, _ := mappers.ToDomainUser(userToCreate)
	applicationService.On("Create", mappedUser).Return(1, nil).Once()
	//Act
	result, err := grpcService.Create(ctx, &proto.CreateUserRequest{User: &userToCreate})
	//Arrange
	assert.Nil(t, err)
	assert.Equal(t, proto.CodeResult_OK, result.Code)
	applicationService.AssertExpectations(t)
}

func Test_Create_InvalidData_ReturnsInvalidDataError(t *testing.T) {
	//Arrange
	userToCreate := proto.User{}
	mappedUser, _ := mappers.ToDomainUser(userToCreate)
	applicationService.On("Create", mappedUser).Return(0, errors.New("some data is missing")).Once()
	//Act
	result, err := grpcService.Create(ctx, &proto.CreateUserRequest{User: &userToCreate})
	//Arrange
	assert.NotNil(t, err)
	assert.Equal(t, proto.CodeResult_INVALIDINPUT, result.Code)
	applicationService.AssertExpectations(t)
}

func Test_Create_DuplicatedData_ReturnsAlreadyExistsError(t *testing.T) {
	//Arrange
	userToCreate := proto.User{}
	mappedUser, _ := mappers.ToDomainUser(userToCreate)
	applicationService.On("Create", mappedUser).Return(0, errors.New("user already exists")).Once()
	//Act
	result, err := grpcService.Create(ctx, &proto.CreateUserRequest{User: &userToCreate})
	//Arrange
	assert.NotNil(t, err)
	assert.Equal(t, proto.CodeResult_FAILED, result.Code)
	applicationService.AssertExpectations(t)
}

func Test_Update_ValidData_ReturnsNoError(t *testing.T) {
	//Arrange
	userToCreate := proto.User{}
	mappedUser, _ := mappers.ToDomainUser(userToCreate)
	applicationService.On("Update", mappedUser).Return(nil).Once()
	//Act
	result, err := grpcService.Update(ctx, &proto.UpdateUserRequest{User: &userToCreate})
	//Arrange
	assert.Nil(t, err)
	assert.Equal(t, proto.CodeResult_OK, result.Code)
	applicationService.AssertExpectations(t)
}

func Test_Update_InvalidUserData_ReturnsNotFoundError(t *testing.T) {
	//Arrange
	userToCreate := proto.User{}
	mappedUser, _ := mappers.ToDomainUser(userToCreate)
	applicationService.On("Update", mappedUser).Return(errors.New("user not found")).Once()
	//Act
	result, err := grpcService.Update(ctx, &proto.UpdateUserRequest{User: &userToCreate})
	//Arrange
	assert.NotNil(t, err)
	assert.Equal(t, proto.CodeResult_NOTFOUND, result.Code)
	applicationService.AssertExpectations(t)
}

func Test_Update_InvalidUserData_ReturnsInvalidInputError(t *testing.T) {
	//Arrange
	userToCreate := proto.User{}
	mappedUser, _ := mappers.ToDomainUser(userToCreate)
	applicationService.On("Update", mappedUser).Return(errors.New("invalid input")).Once()
	//Act
	result, err := grpcService.Update(ctx, &proto.UpdateUserRequest{User: &userToCreate})
	//Arrange
	assert.NotNil(t, err)
	assert.Equal(t, proto.CodeResult_INVALIDINPUT, result.Code)
	applicationService.AssertExpectations(t)
}

func Test_Update_InvalidUserData_ReturnsFailedError(t *testing.T) {
	//Arrange
	userToCreate := proto.User{}
	mappedUser, _ := mappers.ToDomainUser(userToCreate)
	applicationService.On("Update", mappedUser).Return(errors.New("cannot update the user")).Once()
	//Act
	result, err := grpcService.Update(ctx, &proto.UpdateUserRequest{User: &userToCreate})
	//Arrange
	assert.NotNil(t, err)
	assert.Equal(t, proto.CodeResult_FAILED, result.Code)
	applicationService.AssertExpectations(t)
}

func Test_Delete_InvalidId_ReturnsNotFoundError(t *testing.T) {
	//Arrange
	applicationService.On("Delete", 1).Return(errors.New("user not found")).Once()
	//Act
	result, err := grpcService.Delete(ctx, &proto.Id{Value: 1})
	//Assert
	assert.NotNil(t, err)
	assert.Equal(t, proto.CodeResult_NOTFOUND, result.Code)
	applicationService.AssertExpectations(t)
}

func Test_Delete_IdZero_ReturnsInvalidInputError(t *testing.T) {
	//Arrange
	applicationService.On("Delete", 0).Return(errors.New("invalid id")).Once()
	//Act
	result, err := grpcService.Delete(ctx, &proto.Id{Value: 0})
	//Assert
	assert.NotNil(t, err)
	assert.Equal(t, proto.CodeResult_INVALIDINPUT, result.Code)
	applicationService.AssertExpectations(t)
}

func Test_Delete_InternalError_ReturnsError(t *testing.T) {
	//Arrange
	applicationService.On("Delete", 1).Return(errors.New("user was not removed")).Once()
	//Act
	result, err := grpcService.Delete(ctx, &proto.Id{Value: 1})
	//Assert
	assert.NotNil(t, err)
	assert.Equal(t, proto.CodeResult_FAILED, result.Code)
	applicationService.AssertExpectations(t)
}
