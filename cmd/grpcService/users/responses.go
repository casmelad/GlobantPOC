package grpc

type postUserResponse struct {
	Error error
	Id    int
}

type getUserResponse struct {
	User
}

type getAllUsersResponse struct {
	Users []User
}

type updateUserResponse struct {
	Error error
}

type deleteUserResponse struct {
	Error error
}
