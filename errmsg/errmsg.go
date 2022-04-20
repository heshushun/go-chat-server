package errmsg

const (
	ErrUsernameExist = iota
	ErrMysqlServer
	ErrPassword
	InvalidToken
	AuthEmpty
	TokenRunTimeError
	FileTooBig
	Success = 200
	Error   = 300
)

var CodeMsg = map[int]string{
	ErrUsernameExist:  "用户名存在",
	ErrMysqlServer:    "数据库处理错误",
	ErrPassword:       "密码错误",
	InvalidToken:      "非法的token",
	AuthEmpty:         "首部为空",
	TokenRunTimeError: "token过期",
	FileTooBig:        "文件过大",
	Success:           "成功",
	Error:             "失败",
}
