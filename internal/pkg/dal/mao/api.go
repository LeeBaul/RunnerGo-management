package mao

type API struct {
	TargetID    int64  `bson:"target_id"`
	URL         string `bson:"url"`
	Header      string `bson:"header"`
	Query       string `bson:"query"`
	Body        string `bson:"body"`
	Auth        string `bson:"auth"`
	Description string `bson:"description"`
}
