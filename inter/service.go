package inter

// 声明服务
type TestServer struct {
	CallName func(string, int) (string, int)
	CallData func(string)
}
