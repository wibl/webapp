package mq

import (
	"net/url"

	"github.com/go-stomp/stomp"
)

//StompSender is a mq.Sender implementation
type StompSender struct {
	*stomp.Conn
}

//SendMessage sends a message to the STOMP server
func (c *StompSender) SendMessage(destination, message string) error {
	return c.Send(destination, "text/plain", []byte(message), stomp.SendOpt.Header("A", "B"))
}

//New creates a new StompSender and a TCP connection to a STOMP server
func NewStompSender(addr url.URL, user, pass string) (*StompSender, error) {
	stompConn, err := stomp.Dial(addr.Scheme, addr.Host, stomp.ConnOpt.Login(user, pass))
	if err != nil {
		return nil, err
	}
	return &StompSender{stompConn}, nil
}
