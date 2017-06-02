package handlers

import "net"

type Limiter interface {
	Check() bool
	Allowed() bool
	Done()
}

type Handlers struct {
	PrivateNets []*net.IPNet
}

func NewHandlers() *Handlers {
	return &Handlers{}
}
