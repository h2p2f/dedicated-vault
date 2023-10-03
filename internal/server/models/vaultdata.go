// Package: models
// in this fale we have models for data
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// VaultData is a struct for data
type VaultData struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	DataUUID string             `json:"data_uuid" bson:"dataUUID"`
	UserUUID string             `json:"user_uuid" bson:"userUUID"`
	Meta     string             `json:"meta" bson:"meta"`
	DataType string             `json:"data_type" bson:"dataType"`
	Data     []byte             `json:"data" bson:"data"`
	Created  int64              `json:"created" bson:"created"`
	Updated  int64              `json:"updated,omitempty" bson:"updated"`
}
