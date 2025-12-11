package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/obfio/tmx-solver-golang/mongo"
	allsites "github.com/obfio/tmx-solver-golang/sites/allSites"
)

type Socket struct {
	APIKey    string
	WriteChan chan string
	Socket    *websocket.Conn
}

var (
	upgrader = websocket.Upgrader{}
	Sockets  = make(map[string]*Socket)
	mutex    = sync.RWMutex{}

	taskChan = make(chan *allsites.ProxyRequest)
)

func NewSocket(Conn *websocket.Conn) *Socket {
	return &Socket{
		Socket:    Conn,
		WriteChan: make(chan string),
	}
}

// possible actions are
// 1.) auth
// 2.) solve
type messageRecv struct {
	Action string                `json:"action"`
	Data   allsites.ProxyRequest `json:"data"`
}

type cookieResponse struct {
	APIKey    string `json:"-"`
	Proxy     string `json:"proxy,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
	Cookie    string `json:"cookie,omitempty"`
	SessionID string `json:"sessionID,omitempty"`
	Error     string `json:"error,omitempty"`
}

func HandleNewConnection(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading connection: %v\n", err)
		return
	}

	socket := NewSocket(conn)
	for {
		messageType, message, err := socket.Socket.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			return
		}
		if !(messageType == websocket.TextMessage) {
			continue
		}
		var recv messageRecv
		err = json.Unmarshal(message, &recv)
		if err != nil {
			fmt.Printf("Error unmarshalling message: %v\n", err)
			return
		}
		var response cookieResponse
		if recv.Action != "auth" {
			response.Error = "first action must be auth"
			f, _ := json.Marshal(response)
			socket.Socket.WriteMessage(websocket.TextMessage, f)
			continue
		}

		// validate API key
		if !mongo.ValidateAPIKey(recv.Data.APIKey) {
			response.Error = "invalid API key"
			f, _ := json.Marshal(response)
			socket.Socket.WriteMessage(websocket.TextMessage, f)
			continue
		}
		// add socket to map
		mutex.Lock()
		socket.APIKey = recv.Data.APIKey
		Sockets[recv.Data.APIKey] = socket
		mutex.Unlock()
		go socket.handleReads()
		go socket.handleWrites()
		socket.WriteChan <- "authorized"
		break
	}
}

func (s *Socket) handleReads() {
	for {
		messageType, message, err := s.Socket.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			return
		}
		if !(messageType == websocket.TextMessage) {
			continue
		}
		var recv messageRecv
		err = json.Unmarshal(message, &recv)
		if err != nil {
			fmt.Printf("Error unmarshalling message: %v\n", err)
			return
		}
		var response cookieResponse
		if recv.Action != "solve" {
			response.Error = "invalid action"
			f, _ := json.Marshal(response)
			s.Socket.WriteMessage(websocket.TextMessage, f)
			continue
		}
		// fmt.Printf("solve: %+v\n", recv.Data)
		recv.Data.APIKey = s.APIKey
		taskChan <- &recv.Data
	}
}

func (s *Socket) handleWrites() {
	for {
		msg := <-s.WriteChan
		s.Socket.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

var (
	locker = sync.RWMutex{}
)

func HandleTasks() {
	for {
		task := <-taskChan
		go func(t *allsites.ProxyRequest) {
			locker.Lock()
			socket := Sockets[t.APIKey]
			locker.Unlock()
			if socket == nil || socket.Socket == nil {
				return
			}
			if t.Mobile {
				t.Site += "MOBILE"
			}
			response := allsites.GetCookies(t)
			if response.Error == "" {
				mongo.UpdateUsesCount(t.APIKey)
			}

			f, _ := json.Marshal(response)
			socket.WriteChan <- string(f)
		}(task)
	}
}
