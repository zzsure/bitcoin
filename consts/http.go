package consts

const (
	ERR_NO_SUCC          = 0   // 成功响应
	HTTP_SUCC_CODE       = 200 // http请求成功码
	ERR_NO_CLIENT_COMMON = 400 // 客户端有误
	ERR_NO_MISSING_PARAM = 401 // 缺少请求参数
	ERR_NO_INVALID_PARAM = 402 // 请求参数非法
	ERR_NO_SYSTEM_COMMON = 500 // 系统错误
	ERR_NO_DATABASE      = 501 // 数据库错误
)

var errMsgMap = map[int]string{
	ERR_NO_SUCC:          "请求成功",
	HTTP_SUCC_CODE:       "http请求成功",
	ERR_NO_SYSTEM_COMMON: "系统错误",
	ERR_NO_CLIENT_COMMON: "客户端有误",
	ERR_NO_MISSING_PARAM: "缺少请求参数",
	ERR_NO_INVALID_PARAM: "请求参数非法",
	ERR_NO_DATABASE:      "数据库错误",
}

func GetErrMsg(errNo int) string {
	errMsg, ok := errMsgMap[errNo]
	if !ok {
		errMsg = errMsgMap[ERR_NO_SYSTEM_COMMON]
	}
	return errMsg
}
