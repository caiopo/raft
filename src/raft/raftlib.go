package raft

import (
	"strings"
	"udp"
)

type Message struct {
	Sender, Content, State string
}

func DecomposeMessage(msg string) *Message {

	msgTemp := strings.Split(msg, ":")

	return &Message{msgTemp[0], msgTemp[1], msgTemp[2]}

}

func (m *Message) ToString() string {

	sMsg := []string{m.Sender, m.Content, m.State}

	return strings.Join(sMsg, ":")

}

func (m *Message) SendTo(target string) {

	message := m.ToString()

	udp.Send(message, target)

}

// func DecomposeMessage(initMsg string) (msg *Message) {

// 	msg = new(Message)

// 	msgTemp := strings.Split(initMsg, ":")

// 	msg.Sender = msgTemp[0]
// 	msg.Content = msgTemp[1]

// 	if len(msgTemp) == 3 {
// 		msg.State = msgTemp[2]
// 	} else {
// 		msg.State = ""
// 	}

// 	return
// }
