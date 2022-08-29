package mao

import "go.mongodb.org/mongo-driver/bson"

type Flow struct {
	SceneID int64    `bson:"scene_id"`
	TeamID  int64    `bson:"team_id"`
	Version int32    `bson:"version"`
	Flows   bson.Raw `bson:"flows"`
}
