package mappers

import (
	proto "github.com/casmelad/GlobantPOC/cmd/grpcService/users/proto"
	domain "github.com/casmelad/GlobantPOC/pkg/users"
)

//ToDomainUser maps a grpc user to domain user
func ToDomainUser(userToMap proto.User) (domain.User, error) {
	return domain.User{
		ID:       int(userToMap.Id),
		Email:    userToMap.Email,
		Name:     userToMap.Name,
		LastName: userToMap.LastName,
	}, nil
}

//ToGrpcUser maps a domain user to a grpc user
func ToGrpcUser(userToMap domain.User) (proto.User, error) {
	return proto.User{
		Id:       int32(userToMap.ID),
		Email:    userToMap.Email,
		Name:     userToMap.Name,
		LastName: userToMap.LastName,
	}, nil

}
