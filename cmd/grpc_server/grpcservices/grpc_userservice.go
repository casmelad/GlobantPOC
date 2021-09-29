package grpcservices

import (
	"context"

	pb "github.com/casmelad/GlobantPOC/cmd/grpc_server/users"
	mapper "github.com/casmelad/GlobantPOC/pkg/mappers"
	users "github.com/casmelad/GlobantPOC/pkg/users"
)

type GrpcUserService struct {
	usersService users.Service
}

func NewGrpcUserService(us users.Service) *GrpcUserService {
	return &GrpcUserService{
		usersService: us,
	}
}

func (u *GrpcUserService) GetUser(ctx context.Context, uid *pb.EmailAddress) (*pb.GetUserResponse, error) {

	usr, err := u.usersService.GetByEmail(ctx, uid.Value)

	if err != nil {
		return nil, err
	}

	userToReturn, err := mapper.ToGrpcUser(usr)

	return &pb.GetUserResponse{User: &userToReturn}, err
}

func (u *GrpcUserService) Create(ctx context.Context, user *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	userToCreate, errMapping := mapper.ToDomainUser(*user.User)

	if errMapping != nil {
		return nil, errMapping
	}

	newUserId, err := u.usersService.Create(ctx, userToCreate)

	if err != nil {
		if err.Error() == "user already exists" {
			return &pb.CreateUserResponse{Code: pb.CodeResult_FAILED}, err
		} else {
			return &pb.CreateUserResponse{Code: pb.CodeResult_INVALIDINPUT}, err
		}
	}

	return &pb.CreateUserResponse{Code: pb.CodeResult_OK, UserId: int32(newUserId)}, nil

}

func (u *GrpcUserService) GetAllUsers(ctx context.Context, v *pb.Filters) (*pb.GetAllUsersResponse, error) {

	users, err := u.usersService.GetAll(ctx)
	response := []*pb.User{}

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

	return &pb.GetAllUsersResponse{Users: response}, nil
}

func (u *GrpcUserService) Update(ctx context.Context, user *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	userToUpdate, err := mapper.ToDomainUser(*user.User)

	if err != nil {
		return &pb.UpdateUserResponse{Code: pb.CodeResult_FAILED}, err
	}

	err_u := u.usersService.Update(ctx, userToUpdate)

	if err_u != nil {
		errorMessage := err_u.Error()
		switch errorMessage {
		case "user not found":
			return &pb.UpdateUserResponse{Code: pb.CodeResult_NOTFOUND}, err_u
		case "cannot update the user":
			return &pb.UpdateUserResponse{Code: pb.CodeResult_FAILED}, err_u
		default:
			return &pb.UpdateUserResponse{Code: pb.CodeResult_INVALIDINPUT}, err_u
		}
	}

	return &pb.UpdateUserResponse{Code: pb.CodeResult_OK}, nil
}

func (u *GrpcUserService) Delete(ctx context.Context, userId *pb.Id) (*pb.DeleteUserResponse, error) {

	err := u.usersService.Delete(ctx, int(userId.Value))

	if err != nil {
		errorMessage := err.Error()
		switch errorMessage {
		case "user not found":
			return &pb.DeleteUserResponse{Code: pb.CodeResult_NOTFOUND}, err
		case "invalid id":
			return &pb.DeleteUserResponse{Code: pb.CodeResult_INVALIDINPUT}, err
		default:
			return &pb.DeleteUserResponse{Code: pb.CodeResult_FAILED}, err
		}
	}

	return &pb.DeleteUserResponse{Code: pb.CodeResult_OK}, nil
}
