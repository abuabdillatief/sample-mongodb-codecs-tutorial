package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func ConnectMongo() *mongo.Client {
	clientOpts := options.Client().
		SetDirect(true).
		SetConnectTimeout(time.Second * 60).
		SetMaxPoolSize(uint64(100)).
		SetMinPoolSize(10).
		ApplyURI("mongodb://localhost:27017/sample")

	log.Println("test")

	var err error
	client, err = mongo.NewClient(clientOpts)
	if err != nil {
		log.Fatalf("error while connecting to db: %v", err)
	}

	// give 10s timeout to connect to database
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("error while connecting to db: %v", err)
	}
	cancel()

	// ping db
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	cancel()

	log.Println("Connected to mongoDB.")
	return client
}

func GetMongoClient() *mongo.Client {
	if client == nil {
		client = ConnectMongo()
	}

	return client
}
