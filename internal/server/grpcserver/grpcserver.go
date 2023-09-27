package grpcserver

import (
	"context"
	"fmt"
	"github.com/h2p2f/dedicated-vault/internal/server/models"
	pb "github.com/h2p2f/dedicated-vault/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

//go:generate mockery --name UserHandler --output ./mocks --filename mocks_userhandler.go
type UserHandler interface {
	Register(ctx context.Context, user models.User) (string, int64, error)
	Login(ctx context.Context, user models.User) (string, int64, error)
	GetUser(ctx context.Context, user string) (models.User, error)
	ChangePassword(ctx context.Context, user models.User, newPassword string) (string, error)
	DeleteUser(ctx context.Context, user models.User) error
}

//go:generate mockery --name DataHandler --output ./mocks --filename mocks_datahandler.go
type DataHandler interface {
	CreateData(ctx context.Context, user models.User, data models.VaultData) (string, int64, error)
	ChangeData(ctx context.Context, user models.User, data models.VaultData) (int64, error)
	GetAllData(ctx context.Context, user models.User) ([]models.VaultData, error)
	DeleteData(ctx context.Context, user models.User, data models.VaultData) (int64, error)
}

type VaultServer struct {
	pb.UnimplementedDedicatedVaultServer
	userHandler UserHandler
	dataHandler DataHandler
	logger      *zap.Logger
}

func NewVaultServer(uh UserHandler, dh DataHandler, logger *zap.Logger) *VaultServer {
	return &VaultServer{
		userHandler: uh,
		dataHandler: dh,
		logger:      logger}
}

func (s *VaultServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	if req.User.Name == "" || req.User.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "login or password is empty")
	}
	token, lastServerUpdated, err := s.userHandler.Register(ctx, models.User{
		Login:    req.User.Name,
		Password: req.User.Password,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := pb.RegisterResponse{
		Token:             token,
		LastServerUpdated: lastServerUpdated,
	}
	return &response, nil
}

func (s *VaultServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.User.Name == "" || req.User.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "login or password is empty")
	}

	token, lastServerUpdated, err := s.userHandler.Login(ctx, models.User{
		Login:    req.User.Name,
		Password: req.User.Password,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := pb.LoginResponse{
		Token:             token,
		LastServerUpdated: lastServerUpdated,
	}
	return &response, nil
}

func (s *VaultServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	if req.User.Name == "" || req.User.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "login or password is empty")
	}

	token, err := s.userHandler.ChangePassword(ctx, models.User{
		Login:    req.User.Name,
		Password: req.User.Password,
	}, req.NewPassword)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := pb.ChangePasswordResponse{
		Token: token,
	}
	return &response, nil
}

func (s *VaultServer) SaveSecret(ctx context.Context, req *pb.SaveSecretRequest) (*pb.SaveSecretResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "metadata is empty")
	}
	userFromContext := md.Get("user")
	fmt.Println(userFromContext)
	if len(userFromContext) == 0 {
		return nil, status.Error(codes.InvalidArgument, "user is empty")
	}
	user, err := s.userHandler.GetUser(ctx, userFromContext[0])
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	dataUUID, created, err := s.dataHandler.CreateData(ctx, user, models.VaultData{
		Meta:     req.Data.Meta,
		DataType: req.Data.Type,
		Data:     req.Data.Value,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := pb.SaveSecretResponse{
		Uuid:              dataUUID,
		Created:           created,
		LastServerUpdated: created,
	}
	return &response, nil
}

func (s *VaultServer) ChangeSecret(ctx context.Context, req *pb.ChangeSecretRequest) (*pb.ChangeSecretResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "metadata is empty")
	}
	userFromContext := md.Get("user")
	if len(userFromContext) == 0 {
		return nil, status.Error(codes.InvalidArgument, "user is empty")
	}
	user, err := s.userHandler.GetUser(ctx, userFromContext[0])
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	_ = user
	updated, err := s.dataHandler.ChangeData(ctx, user, models.VaultData{
		DataUUID: req.Data.Uuid,
		Meta:     req.Data.Meta,
		DataType: req.Data.Type,
		Data:     req.Data.Value,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := pb.ChangeSecretResponse{
		Updated:           updated,
		LastServerUpdated: updated,
	}
	return &response, nil
}

func (s *VaultServer) DeleteSecret(ctx context.Context, req *pb.DeleteSecretRequest) (*pb.DeleteSecretResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "metadata is empty")
	}
	userFromContext := md.Get("user")
	if len(userFromContext) == 0 {
		return nil, status.Error(codes.InvalidArgument, "user is empty")
	}
	user, err := s.userHandler.GetUser(ctx, userFromContext[0])
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	lastServerUpdated, err := s.dataHandler.DeleteData(ctx, user, models.VaultData{
		DataUUID: req.Uuid,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := pb.DeleteSecretResponse{
		Uuid:              req.Uuid,
		LastServerUpdated: lastServerUpdated,
	}
	return &response, nil
}

func (s *VaultServer) ListSecrets(ctx context.Context, req *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "metadata is empty")
	}
	userFromContext := md.Get("user")
	if len(userFromContext) == 0 {
		return nil, status.Error(codes.InvalidArgument, "user is empty")
	}
	user, err := s.userHandler.GetUser(ctx, userFromContext[0])
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	s.logger.Info("user", zap.Any("user", user))
	data, err := s.dataHandler.GetAllData(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var response pb.ListSecretsResponse
	response.LastServerUpdated = user.LastServerUpdated
	for _, d := range data {
		s.logger.Info("data", zap.Any("data", d))
		response.Data = append(response.Data, &pb.SecretData{
			Uuid:  d.DataUUID,
			Meta:  d.Meta,
			Type:  d.DataType,
			Value: d.Data,
		})
	}
	return &response, nil
}
