package mongo_service

import "go.mongodb.org/mongo-driver/mongo"

// GetCollection getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("go-echo").Collection(collectionName)
	return collection
}
