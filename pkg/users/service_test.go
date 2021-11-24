package users

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) Add(ctx context.Context, u User) (int, error) {
	args := r.Called(ctx, u)
	return args.Int(0), args.Error(1)
}

func (r *repositoryMock) Update(ctx context.Context, u User) error {
	args := r.Called(ctx, u)
	return args.Error(0)
}

func (r *repositoryMock) Delete(ctx context.Context, uid int) error {
	args := r.Called(ctx, uid)
	return args.Error(0)
}

func (r *repositoryMock) GetByID(ctx context.Context, uid int) (User, error) {
	args := r.Called(ctx, uid)
	return args.Get(0).(User), args.Error(1)
}

func (r *repositoryMock) GetByEmail(ctx context.Context, email string) (User, error) {
	args := r.Called(ctx, email)
	return args.Get(0).(User), args.Error(1)
}

func (r *repositoryMock) GetAll(ctx context.Context) ([]User, error) {
	args := r.Called(ctx)
	return args.Get(0).([]User), args.Error(1)
}

func Test_Create_ValidData_OkResult(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	userToAdd := User{Email: "test@gmail.com", Name: "John", LastName: "Connor"}
	repository.On("Add", context.Background(), userToAdd).Return(1, nil)
	repository.On("GetByEmail", context.Background(), userToAdd.Email).Return(User{}, nil)
	//Act
	result, err := service.Create(context.Background(), userToAdd)
	//Assert
	assert.Greater(t, result, 0)
	assert.Nil(t, err)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "Add", 1)
	repository.AssertNumberOfCalls(t, "GetByEmail", 1)
}

func Test_Create_DuplicatedData_ReturnsAlreadyExistsError(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	userToAdd := User{ID: 1, Email: "test@gmail.com", Name: "John", LastName: "Connor"}
	repository.On("GetByEmail", context.Background(), userToAdd.Email).Return(userToAdd, errors.New(""))
	//Act
	result, err := service.Create(context.Background(), userToAdd)
	//Assert
	assert.Equal(t, 0, result)
	assert.NotNil(t, err)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetByEmail", 1)
}

func Test_Create_InvalidData_ReturnsError(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	userToAdd := User{Email: "test@gmail.com"}
	//Act
	id, err := service.Create(context.Background(), userToAdd)
	//Assert
	assert.Equal(t, 0, id)
	assert.NotNil(t, err)
}

func Test_Update_ValidData_OkResult(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	userToUpdate := User{ID: 1, Email: "test@gmail.com", Name: "John", LastName: "Connor"}
	repository.On("Update", context.Background(), userToUpdate).Return(nil)
	repository.On("GetByEmail", context.Background(), userToUpdate.Email).Return(userToUpdate, nil)
	//Act
	err := service.Update(context.Background(), userToUpdate)
	//Assert
	assert.Nil(t, err)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetByEmail", 1)
	repository.AssertNumberOfCalls(t, "Update", 1)
}

func Test_Update_InvalidData_ReturnsError(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	userToUpdate := User{ID: 1, Email: "test@gmail.com", Name: "", LastName: "Connor"}
	//Act
	err := service.Update(context.Background(), userToUpdate)
	//Assert
	assert.NotNil(t, err)
}

func Test_Update_InvalidUser_ReturnsError(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	userToUpdate := User{ID: 1, Email: "test@gmail.com", Name: "John", LastName: "Connor"}
	repository.On("GetByEmail", context.Background(), userToUpdate.Email).Return(User{}, errors.New(""))
	//Act
	err := service.Update(context.Background(), userToUpdate)
	//Assert
	assert.NotNil(t, err)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetByEmail", 1)
}

func Test_Delete_ValidId_DeletesUser(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	repository.On("GetByID", context.Background(), 1).Return(User{ID: 1}, nil)
	repository.On("Delete", context.Background(), 1).Return(nil)
	//Act
	result := service.Delete(context.Background(), 1)
	//Assert
	assert.Nil(t, result)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetByID", 1)
	repository.AssertNumberOfCalls(t, "Delete", 1)
}

func Test_Delete_InvalidId_ReturnsError(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	//Act
	result := service.Delete(context.Background(), 0)
	//Assert
	assert.NotNil(t, result)
	assert.Equal(t, "invalid id", result.Error())
}

func Test_Delete_InvalidId_ReturnsNotFoundError(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	repository.On("GetByID", context.Background(), 999).Return(User{}, nil)
	//Act
	result := service.Delete(context.Background(), 999)
	//Assert
	assert.NotNil(t, result)
	assert.Equal(t, "user not found", result.Error())
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetByID", 1)
}

func Test_GetByEmail_ValidId_ReturnsData(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	expectedResult := User{ID: 1, Email: "test@gmail.com"}
	repository.On("GetByEmail", context.Background(), "test@gmail.com").Return(expectedResult, nil)
	//Act
	result, err := service.GetByEmail(context.Background(), "test@gmail.com")
	//
	assert.Equal(t, expectedResult, result)
	assert.Nil(t, err)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetByEmail", 1)
}

func Test_GetByEmail_NotValidId_ReturnsErrorNotFound(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	repository.On("GetByEmail", context.Background(), "test@gmail.com").Return(User{}, nil)
	//Act
	_, err := service.GetByEmail(context.Background(), "test@gmail.com")
	//
	assert.NotNil(t, err)
	assert.Equal(t, "user not found", err.Error())
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetByEmail", 1)
}

func Test_GetAll_ReturnsNoError(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	repository.On("GetAll", context.Background()).Return([]User{}, nil)
	//Act
	_, err := service.GetAll(context.Background())
	//
	assert.Nil(t, err)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetAll", 1)
}

func Test_GetByEmail(t *testing.T) {

	repository := repositoryMock{}
	service := NewUserService(&repository)

	for _, testCase := range getByEmailTestCases {

		repository.On("GetByEmail", context.Background(), testCase.email).Return(testCase.userExpected, testCase.errExpected).Once()

		result, err := service.GetByEmail(context.Background(), testCase.email)

		assert.Equal(t, testCase.userExpected, result)
		assert.Equal(t, testCase.errExpected, err)
	}
}

var getByEmailTestCases []struct {
	name         string
	email        string
	userExpected User
	errExpected  error
} = []struct {
	name         string
	email        string
	userExpected User
	errExpected  error
}{
	{"ValidId_ReturnsData", "test@gmail.com", User{ID: 1, Email: "test@gmail.com"}, nil},
	{"NotValidId_ReturnsErrorNotFound", "test1@gmail.com", User{}, errors.New("user not found")},
}
