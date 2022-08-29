package mao

import "go.mongodb.org/mongo-driver/bson"

type Scene struct {
	TargetID int64    `bson:"target_id"`
	Request  bson.Raw `bson:"request"`
	Script   bson.Raw `bson:"script"`
}
