package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go_mongoDb/internal/model"
	"go_mongoDb/internal/service"
)

type WebSocketServer struct {
	clients             map[*websocket.Conn]bool
	broadcast           chan []byte
	register            chan *websocket.Conn
	unregister          chan *websocket.Conn
	conversationService *service.ConversationService
	messageService      *service.MessageService
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebSocketServer(conversationService *service.ConversationService, messageService *service.MessageService) *WebSocketServer {
	return &WebSocketServer{
		clients:             make(map[*websocket.Conn]bool),
		broadcast:           make(chan []byte),
		register:            make(chan *websocket.Conn),
		unregister:          make(chan *websocket.Conn),
		conversationService: conversationService,
		messageService:      messageService,
	}
}

func (ws *WebSocketServer) Run() {
	http.HandleFunc("/ws/chat", ws.HandleConnections)
	go ws.handleMessages()
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func (ws *WebSocketServer) handleMessages() {
	for {
		select {
		case conn := <-ws.register:
			ws.clients[conn] = true
		case conn := <-ws.unregister:
			if _, ok := ws.clients[conn]; ok {
				delete(ws.clients, conn)
				conn.Close()
			}
		case message := <-ws.broadcast:
			for conn := range ws.clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error writing message: %v", err)
					conn.Close()
					delete(ws.clients, conn)
				}
			}
		}
	}
}

func (ws *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("Client connected: %s", conn.RemoteAddr())

	ws.register <- conn

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			ws.unregister <- conn
			break
		}
		ws.broadcast <- message

		var request map[string]interface{}
		if err := json.Unmarshal(message, &request); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		action, ok := request["action"].(string)
		if !ok {
			log.Println("Invalid action")
			continue
		}

		switch action {
		case "create_conversation":
			senderID, err1 := primitive.ObjectIDFromHex(request["senderId"].(string))
			receiverID, err2 := primitive.ObjectIDFromHex(request["receiverId"].(string))

			if err1 != nil || err2 != nil {
				log.Println("Invalid sender or receiver ID")
				continue
			}

			conversation := model.Conversation{
				SenderId:   senderID,
				ReceiverId: receiverID,
			}

			_, err := ws.conversationService.Create(conversation)
			if err != nil {
				log.Printf("Error creating conversation: %v", err)
				continue
			}

			log.Println("Conversation created successfully")

		case "get_conversationById":
			conversationId, err := primitive.ObjectIDFromHex(request["_id"].(string))
			if err != nil {
				log.Println("Invalid conversationID")
				continue
			}

			conversation, sender, receiver, err := ws.conversationService.GetConversationWithUsers(r.Context(), conversationId)
			if err != nil {
				log.Printf("Error getting conversation: %v", err)
				continue
			}

			response := map[string]interface{}{
				"conversation": conversation,
				"sender":       sender,
				"receiver":     receiver,
			}
			responseMessage, err := json.Marshal(response)
			if err != nil {
				log.Printf("Error marshalling response: %v", err)
				continue
			}

			err = conn.WriteMessage(websocket.TextMessage, responseMessage)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}

			log.Println("Conversation fetched successfully")

		case "send_message":
			conversationID, err1 := primitive.ObjectIDFromHex(request["conversationId"].(string))
			senderID, err2 := primitive.ObjectIDFromHex(request["senderId"].(string))

			if err1 != nil || err2 != nil {
				log.Println("Invalid conversation or sender ID")
				continue
			}

			message := model.Message{
				ConversationId: conversationID,
				SenderId:       senderID,
				Message:        request["message"].(string),
			}

			createdMessage, err := ws.messageService.Create(message)
			if err != nil {
				log.Printf("Error creating message: %v", err)
				continue
			}

			// Create the response
			response := map[string]interface{}{
				"status":  "success",
				"message": createdMessage,
			}
			responseMessage, err := json.Marshal(response)
			if err != nil {
				log.Printf("Error marshalling response: %v", err)
				continue
			}
		
			err = conn.WriteMessage(websocket.TextMessage, responseMessage)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
		
			log.Println("Message sent successfully")

		default:
			log.Println("Unknown action")
		}
	}
}
