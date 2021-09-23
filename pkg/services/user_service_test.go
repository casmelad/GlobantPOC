package services

import (
	"errors"
	"testing"

	"github.com/casmelad/GlobantPOC/pkg/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) Add(u entities.User) int {
	args := r.Called(u)
	return args.Int(0)
}

func (r *repositoryMock) Update(u entities.User) int {
	args := r.Called(u)
	return args.Int(0)
}

func (r *repositoryMock) Delete(uid int) int {
	args := r.Called(uid)
	return args.Int(0)
}

func (r *repositoryMock) GetById(uid int) entities.User {
	args := r.Called(uid)
	return args.Get(0).(entities.User)
}

func (r *repositoryMock) GetByEmail(email string) entities.User {
	args := r.Called(email)
	return args.Get(0).(entities.User)
}

func (r *repositoryMock) GetAll() []entities.User {
	args := r.Called()
	return args.Get(0).([]entities.User)
}

func Test_Create_ValidData_OkResult(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	userToAdd := entities.User{Email: "test@gmail.com", Name: "John", LastName: "Connor"}
	repository.On("Add", userToAdd).Return(1)
	repository.On("GetByEmail", userToAdd.Email).Return(entities.User{})
	//Act
	result, err := service.Create(userToAdd)
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
	userToAdd := entities.User{Id: 1, Email: "test@gmail.com", Name: "John", LastName: "Connor"}
	repository.On("GetByEmail", userToAdd.Email).Return(userToAdd)
	//Act
	result, err := service.Create(userToAdd)
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
	userToAdd := entities.User{Email: "test@gmail.com"}
	//Act
	id, err := service.Create(userToAdd)
	//Assert
	assert.Equal(t, 0, id)
	assert.NotNil(t, err)
}

func Test_Update_ValidData_OkResult(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	userToUpdate := entities.User{Id: 1, Email: "test@gmail.com", Name: "John", LastName: "Connor"}
	repository.On("Update", userToUpdate).Return(1)
	repository.On("GetByEmail", userToUpdate.Email).Return(userToUpdate)
	//Act
	err := service.Update(userToUpdate)
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
	userToUpdate := entities.User{Id: 1, Email: "test@gmail.com", Name: "", LastName: "Connor"}
	//Act
	err := service.Update(userToUpdate)
	//Assert
	assert.NotNil(t, err)
}

func Test_Update_InvalidUser_ReturnsError(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	userToUpdate := entities.User{Id: 1, Email: "test@gmail.com", Name: "John", LastName: "Connor"}
	repository.On("GetByEmail", userToUpdate.Email).Return(entities.User{})
	//Act
	err := service.Update(userToUpdate)
	//Assert
	assert.NotNil(t, err)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetByEmail", 1)
}

func Test_Delete_ValidId_DeletesUser(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	repository.On("GetById", 1).Return(entities.User{Id: 1})
	repository.On("Delete", 1).Return(1)
	//Act
	result := service.Delete(1)
	//Assert
	assert.Nil(t, result)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetById", 1)
	repository.AssertNumberOfCalls(t, "Delete", 1)
}

func Test_Delete_InvalidId_ReturnsError(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	//Act
	result := service.Delete(0)
	//Assert
	assert.NotNil(t, result)
	assert.Equal(t, "invalid id", result.Error())
}

func Test_Delete_InvalidId_ReturnsNotFoundError(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	repository.On("GetById", 999).Return(entities.User{})
	//Act
	result := service.Delete(999)
	//Assert
	assert.NotNil(t, result)
	assert.Equal(t, "user not found", result.Error())
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetById", 1)
}

func Test_GetByEmail_ValidId_ReturnsData(t *testing.T) {
	//Arrange
	repository := repositoryMock{}
	service := NewUserService(&repository)
	expectedResult := entities.User{Id: 1, Email: "test@gmail.com"}
	repository.On("GetByEmail", "test@gmail.com").Return(expectedResult)
	//Act
	result, err := service.GetByEmail("test@gmail.com")
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
	repository.On("GetByEmail", "test@gmail.com").Return(entities.User{})
	//Act
	_, err := service.GetByEmail("test@gmail.com")
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
	repository.On("GetAll").Return([]entities.User{})
	//Act
	_, err := service.GetAll()
	//
	assert.Nil(t, err)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetAll", 1)
}

func Test_GetByEmail(t *testing.T) {

	repository := repositoryMock{}
	service := NewUserService(&repository)

	for _, testCase := range tests {

		repository.On("GetByEmail", testCase.email).Return(testCase.userExpected).Once()

		result, err := service.GetByEmail(testCase.email)

		assert.Equal(t, testCase.userExpected, result)
		assert.Equal(t, testCase.errExpected, err)
	}
}

var tests []struct {
	name         string
	email        string
	userExpected entities.User
	errExpected  error
} = []struct {
	name         string
	email        string
	userExpected entities.User
	errExpected  error
}{
	{"ValidId_ReturnsData", "test@gmail.com", entities.User{Id: 1, Email: "test@gmail.com"}, nil},
	{"NotValidId_ReturnsErrorNotFound", "test1@gmail.com", entities.User{}, errors.New("user not found")},
}
