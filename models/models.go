package models


// type Netflix struct{
// 	Id primitive.ObjectID `json:"_id,omitempty" , bson:"_id,omitempty" `
// 	Moive string  `json:"moive,omitempty" `
// 	watched bool `json:"watched,omitempty"`
// }



import "go.mongodb.org/mongo-driver/bson/primitive"

type Netflix struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Movie    string             `json:"movie,omitempty"`
	Watched  bool               `json:"watched,omitempty"`
	Director string             `json:"director,omitempty"`
}



