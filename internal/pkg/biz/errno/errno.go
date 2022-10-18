// Package errno 定义所有错误码
package errno

const (
	Ok                 = 0
	ErrParam           = 10001
	ErrServer          = 10002
	ErrNonce           = 10003
	ErrTimeStamp       = 10004
	ErrRPCFailed       = 10005
	ErrInvalidToken    = 10006
	ErrMarshalFailed   = 10007
	ErrUnMarshalFailed = 10008
	ErrMustDID         = 10011
	ErrMustSN          = 10012
	ErrHttpFailed      = 10013
	ErrRedisFailed     = 10100
	ErrMongoFailed     = 10101
	ErrMysqlFailed     = 10102
	ErrRecordNotFound  = 10103
	ErrSignError       = 20001
	ErrRepeatRequest   = 20002
	ErrMustLogin       = 20003
	ErrAuthFailed      = 20004
	ErrYetRegister     = 20005
)

// CodeMsgMap 错误码映射错误信息，不展示给用户
var CodeMsgMap = map[int]string{
	Ok:                 "success",
	ErrServer:          "internal server error",
	ErrParam:           "param error",
	ErrSignError:       "signature error",
	ErrRepeatRequest:   "repeat request",
	ErrNonce:           "nonce error",
	ErrTimeStamp:       "timestamp error",
	ErrRecordNotFound:  "record not found",
	ErrRPCFailed:       "rpc failed",
	ErrInvalidToken:    "invalid token",
	ErrMarshalFailed:   "marshal failed",
	ErrUnMarshalFailed: "unmarshal failed",
	ErrRedisFailed:     "redis operate failed",
	ErrMongoFailed:     "mongo operate failed",
	ErrMysqlFailed:     "mysql operate failed",
	ErrMustLogin:       "must login",
	ErrMustDID:         "must DID",
	ErrMustSN:          "must SN",
	ErrHttpFailed:      "http failed",
	ErrAuthFailed:      "username/password failed",
	ErrYetRegister:     "email yet register",
}

// CodeAlertMap 错图码映射错误提示，展示给用户
var CodeAlertMap = map[int]string{
	Ok:                 "成功",
	ErrServer:          "服务器错误",
	ErrParam:           "参数校验错误",
	ErrSignError:       "签名错误",
	ErrRepeatRequest:   "重放请求",
	ErrNonce:           "_nonce参数错误",
	ErrTimeStamp:       "_timestamp参数错误",
	ErrRecordNotFound:  "数据库记录不存在",
	ErrRPCFailed:       "请求下游服务失败",
	ErrInvalidToken:    "无效的token",
	ErrMarshalFailed:   "序列化失败",
	ErrUnMarshalFailed: "反序列化失败",
	ErrRedisFailed:     "redis操作失败",
	ErrMongoFailed:     "mongo操作失败",
	ErrMysqlFailed:     "mysql操作失败",
	ErrMustLogin:       "没有获取到登录态",
	ErrMustDID:         "缺少设备DID信息",
	ErrMustSN:          "缺少设备SN信息",
	ErrHttpFailed:      "请求下游Http服务失败",
	ErrAuthFailed:      "用户名或密码错误",
	ErrYetRegister:     "用户已注册",
}
