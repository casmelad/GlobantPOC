package mappers

import (
	"testing"

	proto "github.com/casmelad/GlobantPOC/cmd/grpcService/users/proto"
	domain "github.com/casmelad/GlobantPOC/pkg/users"
	"github.com/stretchr/testify/assert"
)

func Test_ToDomainUser_ResultOk(t *testing.T) {
	//Arrange
	toMap := proto.User{Id: 999999}
	expectedResult := domain.User{ID: 999999}

	//Act
	result, err := ToDomainUser(toMap)

	//Assert
	assert.Equal(t, expectedResult, result)
	assert.Nil(t, err)
}

func Test_ToGrpcUser_ResultOk(t *testing.T) {

	//Arrange
	toMap := domain.User{ID: 999999}
	expectedResult := proto.User{Id: 999999}

	//Act
	result, err := ToGrpcUser(toMap)

	//Assert
	assert.Equal(t, expectedResult, result)
	assert.Nil(t, err)
}
