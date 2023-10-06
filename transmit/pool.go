package transmit

import (
	"github.com/TremblingV5/TinyCache/ds/consistenthash"
	"sync"
)

type IPeer interface {
	Get(key string) string
	Add(servers ...string)
}

type ServerPool struct {
	mu    sync.Mutex
	peers IPeer
}

func NewServerPool(serversNum int) *ServerPool {
	return &ServerPool{
		mu:    sync.Mutex{},
		peers: consistenthash.New(serversNum, nil),
	}
}

func (p *ServerPool) Add(servers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers.Add(servers...)
}

func (p *ServerPool) Pick(key string) string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.peers.Get(key)
}
