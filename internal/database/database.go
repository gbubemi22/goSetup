package database

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var globalClient *mongo.Client

func init() {
    client, err := DBinstance()
    if err != nil {
        log.Fatalf("Failed to initialize MongoDB client: %v", err)
    }
    globalClient = client
}

func DBinstance() (*mongo.Client, error) {
    MongoDb := os.Getenv("DB_DATABASE")
    client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
    if err != nil {
        return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
    }

    fmt.Println("Connected to MongoDB")
    return client, nil
}

type Service struct {
    db *mongo.Client
}

func New() (*Service, error) {
    client, err := DBinstance()
    if err != nil {
        return nil, fmt.Errorf("failed to initialize service: %w", err)
    }
    return &Service{db: client}, nil
}

func (s *Service) Health() map[string]string {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    err := s.db.Ping(ctx, nil)
    if err != nil {
        log.Fatalf("db down: %v", err)
    }

    return map[string]string{
        "message": "It's healthy",
    }
}

// OpenCollection opens a MongoDB collection
func (s *Service) OpenCollection(collectionName string) *mongo.Collection {
    collection := s.db.Database("GO_Rental").Collection(collectionName)
    return collection
}



// package database

// import (
// 	"context"
// 	"fmt"
// 	//"log"
// 	"os"
// 	"time"

// 	_ "github.com/joho/godotenv/autoload"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// type Service interface {
// 	Health() map[string]string
// }

// type service struct {
// 	db *mongo.Client
// }



// func DBinstance() (*mongo.Client, error) {

// 	MongoDb := os.Getenv("DB_DATABASE")
// 	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	err = client.Connect(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
// 	}

// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
// 	}

// 	fmt.Println("Connected to MongoDB")
// 	return client, nil
// }

// // OpenCollection opens a MongoDB collection
// func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
// 	collection := client.Database("GO_Rental").Collection(collectionName)
// 	return collection
// }



