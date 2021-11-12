package users

type postUserResponse struct {
	Err  error `json:"err,omitempty"`
	User User  `json:"user,omitempty"`
}

type putUserResponse struct {
	Err error `json:"err,omitempty"`
}

type deleteUserResponse struct {
	Err error `json:"err,omitempty"`
}

type getUserResponse struct {
	Err  error `json:"err,omitempty"`
	User User  `json:"user,omitempty"`
}
