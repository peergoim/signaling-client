package main

import (
	"context"
	"fmt"
	signaling_client "github.com/peergoim/signaling-client"
	"strconv"
	"time"
)

func HelloHandler(req *signaling_client.CallRequest) *signaling_client.CallResponse {
	return &signaling_client.CallResponse{
		CallId: req.CallId,
		Method: req.Method,
		Status: 0,
		Data:   []byte("world"),
	}
}

var (
	handers = map[string]signaling_client.MethodHandler{
		"hello": HelloHandler,
	}
)

var c *signaling_client.ServerConnection

func init() {
	c = signaling_client.RegisterPeer(&signaling_client.Config{
		PeerId:   "1234",
		Endpoint: "ws://localhost:31134",
		LogLevel: "debug",
	}, handers)
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	time.Sleep(time.Second)
	response := c.Call(ctx, &signaling_client.CallRequest{
		PeerId: "telegram",
		CallId: strconv.FormatInt(time.Now().UnixNano(), 10),
		Method: "hello",
		Data:   []byte("hello"),
	})
	if response.Status == signaling_client.CodeOK {
		fmt.Printf("response: %#v\n", response)
	} else {
		fmt.Printf("response: error: %#v\n", response.Status.String())
	}
}
