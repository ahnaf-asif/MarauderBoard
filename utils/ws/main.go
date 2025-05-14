package ws

import (
	"log"
	"sync"

	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/ahnafasif/MarauderBoard/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	UserID    int
	GroupID   int
	UserName  string
	Conn      *websocket.Conn
	SendQueue chan []byte
}

var (
	rooms = make(map[int]map[*Client]bool)
	mu    sync.Mutex
)

func JoinRoom(groupID int, client *Client) {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("Client %s joined group %d", client.UserName, groupID)
	if rooms[groupID] == nil {
		rooms[groupID] = make(map[*Client]bool)
	}
	rooms[groupID][client] = true
}

func LeaveRoom(groupID int, client *Client) {
	mu.Lock()
	defer mu.Unlock()

	if clients, ok := rooms[groupID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(rooms, groupID)
		}
	}
}

func BroadcastToRoom(groupID int, senderID int, template string, data fiber.Map) {
	mu.Lock()
	defer mu.Unlock()

	for client := range rooms[groupID] {
		receiver_user, _ := models.GetUserById(database.DB, uint(client.UserID))
		data["ReceiverUser"] = receiver_user
		html, _ := helpers.RenderPartial(template, data)
		log.Printf("%s\n", string(html))
		client.SendQueue <- html
	}
}
