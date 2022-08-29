package mao

import "go.mongodb.org/mongo-driver/bson"

type Group struct {
	TargetID int64    `bson:"target_id"`
	Request  bson.Raw `bson:"request"`
	Script   bson.Raw `bson:"script"`
}
