package db

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "os"
    "ostium/config"
)

func GetConnection() (conn *mongo.Database) {
    var uri string
    if uri = os.Getenv("MONGODB_URI"); uri == "" {
        uri = config.DatabaseURI
    }

    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        panic(err)
    }
    return client.Database(config.DatabaseName)
}

func ReleaseConnection(db *mongo.Database) {
    client := db.Client()

    // TODO pool these
    if err := client.Disconnect(context.TODO()); err != nil {
        panic(err)
    }
    return
}



