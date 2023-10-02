package grpcserver

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/h2p2f/dedicated-vault/internal/server/grpcserver/mocks"
	"github.com/h2p2f/dedicated-vault/internal/server/models"
	pb "github.com/h2p2f/dedicated-vault/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestVaultServer_Register(t *testing.T) {
	mockCtx := context.Background()

	tests := []struct {
		testname string
		name     string
		password string
		wantCode codes.Code
	}{
		{
			testname: "valid",
			name:     "testuser",
			password: "testpassword",
			wantCode: codes.OK,
		},
		{
			testname: "empty name",
			name:     "",
			password: "testpassword",
			wantCode: codes.InvalidArgument,
		},
		{
			testname: "empty password",
			name:     "testuser",
			password: "",
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			server := &VaultServer{
				userHandler: &mocks.UserHandler{},
			}
			mockToken := "mocktoken"
			mockLastServerUpdated := time.Now().Unix()

			if tt.wantCode == codes.OK {
				mockUserHandler := &mocks.UserHandler{}
				mockUserHandler.On("Register", mockCtx, models.User{
					Login:    tt.name,
					Password: tt.password,
				}).Return(mockToken, mockLastServerUpdated, nil)
				server.userHandler = mockUserHandler
			} else {
				mockUserHandler := &mocks.UserHandler{}
				mockUserHandler.On("Register", mockCtx, models.User{
					Login:    tt.name,
					Password: tt.password,
				}).Return("", int64(0), errors.New("error"))
				server.userHandler = mockUserHandler
			}
			req := &pb.RegisterRequest{
				User: &pb.User{
					Name:     tt.name,
					Password: tt.password,
				},
			}
			_, err := server.Register(mockCtx, req)
			assert.Equal(t, tt.wantCode, status.Code(err))
		})
	}
}

func TestVaultServer_Login(t *testing.T) {
	mockCtx := context.Background()

	tests := []struct {
		testname string
		name     string
		password string
		wantCode codes.Code
	}{
		{
			testname: "valid",
			name:     "testuser",
			password: "testpassword",
			wantCode: codes.OK,
		},
		{
			testname: "empty name",
			name:     "",
			password: "testpassword",
			wantCode: codes.InvalidArgument,
		},
		{
			testname: "empty password",
			name:     "testuser",
			password: "",
			wantCode: codes.InvalidArgument,
		},
		{
			testname: "invalid user",
			name:     "testuser",
			password: "testpassword",
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			server := &VaultServer{
				userHandler: &mocks.UserHandler{},
			}
			mockToken := "mocktoken"
			mockLastServerUpdated := time.Now().Unix()

			if tt.wantCode == codes.OK {
				mockUserHandler := &mocks.UserHandler{}
				mockUserHandler.On("Login", mockCtx, models.User{
					Login:    tt.name,
					Password: tt.password,
				}).Return(mockToken, mockLastServerUpdated, nil)
				server.userHandler = mockUserHandler
			} else {
				mockUserHandler := &mocks.UserHandler{}
				mockUserHandler.On("Login", mockCtx, models.User{
					Login:    tt.name,
					Password: tt.password,
				}).Return("", int64(0), errors.New("error"))
				server.userHandler = mockUserHandler
			}
			req := &pb.LoginRequest{
				User: &pb.User{
					Name:     tt.name,
					Password: tt.password,
				},
			}
			_, err := server.Login(mockCtx, req)
			assert.Equal(t, tt.wantCode, status.Code(err))
		})
	}
}

func TestVaultServer_ChangePassword(t *testing.T) {
	mockCtx := context.Background()
	tests := []struct {
		testname    string
		name        string
		password    string
		newPassword string
		wantCode    codes.Code
	}{
		{
			testname:    "valid",
			name:        "testuser",
			password:    "testpassword",
			newPassword: "newtestpassword",
			wantCode:    codes.OK,
		},
		{
			testname:    "empty name",
			name:        "",
			password:    "testpassword",
			newPassword: "newtestpassword",
			wantCode:    codes.InvalidArgument,
		},
		{
			testname:    "empty password",
			name:        "testuser",
			password:    "",
			newPassword: "newtestpassword",
			wantCode:    codes.InvalidArgument,
		},
		{
			testname:    "empty new password",
			name:        "testuser",
			password:    "testpassword",
			newPassword: "",
			wantCode:    codes.InvalidArgument,
		},
		{
			testname:    "invalid user",
			name:        "testuser",
			password:    "testpassword",
			newPassword: "newtestpassword",
			wantCode:    codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			server := &VaultServer{
				userHandler: &mocks.UserHandler{},
			}
			mockToken := "mocktoken"
			if tt.wantCode == codes.OK {
				mockUserHandler := &mocks.UserHandler{}
				mockUserHandler.On("ChangePassword", mockCtx, models.User{
					Login:    tt.name,
					Password: tt.password,
				}, tt.newPassword).Return(mockToken, nil)
				server.userHandler = mockUserHandler
			} else {
				mockUserHandler := &mocks.UserHandler{}
				mockUserHandler.On("ChangePassword", mockCtx, models.User{
					Login:    tt.name,
					Password: tt.password,
				}, tt.newPassword).Return("", errors.New("error"))
				server.userHandler = mockUserHandler
			}
			req := &pb.ChangePasswordRequest{
				User: &pb.User{
					Name:     tt.name,
					Password: tt.password,
				},
				NewPassword: tt.newPassword,
			}
			_, err := server.ChangePassword(mockCtx, req)
			assert.Equal(t, tt.wantCode, status.Code(err))
		})

	}
}

