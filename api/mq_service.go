package api

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/wibl/webapp/mq"
)

//MqService is a rpc service for interaction with mq
type MqService struct {
	stompSender *mq.StompSender
}

//MqConnectArgs is args for Connect method
type MqConnectArgs struct {
	URL, User, Pass string
}

//MqSendArgs is args for Send method
type MqSendArgs struct {
	Queue, Message string
}

//MqReply is a reply for all MqService methods
type MqReply struct{}

//Connect connects with a specific mq server
func (h *MqService) Connect(r *http.Request, args *MqConnectArgs, reply *MqReply) error {
	url, err := url.Parse(args.URL)
	if err != nil {
		return err
	}
	stompSender, err := mq.NewStompSender(*url, args.User, args.Pass)
	if err != nil {
		return err
	}
	h.stompSender = stompSender
	return nil
}

//Send sends message to a specific queue
func (h *MqService) Send(r *http.Request, args *MqSendArgs, reply *MqReply) error {
	if h.stompSender == nil {
		return errors.New("You must connect to the mq service first")
	}
	h.stompSender.SendMessage(args.Queue, args.Message)
	return nil
}
