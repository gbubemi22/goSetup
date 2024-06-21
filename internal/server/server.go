package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go_mongoDb/internal/database"
)

type Server struct {
	port int
	db   *database.Service
	
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
	  
	}
 
	NewServer := &Server{
	    port: port,
	    db:   db,
	}
 
	// Declare Server config
	server := &http.Server{
	    Addr:         fmt.Sprintf(":%d", NewServer.port),
	    Handler:      NewServer.RegisterRoutes(),
	    IdleTimeout:  time.Minute,
	    ReadTimeout:  10 * time.Second,
	    WriteTimeout: 30 * time.Second,
	}
 
	return server
 }
 
