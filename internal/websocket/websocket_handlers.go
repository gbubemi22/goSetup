package websocket

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// 	"go_mongoDb/internal/service"
// 	_"go_mongoDb/internal/websocket"
// )

// type WebSocketServer struct {
// 	clients             map[*websocket.Conn]bool
// 	broadcast           chan []byte
// 	register            chan *websocket.Conn
// 	unregister          chan *websocket.Conn
// 	conversationService *service.ConversationService
// 	conversationHandler *conversation_handler.ConversationHandler
// }

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func NewWebSocketServer(conversationService *service.ConversationService) *WebSocketServer {
// 	return &WebSocketServer{
// 		clients:             make(map[*websocket.Conn]bool),
// 		broadcast:           make(chan []byte),
// 		register:            make(chan *websocket.Conn),
// 		unregister:          make(chan *websocket.Conn),
// 		conversationService: conversationService,
// 		conversationHandler: conversation_handler.NewConversationHandler(ws),

// 	}
// }

// func (ws *WebSocketServer) Run() {
// 	http.HandleFunc("/ws/chat", ws.HandleConnections)
// 	go ws.handleMessages()
// 	log.Fatal(http.ListenAndServe(":8081", nil))
// }

// func (ws *WebSocketServer) handleMessages() {
// 	for {
// 		select {
// 		case conn := <-ws.register:
// 			ws.clients[conn] = true
// 		case conn := <-ws.unregister:
// 			if _, ok := ws.clients[conn]; ok {
// 				delete(ws.clients, conn)
// 				conn.Close()
// 			}
// 		case message := <-ws.broadcast:
// 			for conn := range ws.clients {
// 				err := conn.WriteMessage(websocket.TextMessage, message)
// 				if err != nil {
// 					log.Printf("Error writing message: %v", err)
// 					conn.Close()
// 					delete(ws.clients, conn)
// 				}
// 			}
// 		}
// 	}
// }

//  func (ws *WebSocketServer) Run() {
// 	conversationHandler := NewConversationEvent(ws)
// 	http.HandleFunc("/ws/chat", conversationHandler.HandleConnections)

// 	messageHandler := NewMessageHandler(ws)
// 	go messageHandler.HandleMessages()

// 	log.Fatal(http.ListenAndServe(":8081", nil))
// }

// type WebSocketServer struct {
// 	clients             map[*websocket.Conn]bool
// 	broadcast           chan []byte
// 	register            chan *websocket.Conn
// 	unregister          chan *websocket.Conn
// 	conversationService *service.ConversationService
// }

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func NewWebSocketServer(conversationService *service.ConversationService) *WebSocketServer {
// 	return &WebSocketServer{
// 		clients:             make(map[*websocket.Conn]bool),
// 		broadcast:           make(chan []byte),
// 		register:            make(chan *websocket.Conn),
// 		unregister:          make(chan *websocket.Conn),
// 		conversationService: conversationService,
// 	}
// }

// func (ws *WebSocketServer) Run() {
// 	http.HandleFunc("/ws/chat", ws.HandleConnections)
// 	go ws.handleMessages()
// 	log.Fatal(http.ListenAndServe(":8081", nil))
// }

// func (ws *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Printf("Error upgrading to websocket: %v", err)
// 		return
// 	}
// 	defer conn.Close()

// 	log.Printf("Client connected: %s", conn.RemoteAddr())

// 	ws.register <- conn

// 	for {
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Printf("Error reading message: %v", err)
// 			ws.unregister <- conn
// 			break
// 		}
// 		ws.broadcast <- message

// 		var request map[string]interface{}
// 		if err := json.Unmarshal(message, &request); err != nil {
// 			log.Printf("Error parsing message: %v", err)
// 			continue
// 		}

// 		action, ok := request["action"].(string)
// 		if !ok {
// 			log.Println("Invalid action")
// 			continue
// 		}
// 		switch action {
// 		case "create_conversation":
// 			senderId, err1 := primitive.ObjectIDFromHex(request["senderId"].(string))
// 			receiverId, err2 := primitive.ObjectIDFromHex(request["receiverId"].(string))

// 			if err1 != nil || err2 != nil {
// 				log.Println("Invalid sender or receiver ID")
// 				continue
// 			}

// 			conversation := model.Conversation{
// 				SenderId:   senderId,
// 				ReceiverId: receiverId,
// 			}

// 			_, err := ws.conversationService.Create(conversation)
// 			if err != nil {
// 				log.Printf("Error creating conversation: %v", err)
// 				continue
// 			}

// 			log.Println("Conversation created successfully")
// 		default:
// 			log.Println("Unknown action")
// 		}
// 	}
// }

// func (ws *WebSocketServer) handleMessages() {
// 	for {
// 		select {
// 		case conn := <-ws.register:
// 			ws.clients[conn] = true
// 		case conn := <-ws.unregister:
// 			if _, ok := ws.clients[conn]; ok {
// 				delete(ws.clients, conn)
// 				conn.Close()
// 			}
// 		case message := <-ws.broadcast:
// 			for conn := range ws.clients {
// 				err := conn.WriteMessage(websocket.TextMessage, message)
// 				if err != nil {
// 					log.Printf("Error writing message: %v", err)
// 					conn.Close()
// 					delete(ws.clients, conn)
// 				}
// 			}
// 		}
// 	}
// }
