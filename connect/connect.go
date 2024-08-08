package connect

import (
	"HiChat/global"
	"HiChat/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gopkg.in/fatih/set.v0"
	"net/http"
	"sync"
)

// mapping of userId and MesNode
var clientMap = make(map[uint]*models.MsgNode, 0)

// a lock for binding user and msgNode
var lock sync.RWMutex

func UpgradeConnection(w http.ResponseWriter, r *http.Request, userId int) {
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		zap.S().Info("Failed to upgrade socket: ", err)
		return
	}

	// new a MsgNode
	msgNode := models.MsgNode{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// bind a MsgNode to User
	lock.Lock()
	clientMap[uint(userId)] = &msgNode
	lock.Unlock()

	go SendDataBySocket(&msgNode)
	go RecvDataBySocket(&msgNode)
}

// SendDataBySocket get data from node and write in socket
func SendDataBySocket(node *models.MsgNode) {
	for true {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				zap.S().Info("Failed to write message in websocket")
				return
			}
			fmt.Println("Success to send message in websocket")
		}
	}
}

// RecvDataBySocket Receive message from user, and send it to a UDP Server for sending the message to target User
func RecvDataBySocket(node *models.MsgNode) {
	for true {
		// get Message
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			zap.S().Info("Failed to get Message: ", err)
			return
		}

		dispatch(data)
	}
}

// Parse the Data to Message And send to friend/group
func dispatch(data []byte) {
	// Parse to Message
	msg := models.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		zap.S().Info("Failed to Parse to Message")
		return
	}

	// Send Message
	switch msg.Type {
	case 1:
		// send message to friend
		SendMessageToFriendAndSave(msg.TargetId, data)
	case 2:
		// send message to group
		SendMessageToCommunity(msg.FromId, msg.TargetId, data)
	}

}

// SendMessageToFriendAndSave send message to friend
func SendMessageToFriendAndSave(id uint, msg []byte) {
	lock.Lock()
	node, ok := clientMap[id]
	lock.Unlock()
	if !ok {
		zap.S().Info("Failed to Get Target User Node")
		return
	}

	// send message by socket
	zap.S().Info("Target Id: ", id, "Node: ", node)
	node.DataQueue <- msg

	// Parse to Message
	message := models.Message{}
	err := json.Unmarshal(msg, &message)
	if err != nil {
		zap.S().Info("Failed to Parse data to Message")
		return
	}

	// generate key
	var key string
	if message.FromId < message.TargetId {
		key = fmt.Sprintf("msg_%d_%d", message.FromId, message.TargetId)
	} else {
		key = fmt.Sprintf("msg_%d_%d", message.TargetId, message.FromId)
	}

	// get the number of record
	ctx := context.Background()
	res, err := global.RedisDB.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		zap.S().Info("Failed to create record")
		return
	}

	// store the Message
	score := float64(cap(res)) + 1
	_, err = global.RedisDB.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: msg,
	}).Result()
	if err != nil {
		zap.S().Info("Failed to store Message")
		return
	}
	zap.S().Info("Success to Save Message")
}

// SendMessageToCommunity find all user in the group and send to them
func SendMessageToCommunity(fromId, targetId uint, msg []byte) {
	usersId, err := models.FindMembersId(targetId)
	if err != nil {
		zap.S().Info("Failed to Get Members Id")
		return
	}
	for _, userId := range *usersId {
		if userId != fromId {
			SendMessageToFriendAndSave(userId, msg)
		}
	}
}
