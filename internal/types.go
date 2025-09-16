package internal

type Result struct {
	//5 代表服务带错误
	//4 代表客户端错误
	Code    int //错误码
	Message string
}
