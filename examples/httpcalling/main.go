package main

import (
	"context"
	"fmt"
	signaling_client "github.com/peergoim/signaling-client"
	"strconv"
	"time"
)

func main() {
	c := signaling_client.NewHttpConnection(&signaling_client.Config{
		PeerId:   "",
		Endpoint: "http://localhost:31134",
		LogLevel: "debug",
	})
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	response := c.Call(ctx, &signaling_client.CallRequest{
		PeerId: "123",
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
