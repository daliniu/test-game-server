package selector

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/micro/go-micro/registry"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Random is a random strategy algorithm for node selection
func Random(services []*registry.Service) Next {
	var nodes []*registry.Node

	for _, service := range services {
		nodes = append(nodes, service.Nodes...)
	}

	return func() (*registry.Node, error) {
		if len(nodes) == 0 {
			return nil, errors.New(ErrNoneAvailable.Error())
		}

		i := rand.Int() % len(nodes)
		return nodes[i], nil
	}
}

// RoundRobin is a roundrobin strategy algorithm for node selection
func RoundRobin(services []*registry.Service) Next {
	var nodes []*registry.Node

	for _, service := range services {
		nodes = append(nodes, service.Nodes...)
	}

	var i = rand.Int()
	var mtx sync.Mutex

	return func() (*registry.Node, error) {
		if len(nodes) == 0 {
			return nil, errors.New(ErrNoneAvailable.Error())
		}

		mtx.Lock()
		node := nodes[i%len(nodes)]
		i++
		mtx.Unlock()

		return node, nil
	}
}