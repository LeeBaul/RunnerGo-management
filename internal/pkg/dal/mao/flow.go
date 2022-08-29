package mao

type Flow struct {
	SceneID int64  `bson:"scene_id"`
	TeamID  int64  `bson:"team_id"`
	Version int32  `bson:"version"`
	Flows   string `bson:"flows"`
}
