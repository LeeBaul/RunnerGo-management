package mao

type Scene struct {
	TargetID int64  `bson:"target_id"`
	Request  string `bson:"request"`
	Script   string `bson:"script"`
}
