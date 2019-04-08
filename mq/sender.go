package mq

//Sender Sender
type Sender interface {
	SendMessage(destination, message string) error
}
