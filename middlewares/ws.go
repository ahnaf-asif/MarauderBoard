package middlewares

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/models"
	"github.com/ahnafasif/MarauderBoard/utils/ws"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type ChatMessage struct {
	Content     string `json:"content"`
	UserID      string `json:"user_id"`
	ChatGroupID string `json:"chat_group_id"`
	// Add a catch-all if needed:
	// HEADERS map[string]interface{} `json:"HEADERS"`
}

func ChatSocketHandler(c *fiber.Ctx) error {
	user_id, _ := strconv.Atoi(c.Params("user_id"))
	group_id, _ := strconv.Atoi(c.Params("group_id"))
	return websocket.New(func(conn *websocket.Conn) {
		user, _ := models.GetUserById(database.DB, uint(user_id))
		groupID := group_id

		client := &ws.Client{
			UserID:    int(user.ID),
			GroupID:   groupID,
			UserName:  user.FirstName + " " + user.LastName,
			Conn:      conn,
			SendQueue: make(chan []byte, 256),
		}

		ws.JoinRoom(groupID, client)
		// ws.BroadcastToRoom(groupID, []byte("User "+client.UserName+" has joined the chat."))
		defer func() {
			ws.LeaveRoom(groupID, client)
			conn.Close()
		}()

		go func() {
			for msg := range client.SendQueue {
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					break
				}
			}
		}()

		for {
			_, msg, err := conn.ReadMessage()

			var chatMessage ChatMessage
			json.Unmarshal(msg, &chatMessage)

			if err != nil {
				break
			}

			user_id_int, _ := strconv.Atoi(chatMessage.UserID)

			newMsg := models.ChatMessage{
				Content:     chatMessage.Content,
				UserId:      user_id_int,
				ChatGroupId: groupID,
			}

			if err := database.DB.Create(&newMsg).Error; err != nil {
				continue
			}
			// preload User of newMsg
			if err := database.DB.Preload("User").First(&newMsg).Error; err != nil {
				continue
			}

			sender_user, _ := models.GetUserById(database.DB, uint(user_id_int))

			data := fiber.Map{
				"Message":    newMsg,
				"SenderUser": sender_user,
			}

			log.Printf("Broadcasting message to group %d", groupID)

			ws.BroadcastToRoom(groupID,
				int(sender_user.ID),
				"partials/chat-message",
				data,
			)
		}
	})(c)
}
