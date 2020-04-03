package main

import (
	"fmt"
	"time"

	"github.com/gitcfly/go-rpc/rpc"
)

// 声明服务，服务端声明，客户端引入
type TestServer struct {
	QueryName func(string, int) (string, int)
	QueryData func(string) string
}

// 服务端实现服务
func InstanceTestServer() *TestServer {
	testServer := &TestServer{}
	testServer.QueryData = func(s string) string {
		fmt.Println("call QueryData by client...")
		return "server return some data !"
	}
	testServer.QueryName = func(s2 string, i int) (s string, i2 int) {
		fmt.Println("call CallName by client")
		return "you name is client v2 ", i + 1
	}
	return testServer
}

// golang rpc框架
func main() {
	//服务端
	ServerSimple()

	time.Sleep(1 * time.Second)

	//客户端
	ClientSimple()
}

// 服务端教程
func ServerSimple() {
	// 1、初始化rpc服务端实例，参数1：psm,可以随意指定，参数2：服务地址ip:port形式
	rpcServer := rpc.NewRpcServer("ckj.gitcfly.rpc", ":8080")
	// 2、注册服务
	rpcServer.Service(InstanceTestServer())
	// 3、 开始服务
	rpcServer.Run()
}

// 以下是客户端教程
func ClientSimple() {
	// 1、初始化客户端实例，参数psm要和服务端的一致
	client := rpc.NewRpcClient("ckj.gitcfly.rpc", "127.0.0.1:8080")
	// 2、注册服务
	netServer := client.Client(&TestServer{}).(*TestServer)
	// 3、调用服务方法1
	strResp, intResp := netServer.QueryName("hello , i am client , i want to know my name !", 999)
	fmt.Println(strResp, intResp)
	// 3、调用服务方法2
	resp := netServer.QueryData("i want to get some data")
	fmt.Println(resp)
}
