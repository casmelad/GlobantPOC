package grpc

import (
	"context"

	mapper "github.com/casmelad/GlobantPOC/cmd/grpcService/users/mappers"
	proto "github.com/casmelad/GlobantPOC/cmd/grpcService/users/proto"
	entities "github.com/casmelad/GlobantPOC/pkg/users"
)

type GrpcService struct {
	usersService entities.Service
}

func NewGrpcUserService(us entities.Service) *GrpcService {
	return &GrpcService{
		usersService: us,
	}
}

func (u *GrpcService) GetUser(ctx context.Context, uid *proto.EmailAddress) (*proto.GetUserResponse, error) {

	usr, err := u.usersService.GetByEmail(ctx, uid.Value)

	if err != nil {
		return nil, err
	}

	userToReturn, err := mapper.ToGrpcUser(usr)

	return &proto.GetUserResponse{User: &userToReturn}, err
}

func (u *GrpcService) Create(ctx context.Context, user *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {

	userToCreate, errMapping := mapper.ToDomainUser(*user.User)

	if errMapping != nil {
		return nil, errMapping
	}

	newUserId, err := u.usersService.Create(ctx, userToCreate)

	if err != nil {
		if err.Error() == "user already exists" {
			return &proto.CreateUserResponse{Code: proto.CodeResult_FAILED}, err
		} else {
			return &proto.CreateUserResponse{Code: proto.CodeResult_INVALIDINPUT}, err
		}
	}

	return &proto.CreateUserResponse{Code: proto.CodeResult_OK, UserId: int32(newUserId)}, nil

}

func (u *GrpcService) GetAllUsers(ctx context.Context, v *proto.Filters) (*proto.GetAllUsersResponse, error) {

	users, err := u.usersService.GetAll(ctx)
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

	return &proto.GetAllUsersResponse{Users: response}, nil
}

func (u *GrpcService) Update(ctx context.Context, user *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {

	userToUpdate, err := mapper.ToDomainUser(*user.User)

	if err != nil {
		return &proto.UpdateUserResponse{Code: proto.CodeResult_FAILED}, err
	}

	err_u := u.usersService.Update(ctx, userToUpdate)

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

func (u *GrpcService) Delete(ctx context.Context, userId *proto.Id) (*proto.DeleteUserResponse, error) {

	err := u.usersService.Delete(ctx, int(userId.Value))

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
