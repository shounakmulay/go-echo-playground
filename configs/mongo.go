package configs

import (
	"awesomeProject/util/logger"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func mongoEnvUri() string {
	return os.Getenv("MONGO_URI")
}

func ConnectDB() {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoEnvUri()))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	MongoClient = client
	logger.Sugar.Info("Connected to MongoDB")
}

// MongoClient Client instance
var MongoClient *mongo.Client
