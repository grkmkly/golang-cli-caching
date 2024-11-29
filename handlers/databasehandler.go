package handlers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"main.go/model"
)

// CONNECT DATABASE
func Connect(db *model.Database) error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	uri := os.Getenv("MONGODB")
	if uri == "" {
		return errors.New("NOURI")
	}
	//clientoptions
	clientOptions := options.Client().ApplyURI(uri)
	db.ClientOption = clientOptions
	//context
	db.Ctx = context.TODO()
	// database
	db.Database = os.Getenv("DATABASENAME")

	var err error
	//client
	db.Client, err = mongo.Connect(db.Ctx, db.ClientOption)
	if err != nil {
		return err
	}
	err = db.Client.Ping(db.Ctx, nil)
	if err != nil {
		return err
	}
	fmt.Println("Connected database")
	return nil
}
func SetCollection(db *model.Database, collectionName string, chanCollection chan bool) {
	db.Collection = db.Client.Database(db.Database).Collection(collectionName)
	fmt.Printf("Connected Collection : %v\n", db.Collection.Name())
	chanCollection <- true
}
func InsertLinkPort(db *model.Database, linkport *model.LinkPort) error {
	linkport.CreatedAt = time.Now()
	linkport.Id = primitive.NewObjectID()
	_, err := db.Collection.InsertOne(db.Ctx, linkport)
	if err != nil {
		return err
	}
	fmt.Println("Inserted")
	return nil
}
func CheckLinkPort(db *model.Database, item model.LinkPort) (bool, string) {
	filterLink := bson.D{
		{Key: "link", Value: item.Link},
	}
	filterPort := bson.D{
		{Key: "port", Value: item.Port},
	}
	resultLink := db.Collection.FindOne(db.Ctx, filterLink)
	var linkModel model.LinkPort
	resultLink.Decode(&linkModel)

	resultPort := db.Collection.FindOne(db.Ctx, filterPort)
	var portModel model.LinkPort
	resultPort.Decode(&portModel)
	if linkModel.Link == item.Link {
		return true, linkModel.Port // Link mevcut o yüzden o linkin portuna ulaş ve isteği yap
	} else if portModel.Port == item.Port {
		return false, "ACTIVE" // o port açık o yüzden portu değiştirmesini söyle
	}
	return false, item.Port
}
func DeleteLinkPort(db *model.Database, item model.LinkPort) error {
	filter := bson.D{
		{Key: "link", Value: item.Link},
		{Key: "port", Value: item.Port},
	}
	resultDeleted, err := db.Collection.DeleteOne(db.Ctx, filter)
	if err != nil {
		return err
	}
	fmt.Println("resultDeleted", resultDeleted)
	return nil
}
