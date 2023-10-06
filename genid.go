package signaling_client

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/bwmarrin/snowflake"
)

var (
	c           = make(chan string)
	nd          *snowflake.Node
	nodId       int64
	uniquePodId string
)

func loop() {
	for {
		c <- nd.Generate().String()
	}
}

// GenerateCallId 生成一个唯一的callId
func GenerateCallId() string {
	h := md5.New()
	s := uniquePodId + <-c
	h.Write([]byte(s))
	cipher := h.Sum(nil)
	return hex.EncodeToString(cipher)
}
