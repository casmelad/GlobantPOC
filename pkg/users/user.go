package users

//User - represents a user
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	LastName string `json:"lastname" validate:"required"`
}
