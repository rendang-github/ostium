package models
import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

/// Candle Instance
type Campaign struct {
    Id *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name string `json:"name" bson:"name"`
    Description string `json:"description" bson:"description"`
    Creator primitive.ObjectID `json:"creator" bson:"creator"`
    Created time.Time `json:"created" bson:"created"`
    Modified time.Time `json:"modified" bson:"modified"`
    Root primitive.ObjectID `json:"root" bson:"root"`
    Layout primitive.ObjectID `json:"layout" bson:"layout"`
}

