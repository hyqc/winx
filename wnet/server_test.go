package wnet

import (
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer("winx v0.1")
	go MockClient("tcp", "127.0.0.1:8888")
	s.Serve()
}
