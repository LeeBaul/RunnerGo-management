package mao

import "go.mongodb.org/mongo-driver/bson"

type Flow struct {
	SceneID int64 `bson:"scene_id"`
	TeamID  int64 `bson:"team_id"`
	Version int32 `bson:"version"`
	//Flows   string `bson:"flows"`
	Nodes           bson.Raw `bson:"nodes"`
	Edges           bson.Raw `bson:"edges"`
	MultiLevelNodes bson.Raw `bson:"multi_level_nodes"`
}
