package users

import "fmt"

var (
	UsersBaseUri = "/users/"
	PostUser     = fmt.Sprintf("%s", UsersBaseUri)
	GetUser      = fmt.Sprintf("%s{%s}", UsersBaseUri, Email)
	PutUser      = GetUser
	DeleteUser   = fmt.Sprintf("%s{%s}", UsersBaseUri, UserID)
)
