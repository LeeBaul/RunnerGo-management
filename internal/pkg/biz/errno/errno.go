// Package errno 定义所有错误码
package errno

const (
	Ok                         = 0
	ParamError                 = 10001
	ServerError                = 10002
	NonceError                 = 10003
	TimeStampError             = 10004
	RPCFailed                  = 10005
	InvalidToken               = 10006
	MarshalFailed              = 10007
	UnMarshalFailed            = 10008
	AvailableIntegralNotEnough = 10009
	IntegralOperTypeError      = 10010
	MustDID                    = 10011
	MustSN                     = 10012
	HttpFailed                 = 10013
	RedisOperFailed            = 10100
	MongoOperFailed            = 10101
	MysqlOperFailed            = 10102
	RecordNotFound             = 10103
	SignError                  = 20001
	RepeatRequest              = 20002
	MustLogin                  = 20003
)

// CodeMsgMap 错误码映射错误信息，不展示给用户
var CodeMsgMap = map[int]string{
	Ok:                         "success",
	ServerError:                "internal server error",
	ParamError:                 "param error",
	SignError:                  "signature error",
	RepeatRequest:              "repeat request",
	NonceError:                 "nonce error",
	TimeStampError:             "timestamp error",
	RecordNotFound:             "record not found",
	RPCFailed:                  "rpc failed",
	InvalidToken:               "invalid token",
	MarshalFailed:              "marshal failed",
	UnMarshalFailed:            "unmarshal failed",
	RedisOperFailed:            "redis operate failed",
	MongoOperFailed:            "mongo operate failed",
	MysqlOperFailed:            "mysql operate failed",
	AvailableIntegralNotEnough: "available integral not enough",
	IntegralOperTypeError:      "integral operate type error",
	MustLogin:                  "must login",
	MustDID:                    "must DID",
	MustSN:                     "must SN",
	HttpFailed:                 "http failed",
}

// CodeAlertMap 错图码映射错误提示，展示给用户
var CodeAlertMap = map[int]string{
	Ok:                         "成功",
	ServerError:                "服务器错误",
	ParamError:                 "参数校验错误",
	SignError:                  "签名错误",
	RepeatRequest:              "重放请求",
	NonceError:                 "_nonce参数错误",
	TimeStampError:             "_timestamp参数错误",
	RecordNotFound:             "数据库记录不存在",
	RPCFailed:                  "请求下游服务失败",
	InvalidToken:               "无效的token",
	MarshalFailed:              "序列化失败",
	UnMarshalFailed:            "反序列化失败",
	RedisOperFailed:            "redis操作失败",
	MongoOperFailed:            "mongo操作失败",
	MysqlOperFailed:            "mysql操作失败",
	AvailableIntegralNotEnough: "积分余额不足",
	IntegralOperTypeError:      "积分操作类型错误",
	MustLogin:                  "没有获取到登录态",
	MustDID:                    "缺少设备DID信息",
	MustSN:                     "缺少设备SN信息",
	HttpFailed:                 "请求下游Http服务失败",
}
