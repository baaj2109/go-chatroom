package main

// import (
// 	"bufio"
// 	"chatroom/model"
// 	"fmt"
// 	"log"
// 	"net"
// 	"strconv"
// 	"sync"
// 	"time"
// )

// var (
// 	// 新用戶到来，通过该 channel 進行登记
// 	enteringChannel = make(chan *model.User)
// 	// 用戶离開，通过该 channel 進行登记
// 	leavingChannel = make(chan *model.User)
// 	// 广播专用的用戶普通消息 channel，缓冲是尽可能避免出现异常情况堵塞
// 	messageChannel = make(chan model.Message, 8)
// )

// func main() {

// 	listener, err := net.Listen("tcp", ":2020")
// 	if err != nil {
// 		panic(err)
// 	}
// 	go broadcaster()
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 		go handleConn(conn)
// 	}
// }

// // broadcaster 用于记录聊天室用戶，并進行消息广播：
// // 1. 新用戶進来；2. 用戶普通消息；3. 用戶离開
// func broadcaster() {
// 	users := make(map[*model.User]struct{})

// 	for {
// 		select {
// 		case user := <-enteringChannel:
// 			// 新用戶進入
// 			users[user] = struct{}{}
// 		case user := <-leavingChannel:
// 			// 用戶离開
// 			delete(users, user)
// 			// 避免 goroutine 泄露
// 			close(user.MessageChannel)
// 		case msg := <-messageChannel:
// 			// 给所有在线用戶發送消息
// 			for user := range users {
// 				if user.ID == msg.OwnerID {
// 					continue
// 				}
// 				user.MessageChannel <- msg.Content
// 			}
// 		}
// 	}
// }

// func handleConn(conn net.Conn) {
// 	defer conn.Close()
// 	// 1. 新用戶進来，构建该用戶的实例

// 	user := &model.User{
// 		ID:             GetUserID(),
// 		Addr:           conn.RemoteAddr().String(),
// 		EnterAt:        time.Now(),
// 		MessageChannel: make(chan string, 8),
// 	}
// 	// 2. 当前在一個新的 goroutine 中，用来進行读操作，因此需要開一個 goroutine 用于写操作
// 	// 读写 goroutine 之间可以通过 channel 進行通信
// 	go sendMessage(conn, user.MessageChannel)

// 	// 3. 给当前用戶發送歡迎信息；给所有用戶告知新用戶到来
// 	user.MessageChannel <- "Welcome, " + user.String()
// 	msg := model.Message{
// 		OwnerID: user.ID,
// 		Content: "user:`" + strconv.Itoa(user.ID) + "` has enter",
// 	}
// 	messageChannel <- msg

// 	// 4. 将该记录到全局的用戶列表中，避免用锁
// 	enteringChannel <- user

// 	// 控制超时用戶踢出
// 	var userActive = make(chan struct{})
// 	go func() {
// 		d := 1 * time.Minute
// 		timer := time.NewTimer(d)
// 		for {
// 			select {
// 			case <-timer.C:
// 				conn.Close()
// 			case <-userActive:
// 				timer.Reset(d)
// 			}
// 		}
// 	}()

// 	// 5. 循环读取用戶的输入
// 	input := bufio.NewScanner(conn)
// 	for input.Scan() {
// 		msg.Content = strconv.Itoa(user.ID) + ":" + input.Text()
// 		messageChannel <- msg

// 		// 用戶活跃
// 		userActive <- struct{}{}
// 	}

// 	if err := input.Err(); err != nil {
// 		log.Println("读取錯誤：", err)
// 	}

// 	// 6. 用戶离開
// 	leavingChannel <- user
// 	msg.Content = "user:`" + strconv.Itoa(user.ID) + "` has left"
// 	messageChannel <- msg
// }

// func sendMessage(conn net.Conn, ch <-chan string) {
// 	for msg := range ch {
// 		fmt.Fprintln(conn, msg)
// 	}
// }

// // 生成用戶 ID
// var (
// 	globalID int
// 	idLocker sync.Mutex
// )

// func GetUserID() int {
// 	idLocker.Lock()
// 	defer idLocker.Unlock()

// 	globalID++
// 	return globalID
// }
