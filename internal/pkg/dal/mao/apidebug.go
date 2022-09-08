package mao

import "go.mongodb.org/mongo-driver/bson"

type APIDebug struct {
	ApiID                 int64    `bson:"api_id"`
	APIName               string   `bson:"api_name"`
	Assertion             bson.Raw `bson:"assertion"`
	EventID               string   `bson:"event_id"`
	Regex                 bson.Raw `bson:"regex"`
	RequestBody           bson.Raw `bson:"request_body"`
	RequestCode           int64    `bson:"request_code"`
	RequestHeader         string   `bson:"request_header"`
	RequestTime           int64    `bson:"request_time"`
	ResponseBody          bson.Raw `bson:"response_body"`
	ResponseBytes         int64    `bson:"response_bytes"`
	ResponseHeader        string   `bson:"response_header"`
	ResponseTime          string   `bson:"response_time"`
	ResponseLen           int32    `bson:"response_len"`
	ResponseStatusMessage string   `bson:"response_status_message"`
	UUID                  string   `bson:"uuid"`
}

type Assertion struct {
	Code      int    `bson:"code"`
	IsSucceed bool   `bson:"isSucceed"`
	Msg       string `bson:"msg"`
}
