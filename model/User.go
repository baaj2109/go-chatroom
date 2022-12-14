package model

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var System = &User{}
var globalUID uint32 = 0

type User struct {
	ID             int
	NickName       string
	Addr           string
	EnterAt        time.Time
	MessageChannel chan *Message
	Token          string
	conn           *websocket.Conn
	IsNew          bool
	message        map[string]string
}

func (u *User) String() string {
	return u.Addr + ", UID:" + strconv.Itoa(u.ID) + ", Enter At:" +
		u.EnterAt.Format("2006-01-02 15:04:05+8000")
}

func (u *User) CloseMessageChannel() {
	close(u.MessageChannel)
}

func (u *User) ReceiveMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		wsjson.Write(ctx, u.conn, msg)
	}
}

func (u *User) SendMessage() map[string]string {
	return u.message
}

func (u *User) LoadMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]string
		err        error
	)
	for {
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		if err != nil {
			/// 判斷連接是否關閉，正常關閉，不認為是錯誤
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			} else if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		u.message = receiveMsg

		// 内容发送到聊天室
		// sendMsg := NewMessage(u, receiveMsg["content"], receiveMsg["send_time"])
		// // sendMsg.Content = FilterSensitive(sendMsg.Content)

		// // 解析 content，看看 @ 谁了
		// reg := regexp.MustCompile(`@[^\s@]{2,20}`)
		// sendMsg.Ats = reg.FindAllString(sendMsg.Content, -1)

		// logic.Broadcaster.Broadcast(sendMsg)
	}
}

func NewUser(conn *websocket.Conn, token, nickname, addr string) *User {
	// func NewUser(conn *websocket.Conn, nickname string, addr string) *User {
	user := &User{
		NickName:       nickname,
		Addr:           addr,
		EnterAt:        time.Now(),
		MessageChannel: make(chan *Message, 32),
		Token:          token,

		conn: conn,
	}

	if user.Token != "" {
		uid, err := parseTokenAndValidate(token, nickname)
		if err == nil {
			user.ID = uid
		}
	}

	if user.ID == 0 {
		user.ID = int(atomic.AddUint32(&globalUID, 1))
		user.Token = genToken(user.ID, user.NickName)
		user.IsNew = true
	}

	return user
}
func genToken(uid int, nickname string) string {
	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid)

	messageMAC := macSha256([]byte(message), []byte(secret))

	return fmt.Sprintf("%suid%d", base64.StdEncoding.EncodeToString(messageMAC), uid)
}

// / HMAC - SHA256 加密
func macSha256(message, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}
func validateMAC(message, messageMAC, secret []byte) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}
func parseTokenAndValidate(token, nickname string) (int, error) {
	pos := strings.LastIndex(token, "uid")
	messageMAC, err := base64.StdEncoding.DecodeString(token[:pos])
	if err != nil {
		return 0, err
	}
	uid := cast.ToInt(token[pos+3:])

	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid)

	ok := validateMAC([]byte(message), messageMAC, []byte(secret))
	if ok {
		return uid, nil
	}
	return 0, errors.New("token is illegal")
}
