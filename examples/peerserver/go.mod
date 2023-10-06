module github.com/peergoim/signaling-client/examples/peerserver

go 1.20

require github.com/peergoim/signaling-client v0.0.0-00010101000000-000000000000

require nhooyr.io/websocket v1.8.7 // indirect

require (
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/klauspost/compress v1.10.3 // indirect
)

replace github.com/peergoim/signaling-client v0.0.0-00010101000000-000000000000 => ../..
