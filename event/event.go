package event

import (
	"github.com/nikunjgit/crypto/client"
	exchange "github.com/preichenberger/go-coinbase-exchange"

	"fmt"
	"sync"
)

type Generator struct {
	socket    *client.SocketClient
	listeners []Listener
	mutex     *sync.Mutex
}

func NewEventGenerator(socket *client.SocketClient) *Generator {
	return &Generator{socket, make([]Listener, 0, 10), &sync.Mutex{}}
}

type Listener interface {
	OnData(message *exchange.Message)
}

func (e *Generator) Register(listener Listener) {
	e.mutex.Lock()
	e.listeners = append(e.listeners, listener)
	e.mutex.Unlock()
}

func (e *Generator) Start() {
	go func() {
		for {
			message := &exchange.Message{}
			err := e.socket.Read(message)
			if err != nil {
				fmt.Println("Found error while read", err)
			} else {
				e.mutex.Lock()
				for _, l := range e.listeners {
					l.OnData(message)
				}
				e.mutex.Unlock()
			}
		}
	}()
}
