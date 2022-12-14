package server

import (
	"chatroom/logic"
	"chatroom/model"
	"log"
	"net/http"
	"regexp"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	/*
		accept 從用戶端接收 websocket驗證 並將連接升級到 websocket
		如果 origin 域與主機不同 accpet將拒絕驗證 除非設定了
		insecurseskipverify選項 透過第三個參數 acceptoptions設定
	*/
	conn, err := websocket.Accept(w, req, nil)
	if err != nil {
		log.Println("websocket accept error", err)
		return
	}

	// 1. new user join
	token := req.FormValue("token")
	nickname := req.FormValue("nickname")
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname length should between 2-20")
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("illegal nick length"))
		return
	}

	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("nickname already exist")
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage(" nickname already exists"))
		return
	}
	user := model.NewUser(conn, nickname, token, req.RemoteAddr)

	/// 開啟給使用者發送訊息的 go routine
	go user.ReceiveMessage(req.Context())

	///	給新使用者發送歡迎
	user.MessageChannel <- logic.NewWelcomeMessage(user)

	/// 通知所有使用者新加入者
	msg := logic.NewUserEnterMessage(user)
	logic.Broadcaster.Broadcast(msg)

	/// add user to user list
	logic.Broadcaster.UserEntering(user)
	log.Println("user:", nickname, " join")

	/// receive user message
	err = user.LoadMessage(req.Context())

	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error:", err)
		conn.Close(websocket.StatusTryAgainLater, "Read from client error")
	}

	textFromUser := user.SendMessage()
	sendMsg := logic.NewMessage(user, textFromUser["content"], textFromUser["send_time"])
	// sendMsg.Content = FilterSensitive(sendMsg.Content)

	// // 解析 content，看看 @ 谁了
	reg := regexp.MustCompile(`@[^\s@]{2,20}`)
	sendMsg.Ats = reg.FindAllString(sendMsg.Content, -1)
	logic.Broadcaster.Broadcast(sendMsg)

	// user leave
	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcaster.Broadcast(msg)
}
