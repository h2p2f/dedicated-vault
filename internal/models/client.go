package models

import (
	"github.com/google/uuid"
	"time"
)

type Client struct {
	UUID     uuid.UUID `json:"uuid" bson:"UUID"`
	User     User      `json:"user" bson:"user"`
	LastSync time.Time `json:"last_sync" bson:"lastSync"`
}
