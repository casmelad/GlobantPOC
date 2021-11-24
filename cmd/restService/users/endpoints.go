package users

import (
	"context"
	"errors"
	"fmt"

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

func MakeServerEndpoints(s GrpcUsersProxy) Endpoints {
	return Endpoints{
		PostUserEndpoint:     MakePostUserEndpoint(s),
		PostManyUserEndpoint: MakePostUserEndpoint(s),
		GetUserEndpoint:      MakeGetUserEndpoint(s),
		GetAllUsersEndpoint:  MakeGetAllUsersEndpoint(s),
		PutUserEndpoint:      MakePutUserEndpoint(s),
		DeleteUserEndpoint:   MakeDeleteUserEndpoint(s),
	}
}

// MakePostUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostUserEndpoint(s GrpcUsersProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		reqData, validCast := request.(postUserRequest)
		if !validCast {
			return nil, errors.New("invalid input data")
		}

		usr, e := s.Create(ctx, reqData.User)

		if e != nil {
			return WrapError(e), nil
		}

		return postUserResponse{Href: fmt.Sprintf("%s%s", UsersBaseUri, usr.Email), Err: e}, nil
	}
}

// MakeGetUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetUserEndpoint(s GrpcUsersProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		reqData, validCast := request.(getUserRequest)
		if !validCast {
			return nil, errors.New("invalid input data")
		}
		p, e := s.GetByEmail(ctx, reqData.Email) //pasar el context hasta el grpc

		if e != nil {
			return WrapError(e), nil
		}

		return getUserResponse{User: p, Err: e}, nil
	}
}

// MakeGetUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetAllUsersEndpoint(s GrpcUsersProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		_, validCast := request.(getAllUsersRequest)
		if !validCast {
			return nil, errors.New("invalid input data")
		}
		p, e := s.GetAll(ctx) //pasar el context hasta el grpc

		if e != nil {
			return WrapError(e), nil
		}

		return getAllUsersResponse{Users: p, Err: e}, nil
	}
}

// MakePutUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePutUserEndpoint(s GrpcUsersProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqData, validCast := request.(putUserRequest)
		if !validCast {
			return nil, errors.New("invalid input data")
		}
		_, e := s.Update(ctx, reqData.User)

		if e != nil {
			return WrapError(e), nil
		}

		return putUserResponse{Err: e}, nil
	}
}

// MakeDeleteUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteUserEndpoint(s GrpcUsersProxy) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqData, validCast := request.(deleteUserRequest)
		if !validCast {
			return nil, errors.New("invalid input data")
		}
		_, e := s.Delete(ctx, reqData.UserID)

		if e != nil {
			return WrapError(e), nil
		}

		return deleteUserResponse{Err: e}, nil
	}
}
