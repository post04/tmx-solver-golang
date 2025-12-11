package config

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DbCtx = context.Background()
	DbCnx = DatabaseConnection().Database("tmxapi")
)

type config struct {
	MongoURL string `json:"mongodb"`
}

func DatabaseConnection() *mongo.Client {
	f, err := os.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	var c config
	err = json.Unmarshal(f, &c)
	if err != nil {
		panic(err)
	}
	clientOptions := options.Client().ApplyURI(c.MongoURL)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
	return client
}
