package models

//import (
//	pb "github.com/h2p2f/dedicated-vault/proto"
//)
//
//type Data struct {
//	UUID        string      `json:"uuid" bson:"UUID"`
//	Meta        string      `json:"meta" bson:"meta"`
//	DataType    string      `json:"data_type" bson:"dataType"`
//	Card        CreditCard  `json:"card" bson:"card"`
//	Credentials Credentials `json:"credentials" bson:"credentials"`
//	Text        TextData    `json:"text_data" bson:"text"`
//	Binary      BinaryData  `json:"binary_data" bson:"binary"`
//}
//
//type CreditCard struct {
//	Number     string `json:"number" bson:"number"`
//	ExpireDate string `json:"expire_date" bson:"expireDate"`
//	CVV        string `json:"cvv" bson:"CVV"`
//}
//
//type Credentials struct {
//	Login    string `json:"login" bson:"login"`
//	Password string `json:"password" bson:"password"`
//}
//
//type TextData struct {
//	Text string `json:"text" bson:"text"`
//}
//
//type BinaryData struct {
//	name string `json:"name" bson:"name"`
//	data []byte `json:"data" bson:"data"`
//}
//
//func (d *Data) FromPB(pb *pb.Data) {
//	d.Meta = pb.Meta
//	d.Card.Number = pb.GetCreditCard().Number
//	d.Card.ExpireDate = pb.GetCreditCard().Expiration
//	d.Card.CVV = pb.GetCreditCard().Cvv
//	d.Credentials.Login = pb.GetCredentials().Username
//	d.Credentials.Password = pb.GetCredentials().Password
//	d.Text.Text = pb.GetText().Text
//	d.Binary.name = pb.GetBinary().Name
//	d.Binary.data = pb.GetBinary().Data
//}
//
//func (d *Data) ToPB() *pb.Data {
//	switch d.DataType {
//	case "credit_card":
//		return &pb.Data{
//			Meta: d.Meta,
//			Data: &pb.Data_CreditCard{
//				CreditCard: &pb.CreditCard{
//					Number:     d.Card.Number,
//					Expiration: d.Card.ExpireDate,
//					Cvv:        d.Card.CVV,
//				},
//			},
//		}
//	case "credentials":
//		return &pb.Data{
//			Meta: d.Meta,
//			Data: &pb.Data_Credentials{
//				Credentials: &pb.Credentials{
//					Username: d.Credentials.Login,
//					Password: d.Credentials.Password,
//				},
//			},
//		}
//	case "text":
//		return &pb.Data{
//			Meta: d.Meta,
//			Data: &pb.Data_Text{
//				Text: &pb.TextData{
//					Text: d.Text.Text,
//				},
//			},
//		}
//	case "binary":
//		return &pb.Data{
//			Meta: d.Meta,
//			Data: &pb.Data_Binary{
//				Binary: &pb.BinaryData{
//					Name: d.Binary.name,
//					Data: d.Binary.data,
//				},
//			},
//		}
//	default:
//		return &pb.Data{}
//	}
//
//}
