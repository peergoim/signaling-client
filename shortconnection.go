package signaling_client

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type HttpConnection struct {
	client *http.Client
	c      *Config
}

func NewHttpConnection(c *Config) *HttpConnection {
	setLogLevel(c.LogLevel)
	return &HttpConnection{
		client: http.DefaultClient,
		c:      c,
	}
}

func (c *HttpConnection) Call(ctx context.Context, request *CallRequest) *CallResponse {
	if request.CallId == "" {
		request.CallId = GenerateCallId()
	}
	Debugf("call: %#v", request)
	body := bytes.NewBuffer(request.ToBytes())
	r, _ := http.NewRequest("POST", c.c.Endpoint+"/call", body)
	r.Header.Set("Content-Type", "application/json")
	select {
	case <-ctx.Done():
		Errorf("call deadline exceeded: %v", ctx.Err())
		return &CallResponse{
			CallId: request.CallId,
			Method: request.Method,
			Status: CodeDeadlineExceeded,
			Data:   nil,
		}
	default:
		do, err := c.client.Do(r)
		if err != nil {
			Errorf("call failed: %v", err)
			return &CallResponse{
				CallId: request.CallId,
				Method: request.Method,
				Status: CodeUnavailable,
				Data:   nil,
			}
		}
		defer do.Body.Close()
		all, err := io.ReadAll(do.Body)
		Debugf("call response: %s", string(all))
		if err != nil {
			return &CallResponse{
				CallId: request.CallId,
				Method: request.Method,
				Status: CodeUnavailable,
				Data:   nil,
			}
		}
		response := &CallResponse{}
		_ = response.FromBytes(all)
		return response
	}
}
