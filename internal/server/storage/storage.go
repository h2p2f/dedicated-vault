package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/h2p2f/dedicated-vault/internal/server/config"
	"github.com/h2p2f/dedicated-vault/internal/server/jwtprocessing"
	"github.com/h2p2f/dedicated-vault/internal/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Storage struct {
	users  *mongo.Collection
	data   *mongo.Collection
	config *config.ServerConfig
	logger *zap.Logger
}

func NewStorage(ctx context.Context, config *config.ServerConfig, logger *zap.Logger) *Storage {
	var storage Storage

	credential := options.Credential{
		Username: config.DBUser,
		Password: config.DBPassword,
	}

	opts := options.Client().ApplyURI(config.StorageAddress).SetAuth(credential)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Fatal("error while connecting to database", zap.Error(err))
	}
	db := client.Database("vault")
	storage.users = db.Collection("users")
	storage.data = db.Collection("data")
	storage.logger = logger
	storage.config = config

	return &storage
}

func (s *Storage) Close(ctx context.Context) error {
	if err := s.users.Database().Client().Disconnect(ctx); err != nil {
		s.logger.Error("error while disconnecting from database", zap.Error(err))
		return err
	}
	return nil
}

func (s *Storage) UpdateLastServerUpdated(ctx context.Context, user models.User) error {
	_, err := s.users.UpdateOne(ctx,
		bson.D{{"login", user.Login}},
		bson.D{{"$set", bson.D{{"lastServerUpdated", user.LastServerUpdated}}}})
	if err != nil {
		s.logger.Error("error while updating lastServerUpdated", zap.Error(err))
		return err
	}
	return nil
}

func (s *Storage) Register(ctx context.Context, user models.User) (string, int64, error) {
	var checkUser models.User
	var token string
	var lastServerUpdated int64
	err := s.users.FindOne(ctx, bson.D{{"login", user.Login}}).Decode(&checkUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			uuidUser := uuid.New()
			lastServerUpdated := time.Now().Unix()
			encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				s.logger.Error("error while encrypting password", zap.Error(err))
				return token, lastServerUpdated, err
			}
			docUser := bson.D{
				{"UUID", uuidUser.String()},
				{"login", user.Login},
				{"password", string(encryptedPassword)},
				{"lastServerUpdated", lastServerUpdated},
			}

			_, err = s.users.InsertOne(ctx, docUser)
			if err != nil {
				s.logger.Error("error while inserting user", zap.Error(err))
				return token, lastServerUpdated, err
			}
			token, err = jwtprocessing.GenerateToken(uuidUser.String(), s.config.JWTKey)
			fmt.Println(token)
			if err != nil {
				s.logger.Error("error while generating token", zap.Error(err))
				return token, lastServerUpdated, err
			}
			return token, lastServerUpdated, nil
		} else {
			s.logger.Error("error while finding user", zap.Error(err))
			return token, lastServerUpdated, err
		}
	}
	return token, lastServerUpdated, errors.New("user already exists")
}

func (s *Storage) Login(ctx context.Context, user models.User) (string, int64, error) {
	var checkUser models.User
	var token string
	var lastServerUpdated int64
	err := s.users.FindOne(ctx, bson.D{{"login", user.Login}}).Decode(&checkUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.Error("error while finding user", zap.Error(err))
			return token, lastServerUpdated, errors.New("user not found")
		} else {
			s.logger.Error("error while finding user", zap.Error(err))
			return token, lastServerUpdated, err
		}
	}
	lastServerUpdated = checkUser.LastServerUpdated
	err = bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(user.Password))
	if err != nil {
		s.logger.Error("error while comparing passwords", zap.Error(err))
		return token, lastServerUpdated, err
	}
	token, err = jwtprocessing.GenerateToken(checkUser.UUID, s.config.JWTKey)
	if err != nil {
		s.logger.Error("error while generating token", zap.Error(err))
		return token, lastServerUpdated, err
	}
	return token, lastServerUpdated, nil
}

func (s *Storage) GetUser(ctx context.Context, user string) (models.User, error) {
	var checkUser models.User
	err := s.users.FindOne(ctx, bson.D{{"UUID", user}}).Decode(&checkUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.Error("error while finding user", zap.Error(err))
			return models.User{}, errors.New("user not found")
		} else {
			s.logger.Error("error while finding user", zap.Error(err))
			return models.User{}, err
		}
	}
	return checkUser, nil
}

