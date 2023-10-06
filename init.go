package signaling_client

import (
	"github.com/bwmarrin/snowflake"
	"math/rand"
)

func init() {
	nodId = int64(rand.Intn(128))
	nd, _ = snowflake.NewNode(nodId)
	uniquePodId = nd.Generate().String()
	go loop()
	for k, v := range strToCode {
		codeToStr[v] = k
	}
}
