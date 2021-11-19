package grpc

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/casmelad/GlobantPOC/pkg/users"
	"github.com/go-kit/kit/endpoint"
)

type grpcUserServerEndpoints struct {
	CreateUserEndpoint     endpoint.Endpoint
	GetUserByEmailEndpoint endpoint.Endpoint
	UpdateUserEndpoint     endpoint.Endpoint
	DeleteUserEndpoint     endpoint.Endpoint
	GetAllUsersEndpoint    endpoint.Endpoint
}

func NewGrpcUsersServer(s domain.Service) *grpcUserServerEndpoints {
	return &grpcUserServerEndpoints{
		CreateUserEndpoint:     MakePostUserEndpoint(s),
		GetUserByEmailEndpoint: MakeGetUserEndpoint(s),
		UpdateUserEndpoint:     MakeUpdateUserEndpoint(s),
		DeleteUserEndpoint:     MakeDeleteUserEndpoint(s),
		GetAllUsersEndpoint:    MakeGetAllUsersEndpoint(s),
	}
}

func MakePostUserEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqData, validCast := request.(postUserRequest)
		if !validCast {
			return nil, errors.New("invalid input data")
		}

		usr := domain.User{Email: reqData.Email, Name: reqData.Name, LastName: reqData.LastName}

		usrID, err := s.Create(ctx, usr)

		fmt.Println("Cool", postUserResponse{Id: usrID, Error: err})

		return postUserResponse{Id: usrID, Error: err}, nil
	}
}

func MakeGetUserEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		reqData, validCast := request.(getUserRequest)
		if !validCast {
			return nil, errors.New("invalid input data")
		}

		usr, err := s.GetByEmail(ctx, reqData.Email)

		return getUserResponse{User: User{Id: int32(usr.ID), Email: usr.Email, Name: usr.Name, LastName: usr.LastName}}, nil
	}
}

func MakeGetAllUsersEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		usrs, err := s.GetAll(ctx)

		responseData := getAllUsersResponse{Users: []User{}}

		for _, usr := range usrs {
			responseData.Users = append(responseData.Users, User{Id: int32(usr.ID), Email: usr.Email, Name: usr.Name, LastName: usr.LastName})
		}

		fmt.Println("users", responseData)

		return responseData, nil
	}
}

func MakeUpdateUserEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		fmt.Println("Llega a endpoint")
		reqData, validCast := request.(updateUserRequest)

		fmt.Println(reqData)

		if !validCast {
			return nil, errors.New("invalid request type")
		}

		usr := domain.User{Email: reqData.Email, Name: reqData.Name, LastName: reqData.LastName}

		if err != nil {
			return nil, errors.New("invalid object cast")
		}

		err = s.Update(ctx, usr)

		return updateUserResponse{Error: err}, nil
	}
}

func MakeDeleteUserEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqData, validCast := request.(deleteUserRequest)

		if !validCast {
			return nil, errors.New("invalid request type")
		}

		err = s.Delete(ctx, int(reqData.Id))

		return deleteUserResponse{Error: err}, nil
	}
}
