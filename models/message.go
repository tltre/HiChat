package models

import (
	"HiChat/global"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

// Message define the structure of message
/*
the params are:
	* FromId: message sender id
	* TargetId: message receiver id
	* Type: type of chat, including chatting in group or to user
	* Media: type of message media, including text and file(such as picture and voice data)
	* Content: content of text message
	* Url: the url of file
	* Desc: description of file
*/
type Message struct {
	gorm.Model
	FromId   uint `json:"userId"`
	TargetId uint `json:"targetId"`
	Type     int
	Media    int
	Content  string
	Url      string `json:"url"`
	Desc     string
}

// MarshalBinary marshal Message to []byte
func (msg Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(msg)
}

// MsgNode is a node bind to a specific User to send and receive Message
/*
the params are:
	* Conn: a connection of websocket
	* Addr: address of user
	* DataQueue: message queue
	* GroupSets: indicate group or friend
*/
type MsgNode struct {
	Conn      *websocket.Conn
	Addr      string
	DataQueue chan []byte
	GroupSets set.Interface
}

// GetMsgFromRedis Get Records From Redis
func GetMsgFromRedis(idA uint, idB uint, start int64, end int64, isRcv bool) []string {
	// get Key
	var key string
	if idA < idB {
		key = fmt.Sprintf("msg_%d_%d", idA, idB)
	} else {
		key = fmt.Sprintf("msg_%d_%d", idB, idA)
	}

	ctx := context.Background()
	var result []string
	var err error
	if isRcv {
		// Get Record from Near to Far
		result, err = global.RedisDB.ZRevRange(ctx, key, start, end).Result()
	} else {
		// Get Record from Far to Near
		result, err = global.RedisDB.ZRange(ctx, key, start, end).Result()
	}
	if err != nil {
		zap.S().Info("Failed to get records")
		return nil
	}
	zap.S().Info("Success to get records")
	return result
}
