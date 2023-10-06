package main

import signaling_client "github.com/peergoim/signaling-client"

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

func main() {
	signaling_client.RegisterPeer(&signaling_client.Config{
		PeerId:   "123",
		Endpoint: "ws://localhost:31134",
		LogLevel: "debug",
	}, handers)
	select {}
}
