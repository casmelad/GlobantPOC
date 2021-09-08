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
		log.Fatalf("did not connect to server: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	defer conn.Close()

	c := NewUsersClient(conn)

	result, errorFromCall := c.GetAllUsers(ctx, &Void{})

	if errorFromCall != nil {
		log.Fatalf("server call did not work: %s", err)
	}

	response := []entities.User{}

	for _, o := range result.Users {
		response = append(response, entities.User{
			Id:       int(o.Id),
			EMail:    o.EMail,
			Name:     o.Name,
			LastName: o.LastName,
		})
	}

	return response, errorFromCall
}

func (up UserProxy) Create(u entities.User) (entities.User, error) {

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	defer conn.Close()

	c := NewUsersClient(conn)

	externalUser := &User{
		Id:       0,
		EMail:    u.EMail,
		Name:     u.Name,
		LastName: u.LastName,
	}

	result, errorFromCall := c.Create(ctx, externalUser)

	if result.Code != TaskResult_Ok {
		return entities.User{}, errors.New("error al crear")
	}

	u.Id = int(result.Result)

	return u, errorFromCall
}

func (up UserProxy) Update(u entities.User) (entities.User, error) {

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	defer conn.Close()

	c := NewUsersClient(conn)

	externalUser := &User{
		Id:       int32(u.Id),
		EMail:    u.EMail,
		Name:     u.Name,
		LastName: u.LastName,
	}

	result, errorFromCall := c.Update(ctx, externalUser)

	if result.Code != TaskResult_Ok || errorFromCall != nil {
		return entities.User{}, errorFromCall
	}

	return u, nil
}

func (up UserProxy) Delete(id int) (bool, error) {

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	defer conn.Close()

	c := NewUsersClient(conn)

	externalUserId := &UserId{
		Id: int32(id),
	}

	result, errorFromCall := c.Delete(ctx, externalUserId)

	if result.Code == TaskResult_Ok {
		return true, nil
	}

	return false, errorFromCall
}

func (up UserProxy) GetByEmail(email string) (entities.User, error) {

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	defer conn.Close()

	c := NewUsersClient(conn)

	userFromGrpc, errorFromCall := c.GetUser(ctx, &UserEmail{EMail: email})

	if errorFromCall != nil {
		fmt.Println("server call did not work:", errorFromCall)
		return entities.User{}, errorFromCall
	}

	response := entities.User{
		Id:       int(userFromGrpc.Id),
		EMail:    userFromGrpc.EMail,
		Name:     userFromGrpc.Name,
		LastName: userFromGrpc.LastName,
	}

	return response, errorFromCall
}
