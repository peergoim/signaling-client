package signaling_client

import (
	"context"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

type Config struct {
	PeerId   string
	Endpoint string
	LogLevel string // debug, info, warn, error
}

func (c *Config) Url() string {
	return c.Endpoint + "/ws?peerId=" + c.PeerId
}

func (c *ServerConnection) onReply(response *CallResponse) {
	defer func() {
		if err := recover(); err != nil {
			Errorf("onReply panic: %v", err)
		}
	}()
	// 1. 从callResponseChannel中获取响应channel
	{
		ch, ok := c.callResponseChannel.Load(response.CallId)
		if !ok {
			return
		}
		// 2. 发送响应
		ch.(chan *CallResponse) <- response
	}
}

func (c *ServerConnection) registerCallResponseChannel(id string, ch chan *CallResponse) {
	c.callResponseChannel.Store(id, ch)
}

func (c *ServerConnection) unregisterCallResponseChannel(id string) {
	c.callResponseChannel.Delete(id)
}

func (c *ServerConnection) Call(ctx context.Context, request *CallRequest) *CallResponse {
	if request.CallId == "" {
		request.CallId = GenerateCallId()
	}
	Debugf("call: %#v", request)
	ch := make(chan *CallResponse, 1)
	c.registerCallResponseChannel(request.CallId, ch)
	defer c.unregisterCallResponseChannel(request.CallId)
	// 发送请求
	c.l.RLock()
	conn := c.conn
	c.l.RUnlock()
	if conn == nil {
		return &CallResponse{
			CallId: request.CallId,
			Method: request.Method,
			Status: CodeUnavailable,
			Data:   nil,
		}
	}
	_ = conn.Write(ctx, websocket.MessageBinary, request.ToBytes())
	Debugf("send request finished")
	// 等待响应
	select {
	case <-ctx.Done():
		return &CallResponse{
			CallId: request.CallId,
			Method: request.Method,
			Status: CodeDeadlineExceeded,
			Data:   nil,
		}
	case response := <-ch:
		return response
	}
}

type ServerConnection struct {
	config              *Config
	conn                *websocket.Conn
	l                   sync.RWMutex
	methodHandlers      map[string]MethodHandler
	callResponseChannel sync.Map
}

func (c *ServerConnection) reconnect() {
	c.l.Lock()
	if c.conn != nil {
		c.conn.Close(websocket.StatusNormalClosure, "")
		c.conn = nil
	}
	c.l.Unlock()
	for {
		conn, _, err := websocket.Dial(context.Background(), c.config.Url(), nil)
		if err != nil {
			Errorf("failed to connect to server: %v, retry after 1s", err)
			time.Sleep(time.Second)
			continue
		}
		c.l.Lock()
		c.conn = conn
		go c.loopRead()
		c.l.Unlock()
		break
	}
}

func (c *ServerConnection) loopRead() {
	for {
		c.l.RLock()
		conn := c.conn
		c.l.RUnlock()
		if conn == nil {
			return
		}
		typ, msg, err := conn.Read(context.Background())
		if err != nil {
			Errorf("failed to read message: %v", err)
			c.reconnect()
			return
		}
		if typ == websocket.MessageBinary {
			//请求来了
			Debugf("receive request: %v", string(msg))
			request := &CallRequest{}
			err = request.FromBytes(msg)
			if err != nil {
				Errorf("failed to unmarshal request: %v", err)
				continue
			}
			var response *CallResponse
			handler, ok := c.methodHandlers[request.Method]
			if !ok {
				// notfound
				response = &CallResponse{
					CallId: request.CallId,
					Method: request.Method,
					Status: CodeNotFound,
					Data:   nil,
				}
			} else {
				response = handler(request)
			}
			// return response
			c.l.RLock()
			w := c.conn
			c.l.RUnlock()
			if w == nil {
				return
			}
			_ = w.Write(context.Background(), websocket.MessageText, response.ToBytes())
		} else if typ == websocket.MessageText {
			//响应来了
			response := &CallResponse{}
			err = response.FromBytes(msg)
			if err != nil {
				Errorf("failed to unmarshal response: %v", err)
				continue
			}
			c.onReply(response)
		}
	}
}

type MethodHandler func(request *CallRequest) *CallResponse

func RegisterPeer(c *Config, methodHandlers map[string]MethodHandler) *ServerConnection {
	if len(methodHandlers) == 0 {
		panic("no method handler")
	}
	setLogLevel(c.LogLevel)
	s := &ServerConnection{
		config:         c,
		methodHandlers: methodHandlers,
	}
	go s.reconnect()
	return s
}
