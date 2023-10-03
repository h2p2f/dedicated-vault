// Package: grpcclient
// in this file we have grpc client
package grpcclient

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/h2p2f/dedicated-vault/internal/client/config"
	"github.com/h2p2f/dedicated-vault/internal/client/grpcclient/middlewares"
	pb "github.com/h2p2f/dedicated-vault/proto"
	//"google.golang.org/grpc/credentials"
)

// Client is a struct for grpc client
type Client struct {
	pb.DedicatedVaultClient
	config *config.ClientConfig
	logger *zap.Logger
}

// NewClient creates a new Client
func NewClient(config *config.ClientConfig, logger *zap.Logger) *Client {
	return &Client{
		config: config,
		logger: logger,
	}
}

// Connect connects to the server
func (c *Client) Connect() (*grpc.ClientConn, error) {

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(c.config.TLSConfig),
		grpc.WithUnaryInterceptor(middlewares.JWTInjectorUnaryClientInterceptor(c.config.Token)),
	}
	conn, err := grpc.Dial(c.config.StorageAddress, opts...)
	if err != nil {
		return nil, err
	}
	c.DedicatedVaultClient = pb.NewDedicatedVaultClient(conn)
	return conn, nil
}

// Register registers a new user
func (c *Client) Register(ctx context.Context, user *pb.User) (string, error) {
	conn, err := c.Connect()
	if err != nil {
		return "", err
	}

	resp, err := c.DedicatedVaultClient.Register(ctx, &pb.RegisterRequest{
		User: user,
	})
	if err != nil {
		return "", err
	}
	c.config.Token = resp.Token
	c.config.LastServerUpdated = resp.LastServerUpdated
	err = conn.Close()
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

// Login logs in a user
func (c *Client) Login(ctx context.Context, user *pb.User) (string, error) {
	conn, err := c.Connect()
	if err != nil {
		return "", err
	}
	resp, err := c.DedicatedVaultClient.Login(ctx, &pb.LoginRequest{
		User: user,
	})
	if err != nil {
		return "", err
	}
	c.config.Token = resp.Token
	c.config.LastServerUpdated = resp.LastServerUpdated
	err = conn.Close()
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

// ChangePassword changes password for a user
func (c *Client) ChangePassword(ctx context.Context, user *pb.User, newPassword string) (string, error) {
	resp, err := c.DedicatedVaultClient.ChangePassword(ctx, &pb.ChangePasswordRequest{
		User:        user,
		NewPassword: newPassword,
	})
	if err != nil {
		return "", err
	}
	c.config.Token = resp.Token
	return resp.Token, nil
}

// SaveSecret gets a secret
func (c *Client) SaveSecret(ctx context.Context, data *pb.SecretData) error {
	conn, err := c.Connect()
	if err != nil {
		return err
	}
	resp, err := c.DedicatedVaultClient.SaveSecret(ctx, &pb.SaveSecretRequest{
		Data: data,
	})
	if err != nil {
		return err
	}
	c.config.LastServerUpdated = resp.LastServerUpdated
	err = conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// ChangeSecret gets a secret
func (c *Client) ChangeSecret(ctx context.Context, data *pb.SecretData) error {
	conn, err := c.Connect()
	if err != nil {
		return err
	}
	resp, err := c.DedicatedVaultClient.ChangeSecret(ctx, &pb.ChangeSecretRequest{
		Data: data,
	})
	if err != nil {
		return err
	}
	c.config.LastServerUpdated = resp.LastServerUpdated
	err = conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// DeleteSecret gets a secret
func (c *Client) DeleteSecret(ctx context.Context, uuid string) error {
	conn, err := c.Connect()
	if err != nil {
		return err
	}
	resp, err := c.DedicatedVaultClient.DeleteSecret(ctx, &pb.DeleteSecretRequest{
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	c.config.LastServerUpdated = resp.LastServerUpdated
	err = conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// ListSecrets gets a secret
func (c *Client) ListSecrets(ctx context.Context) ([]*pb.SecretData, error) {
	conn, err := c.Connect()
	if err != nil {
		return nil, err
	}
	resp, err := c.DedicatedVaultClient.ListSecrets(ctx, &pb.ListSecretsRequest{})
	if err != nil {
		return nil, err
	}
	c.config.LastServerUpdated = resp.LastServerUpdated
	err = conn.Close()
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
