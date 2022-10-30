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
    Created time.Time `json:"created" bson:"created"`
    Modified time.Time `json:"modified" bson:"modified"`
}

