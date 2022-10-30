package db

import (
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "os"
    "log"
    "context"
)

func GetConnection() (conn *mongo.Database) {
    var uri string
    if uri = os.Getenv("MONGODB_URI"); uri == "" {
        log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
    }

    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        panic(err)
    }
    return client.Database("ostium")
}

func ReleaseConnection(db *mongo.Database) {
    client := db.Client()

    // TODO pool these
    if err := client.Disconnect(context.TODO()); err != nil {
        panic(err)
    }
    return
}



