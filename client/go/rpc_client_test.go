package client

import (
	"encoding/json"
	"fmt"
	"testing"

	rpc_serivce "cybervein.org/CyberveinDB/grpc/proto"
)

func TestAuth(t *testing.T) {
	//after server started
	client, e := NewRpcClient("127.0.0.1:40001", "App@1234")
	if e != nil {
		t.Error(e)
	}
	token := client.GetContextWithToken()
	fmt.Println(token)
}

func TestExecute(t *testing.T) {
	//after server started
	client, e := NewRpcClient("127.0.0.1:40001", "App@1234")
	if e != nil {
		t.Error(e)
	}
	request := &rpc_serivce.CommandRequest{Cmd: "set k v"}
	response, e := client.Execute(request)
	if e != nil {
		t.Error(e)
	}
	responseData, _ := json.Marshal(response)
	fmt.Printf(string(responseData))
}
