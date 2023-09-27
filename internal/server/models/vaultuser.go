package models

import (
	pb "github.com/h2p2f/dedicated-vault/proto"
)

type User struct {
	UUID              string `json:"uuid" bson:"UUID"`
	Login             string `json:"login" bson:"login"`
	Password          string `json:"password" bson:"password"`
	LastServerUpdated int64  `json:"last_server_updated" bson:"lastServerUpdated"`
}

func (u *User) FromPB(pb *pb.User) {
	u.Login = pb.Name
	u.Password = pb.Password
}

func (u *User) ToPB() *pb.User {
	return &pb.User{
		Name:     u.Login,
		Password: u.Password,
	}
}
