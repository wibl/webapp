package main

import (
	"testing"
)

type fakeMess struct {
	destination, message string
}

type fakeSender struct {
	messages []*fakeMess
}

func (fs *fakeSender) SendMessage(destination, message string) error {
	fs.messages = append(fs.messages, &fakeMess{destination, message})
	return nil
}

func TestSendTestMessage(t *testing.T) {

	fakeSender := &fakeSender{}

	context := &appContext{
		sender: fakeSender,
	}

	sendTestMessage(context)

	got := fakeSender.messages
	if len(got) != 1 {
		t.Errorf("Count of sended messages; want 1 got %d", len(got))
	}
	if got[0].destination != "/queue/test-1" {
		t.Errorf("Message destination; want /queue/test-1 got %s", got[0].destination)
	}
	if got[0].message != "TEST" {
		t.Errorf("Message body; want TEST got %s", got[0].message)
	}
}
