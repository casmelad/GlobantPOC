package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/zipkin"

	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"

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

		create:      grpctransport.NewServer(endpoints.CreateUserEndpoint, decodeCreateUserRequest, encodeCreateUserResponse, options...),
		getUser:     grpctransport.NewServer(endpoints.GetUserByEmailEndpoint, decodeGetUserRequest, encodeGetUserResponse, options...),
		getAllUsers: grpctransport.NewServer(endpoints.GetAllUsersEndpoint, decodeGetAllUsersRequest, encodeGetAllUsersResponse, options...),
		update:      grpctransport.NewServer(endpoints.UpdateUserEndpoint, decodeUpdateUserRequest, encodeUpdateUserResponse, options...),
		delete:      grpctransport.NewServer(endpoints.DeleteUserEndpoint, decodeDeleteUserRequest, encodeDeleteUserResponse, options...),
	}

	return server
}

func (u grpcUserServer) GetUser(ctx context.Context, uid *proto.EmailAddress) (*proto.GetUserResponse, error) {

	_, grpcResponse, err := u.getUser.ServeGRPC(ctx, uid)

	return grpcResponse.(*proto.GetUserResponse), err
}

func (u grpcUserServer) Create(ctx context.Context, user *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {

	_, grpcResponse, err := u.create.ServeGRPC(ctx, user)

	fmt.Println("xvbxbbcbcvbc", grpcResponse)

	return grpcResponse.(*proto.CreateUserResponse), err

}

func (u grpcUserServer) GetAllUsers(ctx context.Context, filters *proto.Filters) (*proto.GetAllUsersResponse, error) {

	_, grpcResponse, err := u.getAllUsers.ServeGRPC(ctx, filters)

	return grpcResponse.(*proto.GetAllUsersResponse), err
}

func (u grpcUserServer) Update(ctx context.Context, userInfo *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {

	fmt.Println("Update method", u.update)

	_, grpcResponse, err := u.update.ServeGRPC(ctx, userInfo)

	return grpcResponse.(*proto.UpdateUserResponse), err
}

func (u grpcUserServer) Delete(ctx context.Context, userId *proto.Id) (*proto.DeleteUserResponse, error) {

	ctx, grpcResponse, err := u.delete.ServeGRPC(ctx, userId)

	return grpcResponse.(*proto.DeleteUserResponse), err
}

// decodeGRPCSumRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC sum request to a user-domain sum request. Primarily useful in a server.
func decodeCreateUserRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	reqData, validCast := grpcReq.(*proto.CreateUserRequest)
	if !validCast {
		return nil, errors.New("invalid input data")
	}
	usr := User{Email: reqData.User.Email, Name: reqData.User.Name, LastName: reqData.User.LastName}

	return postUserRequest{User: usr}, nil
}

func encodeCreateUserResponse(ctx context.Context, resp interface{}) (interface{}, error) {

	respData, validCast := resp.(postUserResponse)

	if !validCast {
		return nil, errors.New("invalid input data")
	}

	if respData.Error != nil {
		if respData.Error.Error() == "user already exists" {
			return &proto.CreateUserResponse{Code: proto.CodeResult_FAILED}, nil
		} else {
			return &proto.CreateUserResponse{Code: proto.CodeResult_INVALIDINPUT}, nil
		}
	}

	return &proto.CreateUserResponse{UserId: int32(respData.Id), Code: proto.CodeResult_OK}, nil
}

// decodeGRPCSumRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC sum request to a user-domain sum request. Primarily useful in a server.
func decodeGetUserRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	reqData, validCast := grpcReq.(*proto.EmailAddress)
	if !validCast {
		return nil, errors.New("invalid input data")
	}
	return getUserRequest{Email: reqData.Value}, nil
}

func encodeGetUserResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	respData, validCast := resp.(getUserResponse)
	if !validCast {
		return nil, errors.New("invalid input data to encode")
	}
	usr := proto.User{Id: respData.Id, Name: respData.Name, Email: respData.Email, LastName: respData.LastName}

	return &proto.GetUserResponse{User: &usr}, nil
}

// decodeGetAllUsersRequest : param Filters is not being used yet
func decodeGetAllUsersRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	_, validCast := grpcReq.(*proto.Filters)
	if !validCast {
		return nil, errors.New("invalid input data decode")
	}
	return getUserRequest{}, nil
}

func encodeGetAllUsersResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	respData, validCast := resp.(getAllUsersResponse)
	if !validCast {
		return nil, errors.New("invalid input data to encode")
	}

	response := &proto.GetAllUsersResponse{Users: []*proto.User{}}

	for _, usr := range respData.Users {
		pbUser := proto.User{Id: usr.Id, Name: usr.Name, Email: usr.Email, LastName: usr.LastName}

		response.Users = append(response.Users, &pbUser)
	}

	return response, nil
}

func decodeUpdateUserRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {

	reqData, validCast := grpcReq.(*proto.UpdateUserRequest)

	fmt.Println("decode", reqData)

	if !validCast {
		return nil, errors.New("invalid input data to decode")
	}

	usr := User{Id: reqData.User.Id, Email: reqData.User.Email, Name: reqData.User.Name, LastName: reqData.User.LastName}

	fmt.Println("user", usr)

	return updateUserRequest{User: usr}, nil
}

func encodeUpdateUserResponse(ctx context.Context, resp interface{}) (interface{}, error) {

	fmt.Println("encode", resp)
	respData, validCast := resp.(updateUserResponse)

	fmt.Println("encode 2", respData)

	if !validCast {
		return nil, errors.New("invalid input data to encode")
	}

	if respData.Error != nil {
		switch respData.Error.Error() {
		case "user not found":
			return &proto.UpdateUserResponse{Code: proto.CodeResult_NOTFOUND}, nil
		case "cannot update the user information":
			return &proto.UpdateUserResponse{Code: proto.CodeResult_FAILED}, nil
		default:
			return &proto.UpdateUserResponse{Code: proto.CodeResult_INVALIDINPUT}, nil
		}
	}

	return &proto.UpdateUserResponse{Code: proto.CodeResult_OK}, nil
}

func decodeDeleteUserRequest(ctx context.Context, req interface{}) (interface{}, error) {
	fmt.Println(req)
	reqData, validCast := req.(*proto.Id)

	if !validCast {
		return nil, errors.New("invalid input data to decode")
	}

	return deleteUserRequest{Id: reqData.Value}, nil
}

func encodeDeleteUserResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	fmt.Println(resp)
	respData, validCast := resp.(deleteUserResponse)

	if !validCast {
		return nil, errors.New("invalid input data to encode")
	}

	if respData.Error != nil {
		errorMessage := respData.Error.Error()
		switch errorMessage {
		case "user not found":
			return &proto.DeleteUserResponse{Code: proto.CodeResult_NOTFOUND}, nil
		case "invalid id":
			return &proto.DeleteUserResponse{Code: proto.CodeResult_INVALIDINPUT}, nil
		default:
			return &proto.DeleteUserResponse{Code: proto.CodeResult_FAILED}, nil
		}
	}

	return &proto.DeleteUserResponse{Code: proto.CodeResult_OK}, nil
}
