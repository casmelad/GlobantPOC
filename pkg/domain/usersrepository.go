package domain

import "github.com/casmelad/GlobantPOC/pkg/domain/entities"

type UsersRepositoryInterface interface {
	Add(entities.User) int
	GetById(int) entities.User
	GetByEmail(string) entities.User
	GetAll() []entities.User
	Update(entities.User) int
	Delete(int) int
}
