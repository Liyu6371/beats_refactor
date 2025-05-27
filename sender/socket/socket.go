package socket

import "sync"

const (
	Name = "socket"
)

type SocketInstance struct {
	instName string
	wg       sync.WaitGroup
}
