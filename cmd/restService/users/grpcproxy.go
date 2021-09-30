package users

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"

	proto "github.com/casmelad/GlobantPOC/cmd/restService/users/proto"
	"google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"
)

type UserProxy struct {
	grpcLog glog.LoggerV2
}

type config struct {
	Port int    `env:"proto_PORT" envDefault:"9000"`
	Host string `env:"proto_HOST" envDefault:"127.0.0.1"`
}

func NewUserProxy() *UserProxy {
	return &UserProxy{
		grpcLog: glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout),
	}
}

func (up UserProxy) GetAll() ([]User, error) {

	serverCon, err := OpenServerConection(up)

	if err != nil {
		log.Fatalf(err.Error())
	}

	defer serverCon.dispose()
	c := serverCon.client
	result, errorFromCall := c.GetAllUsers(serverCon.context, &proto.Filters{})

	if errorFromCall != nil {
		log.Fatalf(errorFromCall.Error())
	}

	response := []User{}

	for _, o := range result.Users {
		response = append(response, User{
			Id:       int(o.Id),
			Email:    o.Email,
			Name:     o.Name,
			LastName: o.LastName,
		})
	}

	return response, errorFromCall
}

func (up UserProxy) Create(u User) (User, error) {

	serverCon, err := OpenServerConection(up)

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	defer serverCon.dispose()
	c := serverCon.client
	externalUser := &proto.User{
		Id:       0,
		Email:    u.Email,
		Name:     u.Name,
		LastName: u.LastName,
	}

	result, errorFromCall := c.Create(serverCon.context, &proto.CreateUserRequest{User: externalUser})

	if errorFromCall != nil {
		return User{}, errorFromCall
	}

	u.Id = int(result.UserId)
	return u, errorFromCall
}

func (up UserProxy) Update(u User) (User, error) {

	serverCon, err := OpenServerConection(up)

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	defer serverCon.dispose()
	c := serverCon.client
	externalUser := proto.User{
		Id:       int32(u.Id),
		Email:    u.Email,
		Name:     u.Name,
		LastName: u.LastName,
	}

	_, errorFromCall := c.Update(serverCon.context, &proto.UpdateUserRequest{User: &externalUser})

	if errorFromCall != nil {
		return User{}, errorFromCall
	}

	return u, nil
}

func (up UserProxy) Delete(id int) (bool, error) {

	serverCon, err := OpenServerConection(up)

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	fmt.Println(id)

	defer serverCon.dispose()
	c := serverCon.client
	externalUserId := &proto.Id{
		Value: int32(id),
	}
	_, errorFromCall := c.Delete(serverCon.context, externalUserId)

	if errorFromCall != nil {
		return true, errorFromCall
	}

	return false, errorFromCall
}

func (up UserProxy) GetByEmail(email string) (User, error) {

	serverCon, err := OpenServerConection(up)

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
	}

	defer serverCon.dispose()
	c := serverCon.client
	result, errorFromCall := c.GetUser(serverCon.context, &proto.EmailAddress{Value: email})

	if errorFromCall != nil {
		fmt.Println("server call did not work:", errorFromCall)
		return User{}, errorFromCall
	}

	userFromGrpc := result.User

	response := User{
		Id:       int(userFromGrpc.Id),
		Email:    userFromGrpc.Email,
		Name:     userFromGrpc.Name,
		LastName: userFromGrpc.LastName,
	}

	return response, nil
}

func OpenServerConection(up UserProxy) (*ServerConnection, error) {

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Host, strconv.Itoa(cfg.Port)), grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect to server: %s", err)
		return nil, err //unreached?
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	c := proto.NewUsersClient(conn)

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
	client  proto.UsersClient
	context context.Context
	dispose func()
}
