package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
)

type FolderDataType string

var FolderData = struct {
	CreditCard  FolderDataType
	Credentials FolderDataType
	Text        FolderDataType
	Binary      FolderDataType
}{
	CreditCard:  "cc",
	Credentials: "cr",
	Text:        "tx",
	Binary:      "bn",
}

type Data struct {
	UUID     string         `json:"uuid"`
	Meta     string         `json:"meta"`
	DataType FolderDataType `json:"data_type"`
	Folder   Folder         `json:"data"`
}

type Folder struct {
	Card        CreditCard  `json:"card"`
	Credentials Credentials `json:"credentials"`
	Text        TextData    `json:"text"`
	Binary      BinaryData  `json:"binary"`
}

type CreditCard struct {
	Number     string `json:"number"`
	NameOnCard string `json:"name_on_card"`
	ExpireDate string `json:"expire_date"`
	CVV        string `json:"cvv"`
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TextData struct {
	Text string `json:"text"`
}

type BinaryData struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

func (d *Data) EncryptData(key []byte) (*StoredData, error) {
	key32 := sha256.Sum256(key)
	c, err := aes.NewCipher(key32[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	jsonFolder, err := json.Marshal(d.Folder)
	if err != nil {
		return nil, err
	}

	encryptedData := gcm.Seal(nonce, nonce, jsonFolder, nil)

	return &StoredData{
		UUID:          d.UUID,
		Meta:          d.Meta,
		DataType:      string(d.DataType),
		EncryptedData: encryptedData,
	}, nil
}

type StoredData struct {
	UUID          string `json:"uuid"`
	Meta          string `json:"meta"`
	DataType      string `json:"data_type"`
	EncryptedData []byte `json:"encrypted_data"`
}

func (s *StoredData) DecryptData(key []byte) (*Data, error) {
	key32 := sha256.Sum256(key)
	c, err := aes.NewCipher(key32[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(s.EncryptedData) < nonceSize {
		return nil, errors.New("encrypted data too short")
	}
	nonce, encryptedData := s.EncryptedData[:nonceSize], s.EncryptedData[nonceSize:]
	decryptedData, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}
	var data Data
	data.UUID = s.UUID
	data.Meta = s.Meta
	data.DataType = FolderDataType(s.DataType)
	var folder Folder
	err = json.Unmarshal(decryptedData, &folder)
	if err != nil {
		return nil, err
	}
	data.Folder = folder
	return &data, nil
}
