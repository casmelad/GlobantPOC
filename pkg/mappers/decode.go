package mappers

import (
	"github.com/casmelad/GlobantPOC/cmd/grpc_server/users"
	"github.com/casmelad/GlobantPOC/pkg/domain/entities"
)

func ToDomainUser(userToMap users.User) (entities.User, error) {
	return entities.User{
		Id:       int(userToMap.Id),
		Email:    userToMap.Email,
		Name:     userToMap.Name,
		LastName: userToMap.LastName,
	}, nil
}

func ToGrpcUser(userToMap entities.User) (users.User, error) {
	return users.User{
		Id:       int32(userToMap.Id),
		Email:    userToMap.Email,
		Name:     userToMap.Name,
		LastName: userToMap.LastName,
	}, nil

}
