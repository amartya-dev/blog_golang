package dao

import (
	"os"

	"blog/database"
	"blog/models"
)

// UserDao is used for communicating with the Db and performing the CRUD operations
type UserDao struct {
	connection *database.Connection
	user       *models.User
}

func (ur *UserDao) initialize() {
	connection := database.Connection{}
	connection.Hostname = os.Getenv("MONGODB_HOST")
	connection.Dbname = os.Getenv("MONGODB_NAME")
	connection.Password = os.Getenv("MONGODB_PASSWORD")
	connection.Port = os.Getenv("MONGODB_PORT")
	connection.Username = os.Getenv("MONGODB_USER")

	if result := connection.ConnectToDb(); !result {
		panic("Could not connect to database")
	} else {
		ur.connection = &connection
		defer connection.Close()
	}
}

// func (ur *UserDao) saveUser(u *models.User) {
// 	if ur.connection != nil {
// 		collection := ur.connection.GetCollection("user")
// 		data, err := bson.Marshal(u)
// 		if err != nil {
// 			panic(err)
// 		}

// 		ur.connection.AddRecord()
// 	} else {
// 		ur.initialize()
// 	}
// }