func TestVaultServer_SaveSecret(t *testing.T) {
	var mockCtx context.Context
	tests := []struct {
		testname string
		user     string
		mdExists bool
		wantCode codes.Code
	}{
		{
			testname: "valid",
			user:     "testuser",
			mdExists: true,
			wantCode: codes.OK,
		},
		{
			testname: "empty user",
			user:     "",
			mdExists: true,
			wantCode: codes.InvalidArgument,
		},
		{
			testname: "user not found",
			user:     "testuser",
			mdExists: true,
			wantCode: codes.Internal,
		},
		{
			testname: "no metadata in context",
			user:     "testuser",
			mdExists: false,
			wantCode: codes.InvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			mockCtx = context.Background()
			mockReq := &pb.SaveSecretRequest{
				Data: &pb.SecretData{
					Meta:  "testmeta",
					Type:  "testtype",
					Value: []byte("testvalue"),
				},
			}
			mockUser := models.User{
				UUID:  uuid.New().String(),
				Login: tt.user}

			mockDataUUID := uuid.New().String()
			mockCreated := time.Now().Unix()
			server := &VaultServer{
				userHandler: &mocks.UserHandler{},
				dataHandler: &mocks.DataHandler{},
			}
			if tt.mdExists {
				md := make(map[string]string)
				md["user"] = tt.user
				mockCtx = metadata.NewIncomingContext(mockCtx, metadata.New(md))
			}
			mockUserHandler := &mocks.UserHandler{}
			mockDataHandler := &mocks.DataHandler{}
			if tt.testname != "user not found" {
				mockUserHandler.On("GetUser", mockCtx, tt.user).Return(mockUser, nil)
				mockDataHandler.On("CreateData", mockCtx, mockUser, models.VaultData{
					Meta:     mockReq.Data.Meta,
					DataType: mockReq.Data.Type,
					Data:     mockReq.Data.Value,
				}).Return(mockDataUUID, mockCreated, nil)
			} else {
				mockUserHandler.On("GetUser", mockCtx, tt.user).Return(models.User{}, errors.New("error"))
			}

			server.userHandler = mockUserHandler
			server.dataHandler = mockDataHandler

			_, err := server.SaveSecret(mockCtx, mockReq)
			assert.Equal(t, tt.wantCode, status.Code(err))

		})
	}
}

func TestVaultServer_ChangeSecret(t *testing.T) {
	var mockCtx context.Context
	tests := []struct {
		testname string
		user     string
		mdExists bool
		wantCode codes.Code
	}{
		{
			testname: "valid",
			user:     "testuser",
			mdExists: true,
			wantCode: codes.OK,
		},
		{
			testname: "empty user",
			user:     "",
			mdExists: true,
			wantCode: codes.InvalidArgument,
		},
		{
			testname: "user not found",
			user:     "testuser",
			mdExists: true,
			wantCode: codes.Internal,
		},
		{
			testname: "no metadata in context",
			user:     "testuser",
			mdExists: false,
			wantCode: codes.InvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			mockCtx = context.Background()
			mockReq := &pb.ChangeSecretRequest{
				Data: &pb.SecretData{
					Meta:  "testmeta",
					Type:  "testtype",
					Value: []byte("testvalue"),
				},
			}
			mockUser := models.User{
				UUID:  uuid.New().String(),
				Login: tt.user}

			mockCreated := time.Now().Unix()
			server := &VaultServer{
				userHandler: &mocks.UserHandler{},
				dataHandler: &mocks.DataHandler{},
			}
			if tt.mdExists {
				md := make(map[string]string)
				md["user"] = tt.user
				mockCtx = metadata.NewIncomingContext(mockCtx, metadata.New(md))
			}
			mockUserHandler := &mocks.UserHandler{}
			mockDataHandler := &mocks.DataHandler{}
			if tt.testname != "user not found" {
				mockUserHandler.On("GetUser", mockCtx, tt.user).Return(mockUser, nil)
				mockDataHandler.On("ChangeData", mockCtx, mockUser, models.VaultData{
					Meta:     mockReq.Data.Meta,
					DataType: mockReq.Data.Type,
					Data:     mockReq.Data.Value,
				}).Return(mockCreated, nil)
			} else {
				mockUserHandler.On("GetUser", mockCtx, tt.user).Return(models.User{}, errors.New("error"))
			}

			server.userHandler = mockUserHandler
			server.dataHandler = mockDataHandler

			_, err := server.ChangeSecret(mockCtx, mockReq)
			assert.Equal(t, tt.wantCode, status.Code(err))

		})
	}
}

