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

	serverCon, err := OpenServerConection()

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	defer serverCon.dispose()
	c := serverCon.client
	result, errorFromCall := c.GetAllUsers(serverCon.context, &Void{})

	if errorFromCall != nil {
		log.Fatalf("server call did not work: %s", err)
	}

	response := []entities.User{}

	for _, o := range result.Users {
		response = append(response, entities.User{
			Id:       int(o.Id),
			EMail:    o.Email,
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
	externalUser := &UserRequest{
		Id:       0,
		Email:    u.EMail,
		Name:     u.Name,
		LastName: u.LastName,
	}

	result, errorFromCall := c.Create(serverCon.context, externalUser)

	if result.Code != CreateUserResponse_OK {
		return entities.User{}, errors.New("error al crear")
	}

	fmt.Println(result.User)

	u.Id = int(result.User.Id)
	return u, errorFromCall
}

func (up UserProxy) Update(u entities.User) (entities.User, error) {

	serverCon, err := OpenServerConection()

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	defer serverCon.dispose()
	c := serverCon.client
	externalUser := &UserRequest{
		Id:       int32(u.Id),
		Email:    u.EMail,
		Name:     u.Name,
		LastName: u.LastName,
	}

	result, errorFromCall := c.Update(serverCon.context, externalUser)

	if result.Code != UpdateUserResponse_OK || errorFromCall != nil {
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
	externalUserId := &UserId{
		Id: int32(id),
	}
	result, errorFromCall := c.Delete(serverCon.context, externalUserId)

	if result.Code == DeleteUserResponse_OK {
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
	userFromGrpc, errorFromCall := c.GetUser(serverCon.context, &UserEmailRequest{EMail: email})

	if errorFromCall != nil {
		fmt.Println("server call did not work:", errorFromCall)
		return entities.User{}, errorFromCall
	}

	response := entities.User{
		Id:       int(userFromGrpc.Id),
		EMail:    userFromGrpc.Email,
		Name:     userFromGrpc.Name,
		LastName: userFromGrpc.LastName,
	}

	return response, errorFromCall
}

func OpenServerConection() (*ServerConnection, error) {

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
		return nil, err //unreached?
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	c := NewUsersClient(conn)

	return &ServerConnection{client: c, context: ctx, dispose: func() {
		cancel()
		conn.Close()

	}}, nil

}

type ServerConnection struct {
	client  UsersClient
	context context.Context
	dispose func()
}
