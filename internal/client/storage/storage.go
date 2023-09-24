package storage

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/h2p2f/dedicated-vault/internal/client/config"
	"github.com/h2p2f/dedicated-vault/internal/client/models"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

const createTable = `

CREATE TABLE IF NOT EXISTS data (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	user_id INTEGER NOT NULL,
    	uuid TEXT NOT NULL,
    	meta TEXT NOT NULL,
    	type TEXT NOT NULL,
    	data BLOB NOT NULL,
    	FOREIGN KEY (user_id) REFERENCES users (id)
	);
CREATE TABLE IF NOT EXISTS users (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	username TEXT NOT NULL UNIQUE
    	);
`

type ClientStorage struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewClientStorage(logger *zap.Logger, config *config.ClientConfig) *ClientStorage {
	db, err := sql.Open("sqlite3", config.DBPath)
	if err != nil {
		logger.Fatal("failed to open database", zap.Error(err))
	}
	_, err = db.Exec(createTable)
	if err != nil {
		logger.Fatal("failed to create table", zap.Error(err))
	}

	return &ClientStorage{
		db:     db,
		logger: logger,
	}
}

func (s *ClientStorage) Close() error {
	return s.db.Close()
}

func (s *ClientStorage) CreateUser(userName string) error {
	_, err := s.db.Exec("INSERT INTO users (username) VALUES (?)", userName)
	if err != nil {
		s.logger.Error("failed to insert user", zap.Error(err))
		return err
	}
	return nil
}

func (s *ClientStorage) GetUserID(userName string) (int64, error) {
	row := s.db.QueryRow("SELECT id FROM users WHERE username = ?", userName)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		//s.logger.Error("failed to scan user", zap.Error(err))
		return 0, err
	}

	return id, nil

}

func (s *ClientStorage) CreateData(user string, data models.StoredData) error {
	id, err := s.GetUserID(user)
	if err != nil || id == 0 {
		s.logger.Error("failed to get user id", zap.Error(err))
		return err
	}
	data.UUID = uuid.New().String()
	_, err = s.db.Exec("INSERT INTO data (user_id, uuid, meta, type, data) VALUES (?, ?, ?, ?, ?)", id, data.UUID, data.Meta, data.DataType, data.EncryptedData)

	if err != nil {
		s.logger.Error("failed to insert data", zap.Error(err))
		return err
	}
	return nil
}

func (s *ClientStorage) GetDataByUUID(user string, uuid string) (*models.StoredData, error) {
	id, err := s.GetUserID(user)
	if err != nil || id == 0 {
		s.logger.Error("failed to get user id", zap.Error(err))
		return nil, err
	}
	row := s.db.QueryRow("SELECT uuid, meta, type, data FROM data WHERE uuid = ? AND user_id = ?", uuid, id)

	var data models.StoredData
	err = row.Scan(&data.UUID, &data.Meta, &data.DataType, &data.EncryptedData)
	if err != nil {
		s.logger.Error("failed to scan data", zap.Error(err))
		return nil, err
	}
	return &data, nil
}

func (s *ClientStorage) GetData(user string) ([]models.StoredData, error) {
	id, err := s.GetUserID(user)
	if err != nil || id == 0 {
		s.logger.Error("failed to get user id", zap.Error(err))
		return nil, err
	}
	rows, err := s.db.Query("SELECT uuid, meta, type, data FROM data WHERE user_id = ?", id)
	if err != nil {
		s.logger.Error("failed to select data", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	var data []models.StoredData
	for rows.Next() {
		var d models.StoredData
		err := rows.Scan(&d.UUID, &d.Meta, &d.DataType, &d.EncryptedData)
		if err != nil {
			s.logger.Error("failed to scan data", zap.Error(err))
			return nil, err
		}
		data = append(data, d)
	}
	return data, nil
}

func (s *ClientStorage) UpdateData(user string, data models.StoredData) error {
	id, err := s.GetUserID(user)
	if err != nil || id == 0 {
		s.logger.Error("failed to get user id", zap.Error(err))
		return err
	}
	_, err = s.db.Exec("UPDATE data SET meta = ?, data = ? WHERE uuid = ? AND user_id = ?", data.Meta, data.EncryptedData, data.UUID, id)

	if err != nil {
		s.logger.Error("failed to update data", zap.Error(err))
		return err
	}
	return nil
}

func (s *ClientStorage) DeleteData(user string, data models.StoredData) error {
	id, err := s.GetUserID(user)
	if err != nil || id == 0 {
		s.logger.Error("failed to get user id", zap.Error(err))
		return err
	}
	_, err = s.db.Exec("DELETE FROM data WHERE uuid = ? AND user_id = ?", data.UUID, id)

	if err != nil {
		s.logger.Error("failed to delete data", zap.Error(err))
		return err
	}
	return nil
}

func (s *ClientStorage) DeleteAllData(user string) error {
	id, err := s.GetUserID(user)
	if err != nil || id == 0 {
		s.logger.Error("failed to get user id", zap.Error(err))
		return err
	}
	_, err = s.db.Exec("DELETE FROM data WHERE user_id = ?", id)

	if err != nil {
		s.logger.Error("failed to delete all data", zap.Error(err))
		return err
	}
	return nil
}

func (s *ClientStorage) FindByMeta(user string, meta string) ([]models.StoredData, error) {
	id, err := s.GetUserID(user)
	if err != nil || id == 0 {
		s.logger.Error("failed to get user id", zap.Error(err))
		return nil, err
	}
	rows, err := s.db.Query("SELECT uuid, meta,type, data FROM data WHERE meta = ? AND user_id = ?", meta, id)
	if err != nil {
		s.logger.Error("failed to select data", zap.Error(err))
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			s.logger.Error("failed to close rows", zap.Error(err))
		}
	}(rows)
	var data []models.StoredData
	for rows.Next() {
		var d models.StoredData
		err := rows.Scan(&d.UUID, &d.Meta, &d.EncryptedData)
		if err != nil {
			s.logger.Error("failed to scan data", zap.Error(err))
			return nil, err
		}
		data = append(data, d)
	}
	return data, nil
}
