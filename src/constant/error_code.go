package constant

type ErrorStruct struct {
	Code string
	Msg  string
}

// 错误码前四位是模块编号，5-8位是模块内的错误编号（5-6表示错误类型，7-8表示具体错误）
// var Success = ErrorStruct{"10020000", "成功"}

// var Error = ErrorStruct{"10020101", "未知错误"}

var InsertDBError = ErrorStruct{"10020201", "插入数据库失败"}
