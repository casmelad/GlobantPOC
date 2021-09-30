package repository

import (
	"context"
	"testing"

	"github.com/casmelad/GlobantPOC/pkg/users"
	"github.com/stretchr/testify/assert"
)

func Test_Add_ValidData_ReturnsNewId(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	userToAdd := users.User{Email: "test@gmail.com"}
	//Act
	result, err := repository.Add(context.Background(), userToAdd)
	//Assert
	assert.Equal(t, 1, result)
	assert.Nil(t, err)

}

func Test_Add_DuplicatedData_ReturnsInvalidResult(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	userToAdd := users.User{Email: "test@gmail.com"}
	ctx := context.Background()
	repository.Add(ctx, userToAdd)
	//Act
	result, err := repository.Add(ctx, userToAdd)
	//Assert
	assert.Equal(t, result, 0)
	assert.Nil(t, err)
}

func Test_GetByEmail_ReturnsExistingData(t *testing.T) {
	//Arrange
	id := "test@gmail.com"
	repository := NewInMemoryUserRepository()
	userToAdd := users.User{Email: id}
	expected := users.User{ID: 1, Email: id}
	ctx := context.Background()
	repository.Add(ctx, userToAdd)
	//Act
	result, err := repository.GetByEmail(ctx, id)
	//Assert
	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}

func Test_GetByEmail_InvalidId_ReturnsNoData(t *testing.T) {
	//Arrange
	id := "test@gmail.com"
	repository := NewInMemoryUserRepository()
	user := users.User{Email: id}
	//Act
	result, err := repository.GetByEmail(context.Background(), id)
	//Assert
	assert.NotEqual(t, user.Email, result.Email)
	assert.Nil(t, err)
}

func Test_GetByAll_ReturnsNoData(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	//Act
	result, err := repository.GetAll(context.Background())
	//Assert
	assert.Equal(t, []users.User{}, result)
	assert.Nil(t, err)
}

func Test_GetByAll_ReturnsData(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	repository.Add(context.Background(), users.User{Email: "test@gmail.com"})
	repository.Add(context.Background(), users.User{Email: "test2@gmail.com"})
	//Act
	result, err := repository.GetAll(context.Background())
	//Assert
	assert.Equal(t, 2, len(result))
	assert.Nil(t, err)
}

func Test_Update_ValidData_UpdatesData(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	userToAdd := users.User{Email: "test@gmail.com", Name: "Test1", LastName: "LastName1"}
	userToAdd2 := users.User{Email: "test2@gmail.com", Name: "Test1", LastName: "LastName1"}
	newUserData := users.User{Email: "test@gmail.com", Name: "Test1_Updated", LastName: "LastName1_Updated"}
	ctx := context.Background()
	userID, _ := repository.Add(ctx, userToAdd)
	repository.Add(ctx, userToAdd2)
	newUserData.ID = userID
	//Act
	errU := repository.Update(ctx, newUserData)
	userUpdated, _ := repository.GetByEmail(ctx, userToAdd.Email)
	//Assert
	assert.Equal(t, newUserData.Name, userUpdated.Name)
	assert.Equal(t, newUserData.LastName, userUpdated.LastName)
	assert.Nil(t, errU)

}

func Test_Update_InvalidData_ReturnsInvalidResult(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	newUserData := users.User{Email: "test@gmail.com", Name: "Test1_Updated", LastName: "LastName1_Updated"}
	//Act
	err := repository.Update(context.Background(), newUserData)
	//Assert
	assert.Nil(t, err)
}

func Test_Delete_ValidId_DeletesUser(t *testing.T) {
	//Arrange
	repository := NewInMemoryUserRepository()
	userToAdd := users.User{Email: "test@gmail.com", Name: "Test1", LastName: "LastName1"}
	ctx := context.Background()
	userID, _ := repository.Add(ctx, userToAdd)
	//Act
	err := repository.Delete(ctx, userID)
	//Assert
	assert.Nil(t, err)
}
