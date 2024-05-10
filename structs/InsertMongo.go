package structs

import "go.mongodb.org/mongo-driver/bson/primitive"

type InsertKonsumen struct {
	IdKonsumen string `bson:"title"`
	Operation  string `bson:"content"`
}

type PostInsertKonsumen struct {
	Id         primitive.ObjectID `bson:"_id"`
	IdKonsumen string             `bson:"title"`
	Operation  string             `bson:"content"`
}
