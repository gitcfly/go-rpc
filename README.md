# go-rpc
golang rpc框架

获取：go get -u github.com/gitcfly/go-rpc

#服务端教程
```
// 声明服务
type TestServer struct {
	CallName func(string, int) (string, int)
	CallData func(string)
}

//实例化服务
func NewTestServer() *inter.TestServer {
    return &inter.TestServer{
	CallData: func(s string) {
           fmt.Println("call callData by client")
	},
	CallName: func(s string, v int) (string, int) {
	   fmt.Println("call CallName by client")
	   return "you name is client v2 ", 9
	},
    }
}
```
服务端使用示例
```

// 1、初始化rpc服务端实例，参数1，psm,可以随意指定
rpcServer := rpc.NewRpcServer("ckj.gitcfly.rpc", ":8080")
// 2、注册服务
rpcServer.Service(server.NewTestServer())
// 3、 开始服务
rpcServer.Run()
```
#客户端示例
```
// 1、初始化客户端实例，参数psm要和服务端的一致
client := rpc.NewRpcClient("ckj.gitcfly.rpc", "http://127.0.0.1:8080")
// 2、注册服务
netServer := client.Client(&inter.TestServer{}).(*inter.TestServer)
// 3、调用服务方法1
strRes, intRes := netServer.CallName("hello , i am client , i want to know my name !", 999)
fmt.Println(strRes, intRes)
```
