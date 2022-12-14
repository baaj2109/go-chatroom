package logic

import (
	"chatroom/model"
	"container/ring"

	"github.com/spf13/viper"
)

type offlineProcessor struct {
	n int

	// 保存所有用戶最近的 n 條訊息
	recentRing *ring.Ring

	// 保存某個用戶離線訊息（一樣 n 條）
	userRing map[string]*ring.Ring
}

var OfflineProcessor = newOfflineProcessor()

func newOfflineProcessor() *offlineProcessor {
	n := viper.GetInt("offline-num")

	return &offlineProcessor{
		n:          n,
		recentRing: ring.New(n),
		userRing:   make(map[string]*ring.Ring),
	}
}


/*
 將使用者訊息存入 recentRing中 並後移一個位置
 判斷訊息中是某有＠誰 並單獨把這種訊息儲存為一個訊息列表
*/
func (o *offlineProcessor) Save(msg *model.Message) {
	if msg.Type != model.MessageTypeNormal {
		return
	}
	o.recentRing.Value = msg
	o.recentRing = o.recentRing.Next()

	for _, nickname := range msg.Ats {
		nickname = nickname[1:]
		var (
			r  *ring.Ring
			ok bool
		)
		if r, ok = o.userRing[nickname]; !ok {
			r = ring.New(o.n)
		}
		r.Value = msg
		o.userRing[nickname] = r.Next()
	}
}

func (o *offlineProcessor) Send(user *model.User) {
	o.recentRing.Do(func(value interface{}) {
		if value != nil {
			user.MessageChannel <- value.(*model.Message)
		}
	})

	if user.IsNew {
		return
	}

	if r, ok := o.userRing[user.NickName]; ok {
		r.Do(func(value interface{}) {
			if value != nil {
				user.MessageChannel <- value.(*model.Message)
			}
		})

		delete(o.userRing, user.NickName)
	}
}
