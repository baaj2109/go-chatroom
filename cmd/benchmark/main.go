package main

import (
	"chatroom/model"
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	// "github.com/polaris1119/chatroom/logic"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var (
	userNum       int           // 使用者數
	loginInterval time.Duration // 使用者登入時間間隔
	msgInterval   time.Duration // 同一个使用者發送訊息間隔
)

func init() {
	flag.IntVar(&userNum, "u", 10, "登入使用者數")
	flag.DurationVar(&loginInterval, "l", 0, "使用者陸續登人時間間隔")
	flag.DurationVar(&msgInterval, "m", 20*time.Second, "使用者發送訊息時間間隔")
}

func main() {
	flag.Parse()

	for i := 0; i < userNum; i++ {
		go UserConnect("user" + strconv.Itoa(i))
		time.Sleep(loginInterval)
	}

	select {}
}

/*
1. 透過 dial 和服務端建立連接 即進入聊天室
注意函數第二個傳回的數值是 *http.response 因為這裡用不到這個傳回值 並且不能關閉 response body
所以直接用 _ 忽略
*/
func UserConnect(nickname string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws://localhost:2022/ws?nickname="+nickname, nil)
	log.Println("user:", nickname, "connect")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "内部錯誤！")
	/// 開心的 go routine 處理訊息發送
	go sendMessage(conn, nickname)

	ctx = context.Background()

	for {
		var message model.Message
		err = wsjson.Read(ctx, conn, &message)
		if err != nil {
			log.Println("receive msg error:", err)
			continue
		}

		if message.ClientSendTime.IsZero() {
			continue
		}
		if d := time.Now().Sub(message.ClientSendTime); d > 1*time.Second {
			fmt.Printf("接收到伺服器回應(%d):%#v\n", d.Milliseconds(), message)
		}
	}

	conn.Close(websocket.StatusNormalClosure, "")
}

func sendMessage(conn *websocket.Conn, nickname string) {
	ctx := context.Background()
	i := 1
	for {
		msg := map[string]string{
			"content":   "來自" + nickname + "的訊息:" + strconv.Itoa(i),
			"send_time": strconv.FormatInt(time.Now().UnixNano(), 10),
		}
		err := wsjson.Write(ctx, conn, msg)
		if err != nil {
			log.Println("send msg error:", err, "nickname:", nickname, "no:", i)
		}
		i++

		time.Sleep(msgInterval)
	}
}
