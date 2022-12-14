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
	MessageTypeUserList
)

type Message struct {
	User           *User     `json:"user"`
	Type           int       `json:"type"`
	Content        string    `json:"content"`
	MsgTime        time.Time `json:"msg_time"`
	ClientSendTime time.Time `json:"client_send_time"`
	//message @ whom
	Ats []string `json:"ats"`
	// 私密訊息
	To    string           `json:"to"`
	Users map[string]*User `json:"users"`
}
