package infrastructure

import (
	"testing"

	"github.com/casmelad/GlobantPOC/pkg/domain/entities"
	"github.com/stretchr/testify/assert"
)

func Test_Add_ValidData_ReturnsNewId(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	userToAdd := entities.User{Email: "test@gmail.com"}
	//Act
	result := repository.Add(userToAdd)
	//Assert
	assert.Equal(t, 1, result)

}

func Test_Add_DuplicatedData_ReturnsInvalidResult(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	userToAdd := entities.User{Email: "test@gmail.com"}
	repository.Add(userToAdd)
	//Act
	result := repository.Add(userToAdd)
	//Assert
	assert.Equal(t, result, 0)

}

func Test_GetByEmail_ReturnsExistingData(t *testing.T) {
	//Arrange
	id := "test@gmail.com"
	repository := NewInMemoryUserRepository()
	userToAdd := entities.User{Email: id}
	expected := entities.User{Id: 1, Email: id}
	repository.Add(userToAdd)
	//Act
	result := repository.GetByEmail(id)
	//Assert
	assert.Equal(t, expected, result)
}

func Test_GetByEmail_InvalidId_ReturnsNoData(t *testing.T) {
	//Arrange
	id := "test@gmail.com"
	repository := NewInMemoryUserRepository()
	user := entities.User{Email: id}
	//Act
	result := repository.GetByEmail(id)
	//Assert
	assert.NotEqual(t, user.Email, result.Email)
}

func Test_GetByAll_ReturnsNoData(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	//Act
	result := repository.GetAll()
	//Assert
	assert.Equal(t, []entities.User{}, result)
}

func Test_GetByAll_ReturnsData(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	repository.Add(entities.User{Email: "test@gmail.com"})
	repository.Add(entities.User{Email: "test2@gmail.com"})
	//Act
	result := repository.GetAll()
	//Assert
	assert.Equal(t, 2, len(result))
}

func Test_Update_ValidData_UpdatesData(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	userToAdd := entities.User{Email: "test@gmail.com", Name: "Test1", LastName: "LastName1"}
	userToAdd2 := entities.User{Email: "test2@gmail.com", Name: "Test1", LastName: "LastName1"}
	newUserData := entities.User{Email: "test@gmail.com", Name: "Test1_Updated", LastName: "LastName1_Updated"}
	userId := repository.Add(userToAdd)
	repository.Add(userToAdd2)
	newUserData.Id = userId
	//Act
	repository.Update(newUserData)
	userUpdated := repository.GetByEmail(userToAdd.Email)
	//Assert
	assert.Equal(t, newUserData.Name, userUpdated.Name)
	assert.Equal(t, newUserData.LastName, userUpdated.LastName)

}

func Test_Update_InvalidData_ReturnsInvalidResult(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	newUserData := entities.User{Email: "test@gmail.com", Name: "Test1_Updated", LastName: "LastName1_Updated"}
	//Act
	result := repository.Update(newUserData)
	//Assert
	assert.Equal(t, 0, result)
}

func Test_Delete_ValidId_DeletesUser(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	userToAdd := entities.User{Email: "test@gmail.com", Name: "Test1", LastName: "LastName1"}
	userId := repository.Add(userToAdd)
	//Act
	result := repository.Delete(userId)
	//Assert
	assert.Equal(t, 1, result)
}

func Test_Delete_InvalidId_ReturnsInvalidResult(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	userToAdd := entities.User{Email: "test@gmail.com", Name: "Test1", LastName: "LastName1"}
	repository.Add(userToAdd)
	//Act
	result := repository.Delete(-1)
	//Assert
	assert.Equal(t, 0, result)
}
