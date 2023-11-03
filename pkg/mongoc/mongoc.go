package mongoc

import (
	"blackbox-v2/pkg/config"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ConnString string
var DBName string
var MongoClient *mongo.Client
var Database *mongo.Database

func connect() {
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	log.Println("Connecting to MongoDB")
	log.Println(ConnString)
	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(ConnString))
	if err != nil {
		panic(err)
	}
	Database = MongoClient.Database(DBName)
}

func HandleMongoConn(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			MongoClient.Disconnect(ctx)
		}
	}
}
func loadConfig() {
	ConnString = config.Config.MongoDB.ConnString
	DBName = config.Config.MongoDB.Name
}

func init() {
	loadConfig()
	connect()
}
