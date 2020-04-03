package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gitcfly/go-rpc/tools"
	"github.com/gitcfly/tryct/try"
	"github.com/kirinlabs/HttpRequest"
)

var PsmService = make(map[string]*RpcServer)

type RpcServer struct {
	Psm    string
	Addr   string
	SerMap map[string]interface{}
}

func NewRpcServer(psm string, addr string) *RpcServer {
	rpcServer := &RpcServer{
		Psm:    psm,
		Addr:   addr,
		SerMap: map[string]interface{}{},
	}
	PsmService[psm] = rpcServer
	return rpcServer
}

func (rpcServer *RpcServer) Run() {
	go func() {
		gin.SetMode("release")
		g := gin.Default()
		g.Any("/", httpServer)
		g.Run(rpcServer.Addr)
	}()
}

func httpServer(c *gin.Context) {
	try.Try(func() {
		var rpcRequest RpcReqest
		c.BindJSON(&rpcRequest)
		result := invoke(rpcRequest)
		c.JSON(http.StatusOK, result)
	}).Catch(func(err interface{}) {
		fmt.Printf("服务端调用异常,err: %v \n", err)
	})
}

type RpcClient struct {
	Psm  string
	Addr string
}

func NewRpcClient(psm string, addr string) *RpcClient {
	return &RpcClient{Psm: psm, Addr: addr}
}

type RpcResponse struct {
	Outs []interface{}
}

type RpcReqest struct {
	Psm   string
	Path  string
	Fname string
	Args  []interface{}
}

func (rpcServer *RpcServer) Service(obj interface{}) {
	vt := reflect.ValueOf(obj).Elem().Type()
	path := vt.PkgPath() + "/" + vt.Name()
	rpcServer.SerMap[path] = obj
}

func invoke(req RpcReqest) RpcResponse {
	var rpcResp RpcResponse
	try.Try(func() {
		instance := PsmService[req.Psm].SerMap[req.Path]
		ve := reflect.ValueOf(instance).Elem()
		mt := ve.FieldByName(req.Fname)
		ft := mt.Type()
		var args []reflect.Value
		for idx, arg := range req.Args {
			temp := reflect.New(ft.In(idx)).Interface()
			bytes, _ := json.Marshal(arg)
			json.Unmarshal(bytes, &temp)
			args = append(args, reflect.ValueOf(temp).Elem())
		}
		result := mt.Call(args)
		var outs []interface{}
		for _, e := range result {
			outs = append(outs, e.Interface())
		}
		rpcResp.Outs = outs
	}).Catch(func(err interface{}) {
		fmt.Printf("服务端调用异常，request: %v ,err: %v \n", req, err)
	})
	return rpcResp
}

func (client *RpcClient) Client(obj interface{}) interface{} {
	objValue := reflect.ValueOf(obj).Elem()
	objType := objValue.Type()
	pkgPath := objValue.Type().PkgPath() + "/" + objValue.Type().Name()
	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Field(i)
		method := objType.Field(i).Name
		fun := reflect.MakeFunc(field.Type(), createFunc(client, field.Type(), pkgPath, method))
		field.Set(fun)
	}
	return obj
}

func createFunc(client *RpcClient, funcType reflect.Type, pkg string, method string) func(args []reflect.Value) (results []reflect.Value) {
	return func(args []reflect.Value) (results []reflect.Value) {
		try.Try(func() {
			var interArgs []interface{}
			for _, arg := range args {
				interArgs = append(interArgs, arg.Interface())
			}
			req := HttpRequest.NewRequest()
			rpcR := RpcReqest{Psm: client.Psm, Path: pkg, Fname: method, Args: interArgs}
			response, _ := req.Post(client.Addr, tools.ToJsonString(rpcR))
			resBody, _ := response.Body()
			var rpcResp RpcResponse
			json.Unmarshal(resBody, &rpcResp)
			for idx, out := range rpcResp.Outs {
				bytes, _ := json.Marshal(out)
				rOut := reflect.New(funcType.Out(idx)).Interface()
				json.Unmarshal(bytes, &rOut)
				results = append(results, reflect.ValueOf(rOut).Elem())
			}
		}).Catch(func(err interface{}) {
			fmt.Printf("call func error ,function: %v ,error: %v \n", funcType, err)
		})
		return results
	}
}
