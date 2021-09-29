package mappers

import (
	"testing"

	"github.com/casmelad/GlobantPOC/cmd/grpc_server/users"
	entities "github.com/casmelad/GlobantPOC/pkg/users"
	"github.com/stretchr/testify/assert"
)

func Test_ToDomainUser_ResultOk(t *testing.T) {
	//Arrange
	toMap := users.User{Id: 999999}
	expectedResult := entities.User{ID: 999999}

	//Act
	result, err := ToDomainUser(toMap)

	//Assert
	assert.Equal(t, expectedResult, result)
	assert.Nil(t, err)
}

func Test_ToGrpcUser_ResultOk(t *testing.T) {

	//Arrange
	toMap := entities.User{ID: 999999}
	expectedResult := users.User{Id: 999999}

	//Act
	result, err := ToGrpcUser(toMap)

	//Assert
	assert.Equal(t, expectedResult, result)
	assert.Nil(t, err)
}
