package grpc

import "github.com/casmelad/GlobantPOC/pkg/users"

type postUserResponse struct {
	Id  int
	Err error
}

type getUserResponse struct {
	User users.User
}
