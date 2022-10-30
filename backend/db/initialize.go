package db

import (
    "go.mongodb.org/mongo-driver/mongo"
    "context"
    "log"
)

const (
    namespaceExistsErrCode int32 = 48
)

func createCollection(conn *mongo.Database, label string) {
    if err := conn.CreateCollection(context.TODO(), label); err != nil {
        cmdErr, _ := err.(mongo.CommandError)
        if cmdErr.Code != 48 {
            log.Fatal(err)
        } else {
            log.Printf("Collection '" + label + "' exists")
        }
    } else {
        log.Printf("Collection '" + label + "' created")
    }
}

func Initialize() {
    // Get a connection
    conn := GetConnection()
    defer ReleaseConnection(conn)

    createCollection(conn, "campaign")
}

