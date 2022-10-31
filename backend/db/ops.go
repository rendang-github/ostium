package db

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func Insert(model string, document interface{}) (primitive.ObjectID) {
    // Get a connection
    conn := GetConnection()
    defer ReleaseConnection(conn)
    coll := conn.Collection(model)

    result, err := coll.InsertOne(context.TODO(), document)
    if err != nil {
        panic(err)
    }

    return result.InsertedID.(primitive.ObjectID)
}

func Get(output interface {}, model string, id primitive.ObjectID) (err error) {
    // Get a connection
    conn := GetConnection()
    defer ReleaseConnection(conn)
    coll := conn.Collection(model)

    result := coll.FindOne(context.TODO(), bson.M{"_id": id})
    return result.Decode(output)
}

func Query(output interface {}, model string, filter interface{}) (err error) {
    // Get a connection
    conn := GetConnection()
    defer ReleaseConnection(conn)
    coll := conn.Collection(model)

    result := coll.FindOne(context.TODO(), filter)
    return result.Decode(output)
}

func Delete(model string, id primitive.ObjectID) (err error) {
    // Get a connection
    conn := GetConnection()
    defer ReleaseConnection(conn)
    coll := conn.Collection(model)

    _, err = coll.DeleteOne(context.TODO(), bson.M{"_id": id})
    return
}

func Set(model string, document interface {}, id primitive.ObjectID) (err error) {
    // Get a connection
    conn := GetConnection()
    defer ReleaseConnection(conn)
    coll := conn.Collection(model)

    _, err = coll.ReplaceOne(context.TODO(), bson.M{"_id": id}, document)
    return err
}

func Update(model string, document interface {}, id primitive.ObjectID) (err error) {
    // Get a connection
    conn := GetConnection()
    defer ReleaseConnection(conn)
    coll := conn.Collection(model)

    _, err = coll.UpdateOne(context.TODO(), bson.M{"_id": id}, document)
    return err
}

func All(output interface {}, model string) (err error) {
    // Get a connection
    conn := GetConnection()
    defer ReleaseConnection(conn)
    coll := conn.Collection(model)

    // Get all records
    cursor, err := coll.Find(context.TODO(), bson.D{{}})
    if err != nil {
        return err
    }

    return cursor.All(context.TODO(), output)
}
