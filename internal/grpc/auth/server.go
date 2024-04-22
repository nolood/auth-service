package auth

import (
	"context"
	ssov1 "github.com/nolood/auth-service-protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appId int) (token string, err error)
	RegisterNewUser(ctx context.Context, req *ssov1.RegisterRequest) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {

	// TODO: user validator
	if req.GetEmail() == "" || req.GetPassword() == "" || req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "required-fields")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		// TODO: internal error
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "required-fields")
	}

	userID, err := s.auth.RegisterNewUser(ctx, req)
	if err != nil {
		// TODO: internal error
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ssov1.RegisterResponse{UserId: userID}, nil

}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {

	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "required-fields")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		// TODO: internal error
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil

}
