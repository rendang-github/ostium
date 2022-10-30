package models
import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

/// Candle Instance
type Campaign struct {
    Id *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name string `json:"name" db:"name"`
    Description string `json:"description" db:"description"`
    Created time.Time `json:"created" db:"created"`
    Modified time.Time `json:"modified" db:"modified"`
}

