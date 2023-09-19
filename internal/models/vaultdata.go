package models

import "time"

type VaultData struct {
	DataUUID string    `json:"data_uuid" bson:"dataUUID"`
	UserUUID string    `json:"user_uuid" bson:"userUUID"`
	Meta     string    `json:"meta" bson:"meta"`
	Data     []byte    `json:"data" bson:"data"`
	Created  time.Time `json:"created" bson:"created"`
	Updated  time.Time `json:"updated" bson:"updated"`
}
