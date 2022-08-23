package mao

import "go.mongodb.org/mongo-driver/bson"

type API struct {
	TargetID    int64    `bson:"target_id"`
	URL         string   `bson:"url"`
	Header      bson.Raw `bson:"header"`
	Query       bson.Raw `bson:"query"`
	Body        bson.Raw `bson:"body"`
	Auth        bson.Raw `bson:"auth"`
	Description string   `bson:"description"`
}
