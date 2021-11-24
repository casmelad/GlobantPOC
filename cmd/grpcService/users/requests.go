package grpc

type postUserRequest struct {
	User `json:"user,omitempty"`
}

type updateUserRequest struct {
	User `json:"user,omitempty"`
}

type deleteUserRequest struct {
	Id int32 `json:"id,omitempty"`
}

type getUserRequest struct {
	Email string `json:"email,omitempty"`
}

type GetAllUsersRequest struct {
}

type User struct {
	//The user id to update
	Id int32 `json:"id,omitempty"`
	//The user email
	Email string `json:"email,omitempty"`
	//The user name
	Name string `json:"name,omitempty"`
	//The user last name
	LastName string `json:"last_name,omitempty"`
}
