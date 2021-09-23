package application

import (
	"context"
	"fmt"
	"log"
	"time"

	grpcService "github.com/casmelad/GlobantPOC/cmd/REST_server/application/grpcservices"
	entities "github.com/casmelad/GlobantPOC/cmd/REST_server/entities"
	"google.golang.org/grpc"
)

type UserProxy struct {
}

func NewUserProxy() *UserProxy {
	return &UserProxy{}
}

func (up UserProxy) GetAll() ([]entities.User, error) {

	serverCon, err := OpenServerConection()

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	defer serverCon.dispose()
	c := serverCon.client
	result, errorFromCall := c.GetAllUsers(serverCon.context, &grpcService.Filters{})

	if errorFromCall != nil {
		log.Fatalf("server call did not work: %s", err)
	}

	response := []entities.User{}

	for _, o := range result.Users {
		response = append(response, entities.User{
			Id:       int(o.Id),
			Email:    o.Email,
			Name:     o.Name,
			LastName: o.LastName,
		})
	}

	return response, errorFromCall
}

func (up UserProxy) Create(u entities.User) (entities.User, error) {

	serverCon, err := OpenServerConection()

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	defer serverCon.dispose()
	c := serverCon.client
	externalUser := &grpcService.User{
		Id:       0,
		Email:    u.Email,
		Name:     u.Name,
		LastName: u.LastName,
	}

	result, errorFromCall := c.Create(serverCon.context, &grpcService.CreateUserRequest{User: externalUser})

	if errorFromCall != nil {
		return entities.User{}, errorFromCall
	}

	u.Id = int(result.UserId)
	return u, errorFromCall
}

func (up UserProxy) Update(u entities.User) (entities.User, error) {

	serverCon, err := OpenServerConection()

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	defer serverCon.dispose()
	c := serverCon.client
	externalUser := grpcService.User{
		Id:       int32(u.Id),
		Email:    u.Email,
		Name:     u.Name,
		LastName: u.LastName,
	}

	_, errorFromCall := c.Update(serverCon.context, &grpcService.UpdateUserRequest{User: &externalUser})

	if errorFromCall != nil {
		return entities.User{}, errorFromCall
	}

	return u, nil
}

func (up UserProxy) Delete(id int) (bool, error) {

	serverCon, err := OpenServerConection()

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	defer serverCon.dispose()
	c := serverCon.client
	externalUserId := &grpcService.Id{
		Value: int32(id),
	}
	_, errorFromCall := c.Delete(serverCon.context, externalUserId)

	if errorFromCall != nil {
		return true, nil
	}

	return false, errorFromCall
}

func (up UserProxy) GetByEmail(email string) (entities.User, error) {

	serverCon, err := OpenServerConection()

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	defer serverCon.dispose()
	c := serverCon.client
	result, errorFromCall := c.GetUser(serverCon.context, &grpcService.EmailAddress{Value: email})

	if errorFromCall != nil {
		fmt.Println("server call did not work:", errorFromCall)
		return entities.User{}, errorFromCall
	}

	userFromGrpc := result.User

	response := entities.User{
		Id:       int(userFromGrpc.Id),
		Email:    userFromGrpc.Email,
		Name:     userFromGrpc.Name,
		LastName: userFromGrpc.LastName,
	}

	return response, nil
}

func OpenServerConection() (*ServerConnection, error) {

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
		return nil, err //unreached?
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	c := grpcService.NewUsersClient(conn)

	return &ServerConnection{
		client:  c,
		context: ctx,
		dispose: func() {
			cancel()
			conn.Close()
		},
	}, nil

}

type ServerConnection struct {
	client  grpcService.UsersClient
	context context.Context
	dispose func()
}
