package response

// 业务状态码常量
const (
	CodeSuccess      = 0    // 成功
	CodeServerBusy   = 1    // 系统繁忙/通用错误
	CodeInvalidParam = 1001 // 参数错误
	CodeAuthExpired  = 401  // Token 过期
	CodeTargetError  = 502  // 教务系统挂了
)

// MsgFlags 状态码对应的默认提示信息
var MsgFlags = map[int]string{
	CodeSuccess:      "success",
	CodeServerBusy:   "系统繁忙，请稍后再试",
	CodeInvalidParam: "请求参数错误",
	CodeAuthExpired:  "身份认证已过期，请重新登录",
	CodeTargetError:  "目标系统无响应",
}

// GetMsg 获取状态码对应的消息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[CodeServerBusy]
}
