package users

type postUserRequest struct {
	User User
}

type putUserRequest struct {
	User User
}

type deleteUserRequest struct {
	UserID int
}

type getUserRequest struct {
	Email string
}
