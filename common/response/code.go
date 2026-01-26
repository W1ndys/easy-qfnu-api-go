package response

// 业务状态码常量
const (
	CodeSuccess          = 200  // 成功
	CodeServerBusy       = 1    // 系统繁忙/通用错误
	CodeInvalidParam     = 1001 // 参数错误
	CodeAuthExpired      = 401  // 缺少 Authorization 或 Authorization 过期
	CodeResourceNotFound = 404  // 未查询到数据
	CodeTargetError      = 502  // 教务系统挂了
)

// MsgFlags 状态码对应的默认提示信息
var MsgFlags = map[int]string{
	CodeSuccess:          "success",
	CodeServerBusy:       "系统繁忙，请稍后再试",
	CodeInvalidParam:     "请求参数错误",
	CodeAuthExpired:      "缺少 Authorization 字段 或 Authorization 过期，请重新获取相关系统的Cookie，获取方法参考 https://mp.weixin.qq.com/s/zFK9c4ecpGdRwXSKzaVFnw",
	CodeResourceNotFound: "未查询到数据，请调整查询条件后重试",
	CodeTargetError:      "目标系统无响应",
}

// GetMsg 获取状态码对应的消息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[CodeServerBusy]
}
