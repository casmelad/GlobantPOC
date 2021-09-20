package entities

type User struct {
	Id       int    `json:"id"`
	EMail    string `json:"email"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
}
