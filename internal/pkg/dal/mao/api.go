package mao

import (
	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/dal/rao"
)

type API struct {
	TargetID    int64    `bson:"target_id"`
	URL         string   `bson:"url"`
	Header      bson.Raw `bson:"header"`
	Query       bson.Raw `bson:"query"`
	Body        bson.Raw `bson:"body"`
	Auth        bson.Raw `bson:"auth"`
	Description string   `bson:"description"`
	Assert      bson.Raw `bson:"assert"`
	Regex       bson.Raw `bson:"regex"`
}

type Assert struct {
	Assert []*rao.Assert `bson:"assert"`
}

type Regex struct {
	Regex []*rao.Regex `bson:"regex"`
}
