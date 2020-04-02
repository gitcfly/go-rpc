# go-rpc
golang rpc框架

#服务端教程
```
  
	// 1、初始化rpc服务端实例，参数1，psm,可以随意指定
	rpcServer := rpc.NewRpcServer("ckj.gitcfly.rpc", ":8080")
	// 2、注册服务
	rpcServer.Service(server.NewTestServer())
	// 3、 开始服务
	rpcServer.Run()
```
#客户端教程
```
	// 1、初始化客户端实例，参数psm要和服务端的一致
	client := rpc.NewRpcClient("ckj.gitcfly.rpc", "http://127.0.0.1:8080")
	// 2、注册服务
	netServer := client.Client(&inter.TestServer{}).(*inter.TestServer)
	// 3、调用服务方法1
	strRes, intRes := netServer.CallName("hello , i am client , i want to know my name !", 999)
	fmt.Println(strRes, intRes)
```
