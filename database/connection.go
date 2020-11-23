package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection is to be used for CRUD operations on the mongodb database
type Connection struct {
	hostname string
	port     string
	username string
	password string
	dbname   string
	client   *mongo.Client
	cancel   context.CancelFunc
	context  context.Context
}

// ConnectToDb facilitates connecting to the mongodb
func (connection *Connection) ConnectToDb() bool {
	clientOptions := options.Client().ApplyURI(
		"mongodb://" +
			connection.username +
			":" + connection.password +
			"@" + connection.hostname +
			":" + connection.port,
	)

	ctx, cancel := context.WithCancel(context.Background())

	client, err := mongo.Connect(ctx, clientOptions)
	defer client.Disconnect(ctx)

	connection.client = client
	connection.cancel = cancel
	connection.context = ctx

	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

// Close is used to close the connection to the database
func (connection *Connection) Close() {
	connection.cancel()
}

// GetCollection returns the collection object using the collection name
func (connection *Connection) GetCollection(collectionName string) *mongo.Collection {

	if connection.client != nil {
		return connection.client.Database(connection.dbname).Collection(collectionName)
	}

	return nil
}

// AddRecord can be used to add records to any collection
func (connection *Connection) AddRecord(collection *mongo.Collection, record bson.D) bool {
	_, err := collection.InsertOne(connection.context, record)

	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

// UpdateRecord searches for and updates records based on provided conditions
func (connection *Connection) UpdateRecord(collection *mongo.Collection, filter bson.D, update bson.D) bool {
	_, err := collection.UpdateOne(connection.context, filter, update)

	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

//ToDO add functions to search for and delete records
