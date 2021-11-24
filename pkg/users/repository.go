package users

import "context"

//Repository - repository interface for users
type Repository interface {
	//Add - adds a user to the repository
	Add(context.Context, User) (int, error)
	//GetByID - retrieves a user from the repository based on the integer id
	GetByID(context.Context, int) (User, error)
	//GetByEmail - retrieves a user from the repository based on the email address
	GetByEmail(context.Context, string) (User, error)
	//GetAll - retrieves all the users from the repository
	GetAll(context.Context) ([]User, error)
	//Update -  updates the information of a user
	Update(context.Context, User) error
	//Delete - deletes a user from the repository
	Delete(context.Context, int) error
}
