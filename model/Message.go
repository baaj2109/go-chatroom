package model

import (
	"time"
)

const (
	MessageTypeNormal = iota // normal user msg
	MessageTypeWelcome
	MessageTypeUserEnter
	MessageTypeUserLeave
	MessageTypeError
)

type Message struct {
	User           *User     `json:"user"`
	Type           int       `json:"type"`
	Content        string    `json:"content"`
	MsgTime        time.Time `json:"msg_time"`
	ClientSendTime time.Time `json:"client_send_time"`
	Ats            []string  `json:"ats"` //message @ whom
}
