package mq

import (
	"github.com/go-stomp/stomp"
)

//StompSender is a mq.Sender implementation
type StompSender struct {
	*stomp.Conn
}

//SendMessage sends a message to the STOMP server
func (c *StompSender) SendMessage(destination, message string) error {
	return c.Send(destination, "text/plain", []byte(message))
}

//New creates a new StompSender and a connection to a STOMP server
func New(addr string) (*StompSender, error) {
	stompConn, err := stomp.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &StompSender{stompConn}, nil
}
