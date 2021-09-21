package mappers

import (
	"testing"

	"github.com/casmelad/GlobantPOC/cmd/grpc_server/users"
	"github.com/casmelad/GlobantPOC/pkg/domain/entities"
	"github.com/stretchr/testify/assert"
)

func Test_ToDomainUser_ResultOk(t *testing.T) {
	//Arrange
	toMap := users.User{Id: 999999}
	expectedResult := entities.User{Id: 999999}

	//Act
	result, err := ToDomainUser(toMap)

	//Assert
	assert.Equal(t, expectedResult, result)
	assert.Nil(t, err)
}

func Test_ToGrpcUser_ResultOk(t *testing.T) {

	//Arrange
	toMap := entities.User{Id: 999999}
	expectedResult := users.User{Id: 999999}

	//Act
	result, err := ToGrpcUser(toMap)

	//Assert
	assert.Equal(t, expectedResult, result)
	assert.Nil(t, err)
}
