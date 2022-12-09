package logic

import (
	"chatroom/model"
	"time"

	"github.com/spf13/cast"
)

func NewMessage(user *model.User, content string, clientTime string) *model.Message {
	message := &model.Message{
		User:    user,
		Type:    model.MessageTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}
	if clientTime != "" {
		message.ClientSendTime = time.Unix(0, cast.ToInt64(clientTime))
	}
	return message
}

func NewWelcomeMessage(user *model.User) *model.Message {
	return &model.Message{
		User:    user,
		Type:    model.MessageTypeWelcome,
		Content: user.NickName + " 進入聊天室",
		MsgTime: time.Now(),
	}
}

func NewUserEnterMessage(user *model.User) *model.Message {
	return &model.Message{
		User:    user,
		Type:    model.MessageTypeUserEnter,
		Content: user.NickName + " 加入了聊天室",
		MsgTime: time.Now(),
	}
}

func NewUserLeaveMessage(user *model.User) *model.Message {
	return &model.Message{
		User:    user,
		Type:    model.MessageTypeUserLeave,
		Content: user.NickName + " 離開了聊天室",
		MsgTime: time.Now(),
	}
}

func NewErrorMessage(content string) *model.Message {
	return &model.Message{
		User:    model.System,
		Type:    model.MessageTypeError,
		Content: content,
		MsgTime: time.Now(),
	}
}
