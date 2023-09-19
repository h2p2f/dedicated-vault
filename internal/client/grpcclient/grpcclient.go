package grpcclient

import (
	"context"
	"github.com/h2p2f/dedicated-vault/internal/client/config"
	pb "github.com/h2p2f/dedicated-vault/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	//"google.golang.org/grpc/credentials"
)

//type Storager interface {
//	CreateData(ctx context.Context, data models.Data) error
//	GetData(ctx context.Context) ([]models.Data, error)
//	UpdateData(ctx context.Context, data models.Data) error
//	DeleteData(ctx context.Context, data models.Data) error
//}

type Client struct {
	pb.DedicatedVaultClient
	config *config.ClientConfig
	logger *zap.Logger
	Token  string
}

func NewClient(config *config.ClientConfig, logger *zap.Logger) *Client {
	return &Client{
		//Storage:              storage,
		config: config,
		//DedicatedVaultClient: vaultClient,
		logger: logger,
		Token:  "",
	}
}

func (c *Client) Connect() (*grpc.ClientConn, error) {
	//creds, _ := credentials.NewClientTLSFromFile("./crypto/public.crt", "localhost")
	////creds := credentials.NewClientTLSFromCert(c.config.Cert, "localhost")
	//opts := []grpc.DialOption{
	//	grpc.WithTransportCredentials(creds),
	//	//grpc.WithBlock(),
	//}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(c.config.StorageAddress, opts...)
	if err != nil {
		return nil, err
	}
	c.DedicatedVaultClient = pb.NewDedicatedVaultClient(conn)
	return conn, nil
}

func (c *Client) Register(ctx context.Context, user *pb.User, clientID string) (string, error) {
	conn, err := c.Connect()
	if err != nil {
		return "", err
	}

	resp, err := c.DedicatedVaultClient.Register(ctx, &pb.RegisterRequest{
		User:     user,
		ClientId: clientID,
	})
	if err != nil {
		return "", err
	}
	c.Token = resp.Token
	err = conn.Close()
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (c *Client) Login(ctx context.Context, user *pb.User, clientID string) (string, error) {
	conn, err := c.Connect()
	if err != nil {
		return "", err
	}
	resp, err := c.DedicatedVaultClient.Login(ctx, &pb.LoginRequest{
		User:     user,
		ClientId: clientID,
	})
	if err != nil {
		return "", err
	}
	c.Token = resp.Token
	err = conn.Close()
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (c *Client) ChangePassword(ctx context.Context, user *pb.User, newPassword string) (string, error) {
	resp, err := c.DedicatedVaultClient.ChangePassword(ctx, &pb.ChangePasswordRequest{
		User:        user,
		NewPassword: newPassword,
	})
	if err != nil {
		return "", err
	}
	c.Token = resp.Token
	return resp.Token, nil
}

func (c *Client) SaveSecret(ctx context.Context, data *pb.SecretData) error {
	_, err := c.DedicatedVaultClient.SaveSecret(ctx, &pb.SaveSecretRequest{
		Data: data,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetSecret(ctx context.Context, uuid string) (*pb.SecretData, error) {
	resp, err := c.DedicatedVaultClient.GetSecret(ctx, &pb.GetSecretRequest{
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) ChangeSecret(ctx context.Context, data *pb.SecretData) error {
	_, err := c.DedicatedVaultClient.ChangeSecret(ctx, &pb.ChangeSecretRequest{
		Data: data,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteSecret(ctx context.Context, uuid string) error {
	_, err := c.DedicatedVaultClient.DeleteSecret(ctx, &pb.DeleteSecretRequest{
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ListSecrets(ctx context.Context) ([]*pb.SecretData, error) {
	resp, err := c.DedicatedVaultClient.ListSecrets(ctx, &pb.ListSecretsRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
