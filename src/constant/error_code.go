package constant

type ErrorStruct struct {
	Code string
	Msg  string
}

/*
 错误码前2位是模块编号，3-4位是领域，5-6表示错误类型，7-9表示具体错误
	模块编码：11
 	领域划分：
		用户管理：11
		财务管理：12
	错误类型：00-20是系统错误，21-99是业务错误
		未知系统错误：00
		数据库相关错误：01
		缓存相关错误：02

		未知业务错误：21
*/

var InsertDBError = ErrorStruct{"111101001", "插入数据库失败"}
