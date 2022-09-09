package mao

type SceneDebug struct {
	ApiID          int64         `bson:"api_id"`
	APIName        string        `bson:"api_name"`
	Assertion      []*Assertion  `bson:"assertion"`
	EventID        string        `bson:"event_id"`
	NextList       []string      `bson:"next_list"`
	Regex          []*DebugRegex `bson:"regex"`
	RequestBody    string        `bson:"request_body"`
	RequestCode    int64         `bson:"request_code"`
	RequestHeader  string        `bson:"request_header"`
	ResponseBody   string        `bson:"response_body"`
	ResponseBytes  int64         `bson:"response_bytes"`
	ResponseHeader string        `bson:"response_header"`
	Status         string        `bson:"status"`
	UUID           string        `bson:"uuid"`
}

type DebugRegex struct {
	Code string `json:"code"`
}
