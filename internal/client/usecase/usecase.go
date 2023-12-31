// Package: usecase
// in this file we have main logic for client usecase
package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/h2p2f/dedicated-vault/internal/client/clienterrors"
	"github.com/h2p2f/dedicated-vault/internal/client/config"
	"github.com/h2p2f/dedicated-vault/internal/client/models"
	pb "github.com/h2p2f/dedicated-vault/proto"
)

// Storager interface for working with storage
//
//go:generate mockery --name Storager --output ./mocks --filename mocks_storager.go
type Storager interface {
	CreateUser(userName string) error
	GetUserID(userName string) (int64, error)
	CreateData(user string, data models.StoredData) error
	UpdateData(user string, data models.StoredData) error
	DeleteData(user string, data models.StoredData) error
	GetDataByUUID(user string, uuid string) (*models.StoredData, error)
	GetData(user string) ([]models.StoredData, error)
	FindByMeta(user string, meta string) ([]models.StoredData, error)
	DeleteAllData(user string) error
	UpdateLastServerUpdated(username string, updateTime int64) error
	GetLastServerUpdated(username string) (int64, error)
}

// Transporter interface for working with transporter
//
//go:generate mockery --name Transporter --output ./mocks --filename mocks_transporter.go
type Transporter interface {
	Register(ctx context.Context, user *pb.User) (string, error)
	Login(ctx context.Context, user *pb.User) (string, error)
	ChangePassword(ctx context.Context, user *pb.User, newPassword string) (string, error)
	SaveSecret(ctx context.Context, data *pb.SecretData) error
	ChangeSecret(ctx context.Context, data *pb.SecretData) error
	DeleteSecret(ctx context.Context, uuid string) error
	ListSecrets(ctx context.Context) ([]*pb.SecretData, error)
}

// ClientUseCase is a struct for client usecase
type ClientUseCase struct {
	Storage     Storager
	Transporter Transporter
	Config      *config.ClientConfig
}

// NewClientUseCase creates a new ClientUseCase
func NewClientUseCase(config *config.ClientConfig, storage Storager, transporter Transporter) *ClientUseCase {
	return &ClientUseCase{
		Config:      config,
		Storage:     storage,
		Transporter: transporter,
	}
}

// CreateUser creates a new user
func (c *ClientUseCase) CreateUser(ctx context.Context, userName, password, passphrase string) error {
	err := c.Storage.CreateUser(userName)
	if err != nil {
		return err
	}
	user := &pb.User{
		Name:     userName,
		Password: password,
	}
	token, err := c.Transporter.Register(ctx, user)
	if err != nil {
		return err
	}
	c.Config.Token = token
	c.Config.User = userName
	c.Config.Passphrase = passphrase
	err = c.Storage.UpdateLastServerUpdated(userName, c.Config.LastServerUpdated)
	if err != nil {
		return err
	}
	return nil
}

// LoginUser login user
func (c *ClientUseCase) LoginUser(ctx context.Context, userName, password, passphrase string) error {
	user := &pb.User{
		Name:     userName,
		Password: password,
	}
	token, err := c.Transporter.Login(ctx, user)
	if err != nil {
		return err
	}

	_, err = c.Storage.GetUserID(userName)

	if errors.Is(err, clienterrors.UserNotFound) {
		err = c.Storage.CreateUser(userName)
		if err != nil {
			return err
		}
	}

	c.Config.Passphrase = passphrase
	c.Config.User = userName
	c.Config.Token = token

	dbLastServerUpdated, err := c.Storage.GetLastServerUpdated(userName)
	if err != nil {
		return err
	}
	if dbLastServerUpdated < c.Config.LastServerUpdated {
		err = c.FullSync(ctx)
		if err != nil {
			return err
		}
	}

	err = c.Storage.UpdateLastServerUpdated(userName, c.Config.LastServerUpdated)
	if err != nil {
		return err
	}
	return nil
}

// ChangePassword changes user password
func (c *ClientUseCase) ChangePassword(ctx context.Context, userName, password, newPassword string) error {
	id, err := c.Storage.GetUserID(userName)
	if err != nil {
		return err
	}
	if id == 0 {
		return fmt.Errorf("user not found")
	}

	user := &pb.User{
		Name:     userName,
		Password: password,
	}
	token, err := c.Transporter.ChangePassword(ctx, user, newPassword)
	if err != nil {
		return err
	}
	c.Config.User = userName
	c.Config.Token = token
	return nil
}

