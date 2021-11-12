package users

import (
	"context"
	"errors"

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
	resp, validCast := response.(postUserResponse)

	if !validCast {
		return err
	}

	return resp.Err
}

// GetUser implements UsersProxy. Primarily useful in a client.
func (e Endpoints) GetUser(ctx context.Context, email string) (User, error) {

	request := getUserRequest{Email: email}
	response, err := e.GetUserEndpoint(ctx, request)
	if err != nil {
		return User{}, err
	}
	resp, validCast := response.(getUserResponse)

	if !validCast {
		return User{}, errors.New("incompatible endpoint response type and internal handler response type")
	}

	return resp.User, resp.Err
}

// PutUser implements UsersProxy. Primarily useful in a client.
func (e Endpoints) PutUser(ctx context.Context, id string, p User) error {
	request := putUserRequest{User: p}
	response, err := e.PutUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp, validCast := response.(putUserResponse)

	if !validCast {
		return errors.New("incompatible endpoint response type and internal handler response type")
	}

	return resp.Err
}

// DeleteUser implements UsersProxy. Primarily useful in a client.
func (e Endpoints) DeleteUser(ctx context.Context, id int) error {

	request := deleteUserRequest{UserID: id}
	response, err := e.DeleteUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp, validCast := response.(deleteUserResponse)

	if !validCast {
		return errors.New("incompatible endpoint response type and internal handler response type")
	}

	return resp.Err
}

// MakePostUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostUserEndpoint(s UserProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		reqData, validCast := request.(postUserRequest)
		if !validCast {
			return nil, errors.New("invalid input data")
		}

		usr, e := s.Create(ctx, reqData.User)
		return postUserResponse{User: usr, Err: e}, nil
	}
}

// MakeGetUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetUserEndpoint(s UserProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		reqData, validCast := request.(getUserRequest)
		if !validCast {
			return nil, errors.New("invalid input data")
		}
		p, e := s.GetByEmail(ctx, reqData.Email) //pasar el context hasta el grpc
		return getUserResponse{User: p, Err: e}, nil
	}
}

// MakePutUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePutUserEndpoint(s UserProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqData, validCast := request.(putUserRequest)
		if !validCast {
			return nil, errors.New("invalid input data")
		}
		_, e := s.Update(ctx, reqData.User)
		return putUserResponse{Err: e}, nil
	}
}

// MakeDeleteUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteUserEndpoint(s UserProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqData, validCast := request.(deleteUserRequest)
		if !validCast {
			return nil, errors.New("invalid input data")
		}
		_, e := s.Delete(ctx, reqData.UserID)
		return deleteUserResponse{Err: e}, nil
	}
}
