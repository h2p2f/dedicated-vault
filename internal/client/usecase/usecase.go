package usecase

import (
	"context"
	"fmt"
	"github.com/h2p2f/dedicated-vault/internal/client/config"
	"github.com/h2p2f/dedicated-vault/internal/client/models"
	pb "github.com/h2p2f/dedicated-vault/proto"
)

type Storager interface {
	CreateUser(userName string) error
	GetUserID(userName string) (int64, error)
	CreateData(user string, data models.StoredData) error
	UpdateData(user string, data models.StoredData) error
	DeleteData(user string, data models.StoredData) error
	GetDataByUUID(user string, uuid string) (*models.StoredData, error)
	DeleteAllData(user string) error
}

type Transporter interface {
	Register(ctx context.Context, user *pb.User, clientID string) (string, error)
	Login(ctx context.Context, user *pb.User, clientID string) (string, error)
	ChangePassword(ctx context.Context, user *pb.User, newPassword string) (string, error)
	SaveSecret(ctx context.Context, data *pb.SecretData) error
	ChangeSecret(ctx context.Context, data *pb.SecretData) error
	DeleteSecret(ctx context.Context, uuid string) error
	ListSecrets(ctx context.Context) ([]*pb.SecretData, error)
}
type ClientUseCase struct {
	Storage     Storager
	Transporter Transporter
	Config      *config.ClientConfig
}

func NewClientUseCase(config *config.ClientConfig, storage Storager, transporter Transporter) *ClientUseCase {
	return &ClientUseCase{
		Config:      config,
		Storage:     storage,
		Transporter: transporter,
	}
}

func (c *ClientUseCase) CreateUser(userName, password, passphrase string) error {
	err := c.Storage.CreateUser(userName)
	if err != nil {
		return err
	}
	c.Config.Passphrase = passphrase
	user := &pb.User{
		Name:     userName,
		Password: password,
	}
	fmt.Println("go to transporter")
	token, err := c.Transporter.Register(context.Background(), user, "testClientID")
	if err != nil {
		return err
	}
	c.Config.Token = token
	return nil
}

func (c *ClientUseCase) LoginUser(userName, password, passphrase string) error {
	id, err := c.Storage.GetUserID(userName)
	if err != nil {
		return err
	}
	if id == 0 {
		return fmt.Errorf("user not found")
	}
	c.Config.Passphrase = passphrase
	user := &pb.User{
		Name:     userName,
		Password: password,
	}
	token, err := c.Transporter.Login(context.Background(), user, "")
	if err != nil {
		return err
	}
	c.Config.User = userName
	c.Config.Token = token
	return nil
}

func (c *ClientUseCase) ChangePassword(userName, password, newPassword string) error {
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
	token, err := c.Transporter.ChangePassword(context.Background(), user, newPassword)
	if err != nil {
		return err
	}
	c.Config.User = userName
	c.Config.Token = token
	return nil
}

func (c *ClientUseCase) SaveData(data models.Data) error {
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
	err = c.Transporter.SaveSecret(context.Background(), secretData)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientUseCase) ChangeData(data models.Data) error {
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
	err = c.Transporter.ChangeSecret(context.Background(), secretData)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientUseCase) DeleteData(data models.Data) error {
	if c.Config.Token == "" {
		return fmt.Errorf("user not logged in")
	}
	var storedData models.StoredData
	storedData.UUID = data.UUID
	err := c.Storage.DeleteData(c.Config.User, storedData)
	if err != nil {
		return err
	}
	err = c.Transporter.DeleteSecret(context.Background(), data.UUID)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientUseCase) GetData(uuid string) (*models.Data, error) {
	if c.Config.Token == "" {
		return nil, fmt.Errorf("user not logged in")
	}
	storedData, err := c.Storage.GetDataByUUID(c.Config.User, uuid)
	if err != nil {
		return nil, err
	}
	data, err := storedData.DecryptData([]byte(c.Config.Passphrase))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *ClientUseCase) FullSync() error {
	if c.Config.Token == "" {
		return fmt.Errorf("user not logged in")
	}
	err := c.Storage.DeleteAllData(c.Config.User)
	if err != nil {
		return err
	}
	secrets, err := c.Transporter.ListSecrets(context.Background())
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
	return nil
}
