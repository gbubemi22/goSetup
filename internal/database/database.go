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
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(MongoDb))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("Connected to MongoDB")
	return client, nil
}

func New() (*mongo.Database, error) {
	client, err := DBinstance()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize service: %w", err)
	}
	return client.Database("Gomongodb"), nil
}

type Service struct {
	db *mongo.Client
}

// OpenCollection opens a MongoDB collection
func (s *Service) OpenCollection(collectionName string) *mongo.Collection {
	collection := s.db.Database("Gomongodb").Collection(collectionName)
	return collection
}