func (s *Storage) ChangePassword(ctx context.Context, user models.User, newPassword string) (string, error) {
	var checkUser models.User
	err := s.users.FindOne(ctx, bson.D{{"login", user.Login}}).Decode(&checkUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.Error("error while finding user", zap.Error(err))
			return "", errors.New("user not found")
		} else {
			s.logger.Error("error while finding user", zap.Error(err))
			return "", err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(user.Password))
	if err != nil {
		s.logger.Error("error while comparing passwords", zap.Error(err))
		return "", err
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("error while encrypting password", zap.Error(err))
		return "", err
	}
	_, err = s.users.UpdateOne(ctx,
		bson.D{{"login", user.Login}},
		bson.D{{"$set", bson.D{{"password", string(encryptedPassword)}}}})
	if err != nil {
		s.logger.Error("error while updating password", zap.Error(err))
		return "", err
	}
	token, err := jwtprocessing.GenerateToken(checkUser.UUID, s.config.JWTKey)
	if err != nil {
		s.logger.Error("error while generating token", zap.Error(err))
		return "", err
	}
	return token, nil
}

func (s *Storage) DeleteUser(ctx context.Context, user models.User) error {
	_, err := s.users.DeleteOne(ctx, bson.D{{"login", user.Login}})
	if err != nil {
		s.logger.Error("error while deleting user", zap.Error(err))
		return err
	}
	return nil
}

func (s *Storage) CreateData(ctx context.Context, user models.User, data models.VaultData) (string, int64, error) {
	uuidData := uuid.New()
	data.DataUUID = uuidData.String()
	data.UserUUID = user.UUID
	data.Created = time.Now().Unix()
	doc, err := bson.Marshal(data)
	if err != nil {
		s.logger.Error("error while marshaling data", zap.Error(err))
	}
	_, err = s.data.InsertOne(ctx, doc)
	if err != nil {
		s.logger.Error("error while inserting data", zap.Error(err))
		return "", 0, err
	}
	user.LastServerUpdated = data.Created
	err = s.UpdateLastServerUpdated(ctx, user)
	if err != nil {
		s.logger.Error("error while updating lastServerUpdated", zap.Error(err))
		return "", 0, err
	}
	return data.DataUUID, data.Created, nil
}

func (s *Storage) ChangeData(ctx context.Context, user models.User, data models.VaultData) (int64, error) {
	data.Updated = time.Now().Unix()
	_, err := s.data.ReplaceOne(ctx,
		bson.D{{"userUUID", user.UUID}, {"dataUUID", data.DataUUID}},
		data)
	if err != nil {
		s.logger.Error("error while replacing data", zap.Error(err))
		return 0, err
	}
	user.LastServerUpdated = data.Updated
	err = s.UpdateLastServerUpdated(ctx, user)
	if err != nil {
		s.logger.Error("error while updating lastServerUpdated", zap.Error(err))
		return 0, err
	}
	return data.Updated, nil
}

func (s *Storage) GetAllData(ctx context.Context, user models.User) ([]models.VaultData, error) {
	var data []models.VaultData

	filter := bson.D{
		{"userUUID", user.UUID},
	}

	cur, err := s.data.Find(ctx, filter)

	if err != nil {
		s.logger.Error("error while finding data", zap.Error(err))
		return nil, err
	}

	for cur.Next(ctx) {
		var elem models.VaultData
		err := cur.Decode(&elem)
		if err != nil {
			s.logger.Error("error while decoding data", zap.Error(err))
			return nil, err
		}
		data = append(data, elem)
	}

	err = cur.Close(ctx)
	if err != nil {
		s.logger.Error("error while closing cursor", zap.Error(err))
		return nil, err
	}
	return data, nil
}

func (s *Storage) DeleteData(ctx context.Context, user models.User, data models.VaultData) (int64, error) {

	_, err := s.data.DeleteOne(ctx,
		bson.D{{"user", user.UUID}, {"data.dataUUID", data.DataUUID}})
	if err != nil {
		s.logger.Error("error while deleting data", zap.Error(err))
		return 0, err
	}
	user.LastServerUpdated = time.Now().Unix()
	return user.LastServerUpdated, nil
}
