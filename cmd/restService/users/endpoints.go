package users

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	PostUserEndpoint     endpoint.Endpoint
	PostManyUserEndpoint endpoint.Endpoint
	GetUserEndpoint      endpoint.Endpoint
	GetAllUsersEndpoint  endpoint.Endpoint
	PutUserEndpoint      endpoint.Endpoint
	DeleteUserEndpoint   endpoint.Endpoint
}

func MakeServerEndpoints(s UserProxy) Endpoints {
	return Endpoints{
		PostUserEndpoint:     MakePostUserEndpoint(s),
		PostManyUserEndpoint: MakePostUserEndpoint(s),
		GetUserEndpoint:      MakeGetUserEndpoint(s),
		GetAllUsersEndpoint:  MakeGetUserEndpoint(s),
		PutUserEndpoint:      MakePutUserEndpoint(s),
		DeleteUserEndpoint:   MakeDeleteUserEndpoint(s),
	}
}

// PostUser implements UsersProxy. Primarily useful in a client.
func (e Endpoints) PostUser(ctx context.Context, p User) error {
	request := postUserRequest{User: p}
	response, err := e.PostUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(postUserResponse)
	return resp.Err
}

// GetUser implements UsersProxy. Primarily useful in a client.
func (e Endpoints) GetUser(ctx context.Context, email string) (User, error) {
	request := getUserRequest{Email: email}
	response, err := e.GetUserEndpoint(ctx, request)
	if err != nil {
		return User{}, err
	}
	resp := response.(getUserResponse)
	return resp.User, resp.Err
}

// PutUser implements UsersProxy. Primarily useful in a client.
func (e Endpoints) PutUser(ctx context.Context, id string, p User) error {
	request := putUserRequest{User: p}
	response, err := e.PutUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(putUserResponse)
	return resp.Err
}

// DeleteUser implements UsersProxy. Primarily useful in a client.
func (e Endpoints) DeleteUser(ctx context.Context, id int) error {
	request := deleteUserRequest{UserID: id}
	response, err := e.DeleteUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteUserResponse)
	return resp.Err
}

// MakePostUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostUserEndpoint(s UserProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postUserRequest)
		usr, e := s.Create(req.User)
		return postUserResponse{User: usr, Err: e}, nil
	}
}

// MakeGetUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetUserEndpoint(s UserProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getUserRequest)
		p, e := s.GetByEmail(req.Email)
		return getUserResponse{User: p, Err: e}, nil
	}
}

// MakePutUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePutUserEndpoint(s UserProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putUserRequest)
		_, e := s.Update(req.User)
		return putUserResponse{Err: e}, nil
	}
}

// MakeDeleteUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteUserEndpoint(s UserProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteUserRequest)
		_, e := s.Delete(req.UserID)
		return deleteUserResponse{Err: e}, nil
	}
}

type postUserRequest struct {
	User User
}

type postUserResponse struct {
	Err  error `json:"err,omitempty"`
	User User  `json:"user,omitempty"`
}

type putUserRequest struct {
	User User
}

type putUserResponse struct {
	Err error `json:"err,omitempty"`
}

type deleteUserRequest struct {
	UserID int
}

type deleteUserResponse struct {
	Err error `json:"err,omitempty"`
}

type getUserRequest struct {
	Email string
}

type getUserResponse struct {
	Err  error `json:"err,omitempty"`
	User User  `json:"user,omitempty"`
}
