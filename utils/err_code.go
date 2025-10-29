package utils

// 用户模块错误码
const (
	ErrCodeUsernameTaken      = 1001
	ErrCodeInvalidPassword    = 1002
	ErrCodeUserNotFound       = 1003
	ErrCodeTokenNotFound      = 1004
	ErrCodeTokenExpired       = 1005
	ErrCodeInvalidToken       = 1006
	ErrCodeInvalidTokenFormat = 1007
	ErrCodePermissionDenied   = 1008
	ErrCodeTokenGenerate      = 1009
	ErrCodeInvalidEmailCode   = 1010
	ErrCodeEmailTaken         = 1011
	ErrCodeCreateUser         = 1012
	ErrCodeGetUserList        = 1013
	ErrCodeEmailCodeGenerate  = 1014
	ErrCodeSendEmailCode      = 1015
	ErrCodeUpdateUser         = 1016
	ErrCodeDeleteUser         = 1017
	ErrCodeWechatLogin        = 1018
	ErrCodeWechatBindFailed   = 1019
	ErrCodeInvalidParams      = 1020
)

// other
const (
	SUCCESS                     = 1
	InvalidParameter            = 0
	ErrCodeFullUpdateRunning    = 2000
	ErrCodeGetSettings          = 2001
	ErrCodeGetGoodsTotal        = 2002
	ErrCodeGetGoods             = 2003
	ErrCodeUpdateSetting        = 2004
	ErrCodeUpdateUUToken        = 2005
	ErrCodeUpdateBuffToken      = 2006
	ErrCodeGetTokenExpired      = 2007
	ErrCodeCreateDefaultSetting = 2008
	ErrCodeGetGoodsCategory     = 2009
)

// 错误码与消息映射
var errorCodeToMessage = map[int]string{
	SUCCESS:                  "success",
	ErrCodeFullUpdateRunning: "Full update running",
	ErrCodeGetSettings:       "Get settings error",
	ErrCodeGetGoodsTotal:     "Get goods total error",
	ErrCodeGetGoods:          "Get goods data error",
	InvalidParameter:         "Invalid Parameter",
	ErrCodeUpdateSetting:     "Update setting error",
	ErrCodeUpdateUUToken:     "Update UU token error",
	ErrCodeUpdateBuffToken:   "Update buff token error",
	ErrCodeGetTokenExpired:   "Get token expired error",
	ErrCodeGetGoodsCategory:  "Get goods category error",
	// 用户模块
	ErrCodeInvalidEmailCode:   "The provided email code is incorrect",
	ErrCodeUsernameTaken:      "The requested username is already in use",
	ErrCodeEmailTaken:         "The requested email is already in use",
	ErrCodeInvalidPassword:    "The provided password is incorrect",
	ErrCodeUserNotFound:       "User account not found",
	ErrCodeTokenNotFound:      "Authentication token not found",
	ErrCodeTokenExpired:       "Authentication token has expired",
	ErrCodeInvalidToken:       "Invalid authentication token",
	ErrCodeInvalidTokenFormat: "Malformed authentication token",
	ErrCodePermissionDenied:   "Insufficient permissions for this operation",
	ErrCodeTokenGenerate:      "Generate token error",
	ErrCodeCreateUser:         "Register user error",
	ErrCodeGetUserList:        "Get user list error",
	ErrCodeEmailCodeGenerate:  "Email code generate error",
	ErrCodeSendEmailCode:      "Send email code error",
	ErrCodeUpdateUser:         "Update user error",
	ErrCodeDeleteUser:         "Delete user error",
	ErrCodeWechatLogin:        "Wechat login error",
	ErrCodeWechatBindFailed:   "Wechat bind failed",
	ErrCodeInvalidParams:      "Invalid parameters",
}

// ErrorMessage 返回指定错误码对应的错误消息
func ErrorMessage(code int) string {
	if msg, exists := errorCodeToMessage[code]; exists {
		return msg
	}
	return "Unknown error"
}
