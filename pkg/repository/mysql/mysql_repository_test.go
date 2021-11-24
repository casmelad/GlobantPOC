package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/casmelad/GlobantPOC/pkg/users"
	"github.com/stretchr/testify/assert"
)

var userID int = 0
var repository *MySQLRepository

func TestMain(m *testing.M) {
	var err error
	repository, err = NewMySQLUserRepository()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func Test_Add_ValidData_ReturnsNewId(t *testing.T) {
	//Arrange
	var err error
	userToAdd := users.User{Email: "test@gmail.com", Name: "Test", LastName: "LastName"}
	//Act
	userID, err = repository.Add(context.Background(), userToAdd)
	//Assert
	assert.NotZero(t, userID)
	assert.Nil(t, err)

}

func Test_Add_DuplicatedData_ReturnsInvalidResult(t *testing.T) {
	//Arrange
	userToAdd := users.User{Email: "test@gmail.com", Name: "Test", LastName: "LastName"}
	ctx := context.Background()
	repository.Add(ctx, userToAdd)
	//Act
	result, err := repository.Add(ctx, userToAdd)
	//Assert
	assert.Equal(t, result, 0)
	assert.NotNil(t, err)
}

func Test_GetByEmail_ReturnsExistingData(t *testing.T) {
	//Arrange
	emailAddress := "test@gmail.com"
	expected := users.User{ID: userID, Email: "test@gmail.com", Name: "Test", LastName: "LastName"}
	ctx := context.Background()
	//Act
	result, err := repository.GetByEmail(ctx, emailAddress)
	//Assert
	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}

func Test_GetByEmail_InvalidId_ReturnsNoData(t *testing.T) {
	//Arrange
	//Act
	result, err := repository.GetByEmail(context.Background(), "testMySql@gmail.com")
	//Assert
	assert.Empty(t, result)
	assert.Nil(t, err)
}

func Test_GetByAll_ReturnsData(t *testing.T) {
	//Arrange
	//Act
	result, err := repository.GetAll(context.Background())
	//Assert
	assert.NotZero(t, len(result))
	assert.Nil(t, err)
}

func Test_Update_ValidData_UpdatesData(t *testing.T) {
	//Arrange
	newUserData := users.User{ID: userID, Email: "test@gmail.com", Name: "Test1_Updated", LastName: "LastName1_Updated"}
	ctx := context.Background()
	newUserData.ID = userID
	//Act
	errU := repository.Update(ctx, newUserData)
	//Assert
	assert.Nil(t, errU)

}

/* func Test_Update_InvalidData_ReturnsInvalidResult(t *testing.T) {
	//Arrange
	newUserData := users.User{ID: userID}
	//Act
	err := repository.Update(context.Background(), newUserData)
	//Assert
	assert.NotNil(t, err)
} */

func Test_Delete_ValidId_DeletesUser(t *testing.T) {
	//Arrange
	ctx := context.Background()
	//Act
	err := repository.Delete(ctx, userID)
	//Assert
	assert.Nil(t, err)
}
