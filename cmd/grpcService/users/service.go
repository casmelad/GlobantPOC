package grpc

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/zipkin"

	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"

	mapper "github.com/casmelad/GlobantPOC/cmd/grpcService/users/mappers"
	proto "github.com/casmelad/GlobantPOC/cmd/grpcService/users/proto"
)

type grpcUserServer struct {
	proto.UsersServer
	getUser     grpctransport.Handler
	create      grpctransport.Handler
	getAllUsers grpctransport.Handler
	update      grpctransport.Handler
	delete      grpctransport.Handler
}

func NewGrpcUserServer(endpoints grpcUserServerEndpoints, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) proto.UsersServer {

	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	if zipkinTracer != nil {
		// Zipkin GRPC Server Trace can either be instantiated per gRPC method with a
		// provided operation name or a global tracing service can be instantiated
		// without an operation name and fed to each Go kit gRPC server as a
		// ServerOption.
		// In the latter case, the operation name will be the endpoint's grpc method
		// path if used in combination with the Go kit gRPC Interceptor.
		//
		// In this example, we demonstrate a global Zipkin tracing service with
		// Go kit gRPC Interceptor.
		options = append(options, zipkin.GRPCServerTrace(zipkinTracer))
	}

	server := &grpcUserServer{

		create:  grpctransport.NewServer(endpoints.CreateUserEndpoint, decodeCreateUserRequest, encodeCreateUserResponse, options...),
		getUser: grpctransport.NewServer(endpoints.GetUserByEmailEndpoint, decodeGetUserRequest, encodeGetUserResponse, options...),
	}

	return server
}

func (u grpcUserServer) GetUser(ctx context.Context, uid *proto.EmailAddress) (*proto.GetUserResponse, error) {

	_, usr, err := u.getUser.ServeGRPC(ctx, uid)

	if usr == nil {
		usr = getUserResponse{}
	}

	usrResp := usr.(getUserResponse)

	if err != nil {
		return nil, err
	}
	pbResponse, errDecode := encodeGetUserResponse(ctx, usrResp)

	if pbResponse == nil {
		pbResponse = proto.GetUserResponse{}
	}

	return pbResponse.(*proto.GetUserResponse), errDecode
}

func (u grpcUserServer) Create(ctx context.Context, user *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {

	_, errMapping := mapper.ToDomainUser(*user.User)

	if errMapping != nil {
		return nil, errMapping
	}

	_, newUserId, err := u.create.ServeGRPC(ctx, user)

	if err != nil {
		if err.Error() == "user already exists" {
			return &proto.CreateUserResponse{Code: proto.CodeResult_FAILED}, err
		} else {
			return &proto.CreateUserResponse{Code: proto.CodeResult_INVALIDINPUT}, err
		}
	}

	return &proto.CreateUserResponse{Code: proto.CodeResult_OK, UserId: int32(newUserId.(int))}, nil

}

func (u grpcUserServer) GetAllUsers(ctx context.Context, v *proto.Filters) (*proto.GetAllUsersResponse, error) {

	/* users, err := u.usersService.GetAll(ctx)
	response := []*proto.User{}

	if err != nil {
		return nil, err
	}

	for _, usr := range users {
		userMapped, errMapping := mapper.ToGrpcUser(usr)

		if errMapping != nil {
			return nil, errMapping
		}

		response = append(response, &userMapped)
	}
	*/
	return &proto.GetAllUsersResponse{}, nil
}

func (u grpcUserServer) Update(ctx context.Context, user *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {

	_, err := mapper.ToDomainUser(*user.User)

	if err != nil {
		return &proto.UpdateUserResponse{Code: proto.CodeResult_FAILED}, err
	}

	err_u := errors.New("")

	if err_u != nil {
		errorMessage := err_u.Error()
		switch errorMessage {
		case "user not found":
			return &proto.UpdateUserResponse{Code: proto.CodeResult_NOTFOUND}, err_u
		case "cannot update the user":
			return &proto.UpdateUserResponse{Code: proto.CodeResult_FAILED}, err_u
		default:
			return &proto.UpdateUserResponse{Code: proto.CodeResult_INVALIDINPUT}, err_u
		}
	}

	return &proto.UpdateUserResponse{Code: proto.CodeResult_OK}, nil
}

func (u grpcUserServer) Delete(ctx context.Context, userId *proto.Id) (*proto.DeleteUserResponse, error) {

	err := errors.New("") // u.usersService.Delete(ctx, int(userId.Value))

	if err != nil {
		errorMessage := err.Error()
		switch errorMessage {
		case "user not found":
			return &proto.DeleteUserResponse{Code: proto.CodeResult_NOTFOUND}, err
		case "invalid id":
			return &proto.DeleteUserResponse{Code: proto.CodeResult_INVALIDINPUT}, err
		default:
			return &proto.DeleteUserResponse{Code: proto.CodeResult_FAILED}, err
		}
	}

	return &proto.DeleteUserResponse{Code: proto.CodeResult_OK}, nil
}

// decodeGRPCSumRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC sum request to a user-domain sum request. Primarily useful in a server.
func decodeCreateUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	reqData, validCast := grpcReq.(*proto.CreateUserRequest)
	if !validCast {
		return nil, errors.New("invalid input data")
	}
	return postUserRequest{User: *reqData.User}, nil
}

func encodeCreateUserResponse(_ context.Context, resp interface{}) (interface{}, error) {
	reqData, validCast := resp.(postUserResponse)
	if !validCast {
		return nil, errors.New("invalid input data")
	}
	return &proto.CreateUserResponse{UserId: int32(reqData.Id), Code: proto.CodeResult_FAILED}, nil
}

// decodeGRPCSumRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC sum request to a user-domain sum request. Primarily useful in a server.
func decodeGetUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	reqData, validCast := grpcReq.(*proto.EmailAddress)
	if !validCast {
		return nil, errors.New("invalid input data")
	}
	return getUserRequest{Value: reqData.Value}, nil
}

func encodeGetUserResponse(_ context.Context, resp interface{}) (interface{}, error) {
	reqData, validCast := resp.(getUserResponse)
	if !validCast {
		return nil, errors.New("invalid input data to encode")
	}
	return &proto.GetUserResponse{User: &proto.User{Id: int32(reqData.User.ID), Email: reqData.User.Email, Name: reqData.User.Name, LastName: reqData.User.LastName}}, nil
}
