package mappers

import (
	grpc "github.com/casmelad/GlobantPOC/cmd/grpc_server/users"
	users "github.com/casmelad/GlobantPOC/pkg/users"
)

//ToDomainUser maps a grpc user to domain user
func ToDomainUser(userToMap grpc.User) (users.User, error) {
	return users.User{
		ID:       int(userToMap.Id),
		Email:    userToMap.Email,
		Name:     userToMap.Name,
		LastName: userToMap.LastName,
	}, nil
}

//ToGrpcUser maps a domain user to a grpc user
func ToGrpcUser(userToMap users.User) (grpc.User, error) {
	return grpc.User{
		Id:       int32(userToMap.ID),
		Email:    userToMap.Email,
		Name:     userToMap.Name,
		LastName: userToMap.LastName,
	}, nil

}