// SaveData saves data
func (c *ClientUseCase) SaveData(ctx context.Context, data models.Data) error {
	if c.Config.Token == "" {
		return fmt.Errorf("user not logged in")
	}

	storedData, err := data.EncryptData([]byte(c.Config.Passphrase))
	if err != nil {
		return err
	}
	err = c.Storage.CreateData(c.Config.User, *storedData)
	if err != nil {
		return err
	}
	secretData := &pb.SecretData{
		Uuid:  storedData.UUID,
		Meta:  storedData.Meta,
		Type:  storedData.DataType,
		Value: storedData.EncryptedData,
	}
	err = c.Transporter.SaveSecret(ctx, secretData)
	if err != nil {
		return err
	}
	err = c.Storage.UpdateLastServerUpdated(c.Config.User, c.Config.LastServerUpdated)
	if err != nil {
		return err
	}
	return nil
}

// ChangeData changes data
func (c *ClientUseCase) ChangeData(ctx context.Context, data models.Data) error {
	if c.Config.Token == "" {
		return fmt.Errorf("user not logged in")
	}
	storedData, err := data.EncryptData([]byte(c.Config.Passphrase))
	if err != nil {
		return err
	}
	err = c.Storage.UpdateData(c.Config.User, *storedData)
	if err != nil {
		return err
	}
	secretData := &pb.SecretData{
		Uuid:  storedData.UUID,
		Meta:  storedData.Meta,
		Type:  storedData.DataType,
		Value: storedData.EncryptedData,
	}
	err = c.Transporter.ChangeSecret(ctx, secretData)
	if err != nil {
		return err
	}
	err = c.Storage.UpdateLastServerUpdated(c.Config.User, c.Config.LastServerUpdated)
	if err != nil {
		return err
	}
	return nil
}

// DeleteData deletes data
func (c *ClientUseCase) DeleteData(ctx context.Context, data models.Data) error {
	if c.Config.Token == "" {
		return fmt.Errorf("user not logged in")
	}
	var storedData models.StoredData
	storedData.UUID = data.UUID
	err := c.Storage.DeleteData(c.Config.User, storedData)
	if err != nil {
		return err
	}
	err = c.Transporter.DeleteSecret(ctx, data.UUID)
	if err != nil {
		return err
	}
	err = c.Storage.UpdateLastServerUpdated(c.Config.User, c.Config.LastServerUpdated)
	if err != nil {
		return err
	}
	return nil
}

// GetDataByType gets data by type
func (c *ClientUseCase) GetDataByType(dataType string) ([]models.Data, error) {
	if c.Config.Token == "" {
		return nil, fmt.Errorf("user not logged in")
	}
	storedData, err := c.Storage.GetData(c.Config.User)
	if err != nil {
		return nil, err
	}
	var data []models.Data
	for _, d := range storedData {
		if d.DataType == dataType {
			decryptData, err := d.DecryptData([]byte(c.Config.Passphrase))
			if err != nil {
				return nil, err
			}
			data = append(data, *decryptData)
		}
	}
	return data, nil
}

// FullSync does full sync with remote server
func (c *ClientUseCase) FullSync(ctx context.Context) error {
	if c.Config.Token == "" {
		return fmt.Errorf("user not logged in")
	}
	err := c.Storage.DeleteAllData(c.Config.User)
	if err != nil {
		fmt.Println(err)
		return err
	}
	secrets, err := c.Transporter.ListSecrets(ctx)
	if err != nil {
		return err
	}
	for _, secret := range secrets {
		storedData := models.StoredData{
			UUID:          secret.Uuid,
			Meta:          secret.Meta,
			DataType:      secret.Type,
			EncryptedData: secret.Value,
		}
		err = c.Storage.CreateData(c.Config.User, storedData)
		if err != nil {
			return err
		}
	}
	err = c.Storage.UpdateLastServerUpdated(c.Config.User, c.Config.LastServerUpdated)
	if err != nil {
		return err
	}
	return nil
}
