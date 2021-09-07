package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	entities "github.com/casmelad/GlobantPOC/cmd/REST_server/entities"
	"google.golang.org/grpc"
)

type UserProxy struct {
}

func NewUserProxy() *UserProxy {
	return &UserProxy{}
}

func (up UserProxy) GetAll() ([]entities.User, error) {

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	defer conn.Close()

	c := NewUsersClient(conn)

	result, error2 := c.GetAllUsers(ctx, &Void{})

	fmt.Println(result)

	if error2 != nil {
		log.Fatalf("did not connect: %s", err)
	}

	response := []entities.User{}

	for _, o := range result.Users {
		response = append(response, entities.User{
			Id:       int(o.Id),
			Name:     o.Name,
			LastName: o.LastName,
		})
	}

	return response, error2
}

func (up UserProxy) Create(u entities.User) (entities.User, error) {

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	defer conn.Close()

	c := NewUsersClient(conn)

	externalUser := &User{
		Id:       0,
		Name:     u.Name,
		LastName: u.LastName,
	}

	result, error := c.Create(ctx, externalUser)

	if result.Code == 1 {
		return entities.User{}, errors.New("Error al crear")
	}

	return u, error
}
