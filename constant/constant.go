package constant

// 全局错误码
const (
	GLOBAL_SUCCESS = 1000 //操作成功

	GLOBAL_NO_AUTH     = 1001 //没有权限
	GLOBAL_NO_AUTH_MSG = "没有权限"

	GLOBAL_SYS_ERR     = 1002 //系统错误
	GLOBAL_SYS_ERR_MSG = "系统错误"

	GLOBAL_PARM_ERR     = 1003 //参数错误
	GLOBAL_PARM_ERR_MSG = "参数错误"

	GLOBAL_NO_SESSION     = 1005 //没有Session
	GLOBAL_NO_SESSION_MSG = "身份未验证"
)
