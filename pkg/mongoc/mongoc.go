package mongoc

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ConnString = "mongodb://localhost:27017"
var MongoClient *mongo.Client
var Database *mongo.Database
var DBName = "blackboxv2"

func connect() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

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

func init() {
	connect()
}
