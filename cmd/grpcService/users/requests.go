package grpc

import users "github.com/casmelad/GlobantPOC/cmd/grpcService/users/proto"

type postUserRequest struct {
	users.User
}

type updateUserRequest struct {
	users.User
}

type deleteUserRequest struct {
	Value int
}

type getUserRequest struct {
	Value string
}

type GetAllUsersRequest struct {
}
