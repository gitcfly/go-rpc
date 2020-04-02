package server

import (
	"fmt"

	"github.com/gitcfly/go-rpc/inter"
)

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
