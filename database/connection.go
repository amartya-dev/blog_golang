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
	Hostname string
	Port     string
	Username string
	Password string
	Dbname   string
	Client   *mongo.Client
	cancel   context.CancelFunc
	context  context.Context
}

// ConnectToDb facilitates connecting to the mongodb
func (connection *Connection) ConnectToDb() bool {
	clientOptions := options.Client().ApplyURI(
		"mongodb://" +
			connection.Username +
			":" + connection.Password +
			"@" + connection.Hostname +
			":" + connection.Port,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer connection.Close()

	client, err := mongo.Connect(ctx, clientOptions)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	connection.Client = client
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

	if connection.Client != nil {
		return connection.Client.Database(connection.Dbname).Collection(collectionName)
	}

	return nil
}

// AddRecord can be used to add records to any collection
func (connection *Connection) AddRecord(collection *mongo.Collection, record bson.M) error {
	_, err := collection.InsertOne(connection.context, record)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
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

//DeleteRecord searches for and deletes a record based on provided conditions
func (connection *Connection) DeleteRecord(collection *mongo.Collection, filter bson.D) bool {
	_, err := collection.DeleteOne(connection.context, filter)

	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

//GetRecords finds and returns the records with specified search conditions and options
func (connection *Connection) GetRecords(collection *mongo.Collection, conditions bson.D, opts *options.FindOptions) []bson.M {
	cursor, err := collection.Find(connection.context, conditions, opts)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// get a list of all returned documents and return them
	// see the mongo.Cursor documentation for more examples of using cursors
	var results []bson.M
	if err = cursor.All(connection.context, &results); err != nil {
		log.Fatal(err)
	}

	return results
}

//GetRecord finds and returns a single order with based on the specified query
func (connection *Connection) GetRecord(collection *mongo.Collection, condition bson.D) bson.M {
	var result bson.M
	err := collection.FindOne(connection.context, condition).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Print("No documents match the query")
			return nil
		}
		log.Fatal(err)
	}
	return result
}
