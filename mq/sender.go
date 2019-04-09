package mq

//Sender is a mq service
type Sender interface {
	SendMessage(destination, message string) error
}
