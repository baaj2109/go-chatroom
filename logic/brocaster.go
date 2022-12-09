package logic

import (
	"chatroom/config"
	"chatroom/model"
	"log"
)

type broadcaster struct {
	users           map[string]*model.User
	enteringChannel chan *model.User
	leavingChannel  chan *model.User
	messageChannel  chan *model.Message

	/// for user can join
	checkUserChannel        chan string
	checkUserCanJoinChannel chan bool

	/// for user list
	requestUsersChannel chan struct{}
	usersChannel        chan []*model.User
}

// / singleton pattern
var Broadcaster = &broadcaster{
	users: make(map[string]*model.User),

	enteringChannel: make(chan *model.User),
	leavingChannel:  make(chan *model.User),
	messageChannel:  make(chan *model.Message, config.MessageQueueLen),

	checkUserChannel:        make(chan string),
	checkUserCanJoinChannel: make(chan bool),

	requestUsersChannel: make(chan struct{}),
	// usersChannel:        make(chan []*model.User),
}

func (b *broadcaster) GetUserList() []*model.User {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}

func (b *broadcaster) CanEnterRoom(name string) bool {
	b.checkUserChannel <- name
	return <-b.checkUserCanJoinChannel
}

func (b *broadcaster) Broadcast(message *model.Message) {
	if len(b.messageChannel) >= config.MessageQueueLen {
		log.Print("broadcast queeu full")
	}
}

func (b *broadcaster) UserEntering(u *model.User) {
	b.enteringChannel <- u
}

func (b *broadcaster) UserLeaving(u *model.User) {
	b.leavingChannel <- u
}

func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enteringChannel:
			//newby
			b.users[user.NickName] = user
			// OfflineProcessor.Send(user)

		case user := <-b.leavingChannel:
			delete(b.users, user.NickName)
			user.CloseMessageChannel()

		case <-b.requestUsersChannel:
			userList := make([]*model.User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}

			b.usersChannel <- userList

		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanJoinChannel <- false
			} else {
				b.checkUserCanJoinChannel <- true
			}

		case msg := <-b.messageChannel:
			for _, user := range b.users {
				if user.ID == msg.User.ID {
					continue
				}
				user.MessageChannel <- msg
			}
			// OfflineProcessor.Save(msg)
		}
	}
}
