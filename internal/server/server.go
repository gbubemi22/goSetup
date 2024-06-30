package server

import (
    "fmt"
    "net/http"
    "os"
    "strconv"
    "time"

    _ "github.com/joho/godotenv/autoload"
    "go_mongoDb/internal/database"
    "go_mongoDb/internal/service"
    "go_mongoDb/internal/websocket"
    "go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
    port int
    db   *mongo.Database
    ws   *websocket.WebSocketServer
}

func NewServer() *http.Server {
    portStr := os.Getenv("PORT")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        port = 8080
    }

    db, err := database.New()
    if err != nil {
        fmt.Printf("Error initializing database: %v\n", err)
        os.Exit(1)
    }

    conversationService := service.NewConversationService(db)
    ws := websocket.NewWebSocketServer(conversationService)
    go ws.Run()

    newServer := &Server{
        port: port,
        db:   db,
        ws:   ws,
    }

    server := &http.Server{
        Addr:         fmt.Sprintf(":%d", newServer.port),
        Handler:      newServer.RegisterRoutes(),
        IdleTimeout:  time.Minute,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 30 * time.Second,
    }

    return server
}