func TestVaultServer_DeleteSecret(t *testing.T) {
	var mockCtx context.Context
	tests := []struct {
		testname string
		user     string
		mdExists bool
		wantCode codes.Code
	}{
		{
			testname: "valid",
			user:     "testuser",
			mdExists: true,
			wantCode: codes.OK,
		},
		{
			testname: "empty user",
			user:     "",
			mdExists: true,
			wantCode: codes.InvalidArgument,
		},
		{
			testname: "user not found",
			user:     "testuser",
			mdExists: true,
			wantCode: codes.Internal,
		},
		{
			testname: "no metadata in context",
			user:     "testuser",
			mdExists: false,
			wantCode: codes.InvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			mockCtx = context.Background()
			mockReq := &pb.DeleteSecretRequest{
				Uuid: uuid.New().String(),
			}
			mockUser := models.User{
				UUID:  uuid.New().String(),
				Login: tt.user}

			mockCreated := time.Now().Unix()
			server := &VaultServer{
				userHandler: &mocks.UserHandler{},
				dataHandler: &mocks.DataHandler{},
			}
			if tt.mdExists {
				md := make(map[string]string)
				md["user"] = tt.user
				mockCtx = metadata.NewIncomingContext(mockCtx, metadata.New(md))
			}
			mockUserHandler := &mocks.UserHandler{}
			mockDataHandler := &mocks.DataHandler{}
			if tt.testname != "user not found" {
				mockUserHandler.On("GetUser", mockCtx, tt.user).Return(mockUser, nil)
				mockDataHandler.On("DeleteData", mockCtx, mockUser, models.VaultData{
					DataUUID: mockReq.Uuid,
				}).Return(mockCreated, nil)
			} else {
				mockUserHandler.On("GetUser", mockCtx, tt.user).Return(models.User{}, errors.New("error"))
			}

			server.userHandler = mockUserHandler
			server.dataHandler = mockDataHandler

			_, err := server.DeleteSecret(mockCtx, mockReq)
			assert.Equal(t, tt.wantCode, status.Code(err))
		})
	}
}

func TestVaultServer_ListSecrets(t *testing.T) {
	var mockCtx context.Context

	tests := []struct {
		testname string
		user     string
		mdExists bool
		wantCode codes.Code
	}{
		{
			testname: "valid",
			user:     "testuser",
			mdExists: true,
			wantCode: codes.OK,
		},
		{
			testname: "empty user",
			user:     "",
			mdExists: true,
			wantCode: codes.InvalidArgument,
		},
		{
			testname: "user not found",
			user:     "testuser",
			mdExists: true,
			wantCode: codes.Internal,
		},
		{
			testname: "no metadata in context",
			user:     "testuser",
			mdExists: false,
			wantCode: codes.InvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			mockCtx = context.Background()
			mockReq := &pb.ListSecretsRequest{}

			mockUser := models.User{
				UUID:  uuid.New().String(),
				Login: tt.user}

			server := &VaultServer{
				userHandler: &mocks.UserHandler{},
				dataHandler: &mocks.DataHandler{},
			}
			if tt.mdExists {
				md := make(map[string]string)
				md["user"] = tt.user
				mockCtx = metadata.NewIncomingContext(mockCtx, metadata.New(md))
			}
			mockUserHandler := &mocks.UserHandler{}
			mockDataHandler := &mocks.DataHandler{}
			if tt.testname != "user not found" {
				mockUserHandler.On("GetUser", mockCtx, tt.user).Return(mockUser, nil)
				mockDataHandler.On("GetAllData", mockCtx, mockUser).Return([]models.VaultData{}, nil)
			} else {
				mockUserHandler.On("GetUser", mockCtx, tt.user).Return(models.User{}, errors.New("error"))
			}

			server.userHandler = mockUserHandler
			server.dataHandler = mockDataHandler

			_, err := server.ListSecrets(mockCtx, mockReq)
			assert.Equal(t, tt.wantCode, status.Code(err))
		})
	}
}
