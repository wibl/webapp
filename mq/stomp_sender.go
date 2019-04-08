package mq

import (
	"github.com/go-stomp/stomp"
)

//StompSender StompSender
type StompSender struct {
	*stomp.Conn
}

//SendMessage SendMessage
func (c *StompSender) SendMessage(destination, message string) error {
	return c.Send(destination, "text/plain", []byte(message))
}

//Dial Dial
func Dial(network, addr string) (*StompSender, error) {
	stompConn, err := stomp.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &StompSender{stompConn}, nil
}
