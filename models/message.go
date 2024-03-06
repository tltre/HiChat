package models

import (
	"HiChat/global"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
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

// the map of userId and MesNode
var clientMap = make(map[uint]*MsgNode, 0)

// a global channel to store the message sending to a Parse and HangOut UDP Server
var udpSendChan = make(chan []byte, 1024)

// a lock for binding user and msgNode
var lock sync.RWMutex

// open the UDP sending channel when initial message.go
func init() {
	go UdpSendData()
	go UdpRecData()
}

// Chat call by Service Layer to send message
func Chat(w http.ResponseWriter, r *http.Request) {
	// get User Id
	q := r.URL.Query()
	id := q.Get("userId")
	userId, err := strconv.Atoi(id)
	if err != nil {
		zap.S().Info("Failed to get userId: ", err)
		return
	}

	// upgrade socket
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
	msgNode := MsgNode{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// bind a MsgNode to User
	lock.Lock()
	clientMap[uint(userId)] = &msgNode
	lock.Unlock()

	go SendDataBySocket(&msgNode)
	go RecDataBySocket(&msgNode)
}

// SendDataBySocket get data from node and write in socket
func SendDataBySocket(node *MsgNode) {
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

// RecDataBySocket Receive message from user, and send it to a UDP Server for sending the message to target User
func RecDataBySocket(node *MsgNode) {
	for true {
		// get Message
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			zap.S().Info("Failed to get Message: ", err)
			return
		}

		// store in UDP Channel to send data to UDP Server
		udpSendChan <- data
	}
}

// UdpSendData send the message that user has posted to UDP Server
func UdpSendData() {
	// create a UDP connection to UDP Server(127.0.0.1:3000)
	udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 3000,
		Zone: "",
	})
	if err != nil {
		zap.S().Info("Failed to connect to UDP Server")
		return
	}

	defer udpConn.Close()
	// receive data and send it
	for true {
		select {
		case data := <-udpSendChan:
			_, err := udpConn.Write(data)
			if err != nil {
				zap.S().Info("Failed to send udp data")
				return
			}
		}
	}
}

// UdpRecData Receive data from udp client and HangOut the message to Target User
func UdpRecData() {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 3000,
		Zone: "",
	})
	if err != nil {
		zap.S().Info("Failed to Listen UDP Port")
		return
	}
	defer udpConn.Close()
	for true {
		var buf [1024]byte
		n, err := udpConn.Read(buf[0:])
		if err != nil {
			zap.S().Info("Failed to Read Data")
			return
		}
		dispatch(buf[0:n])
	}
}

// Parse the Data to Message And send to friend/group
func dispatch(data []byte) {
	// Parse to Message
	msg := Message{}
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
	message := Message{}
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
	usersId, err := FindMembersId(targetId)
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
