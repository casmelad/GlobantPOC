package grpc

import (
	"context"
	"errors"
	"fmt"

	mappers "github.com/casmelad/GlobantPOC/cmd/grpcService/users/mappers"
	grpcUsers "github.com/casmelad/GlobantPOC/cmd/grpcService/users/proto"
	entities "github.com/casmelad/GlobantPOC/pkg/users"
	"github.com/go-kit/kit/endpoint"
)

type grpcUserServerEndpoints struct {
	CreateUserEndpoint     endpoint.Endpoint
	GetUserByEmailEndpoint endpoint.Endpoint
	UpdateUserEndpoint     endpoint.Endpoint
	DeleteUserEndpoint     endpoint.Endpoint
	GetAllUsersEndpoint    endpoint.Endpoint
}

func NewGrpcUsersServer(s entities.Service) *grpcUserServerEndpoints {
	return &grpcUserServerEndpoints{
		CreateUserEndpoint:     MakePostUserEndpoint(s),
		GetUserByEmailEndpoint: MakeGetUserEndpoint(s),
	}
}

func MakePostUserEndpoint(s entities.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		fmt.Println(ctx.Value("uuid"))

		reqData, validCast := request.(grpcUsers.CreateUserRequest)
		if !validCast {
			return nil, errors.New("invalid input data")
		}

		usr, err := mappers.ToDomainUser(*reqData.User)

		if err != nil {
			return nil, errors.New("invalid input data")
		}

		usrID, e := s.Create(ctx, usr)
		return postUserResponse{Id: usrID, Err: e}, nil
	}
}

func MakeGetUserEndpoint(s entities.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		reqData, validCast := request.(grpcUsers.EmailAddress)
		if !validCast {
			return nil, errors.New("invalid input data")
		}

		usr, err := s.GetByEmail(ctx, reqData.Value)

		return getUserResponse{User: usr}, err
	}
}

/* func MakePostUserEndpoint(s entities.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		reqData, validCast := request.(postUserRequest)
		if !validCast {
			return nil, errors.New("invalid input data")
		}

		usr, e := s.Create(ctx, reqData.User)
		return postUserResponse{Id: usr, err: e}, nil
	}
} */
