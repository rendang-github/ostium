package config

// MongoDB connection URI
var DatabaseURI string

// Database name
var DatabaseName string

func init() {
    DatabaseURI = "mongodb://127.0.0.1:27017"
    DatabaseName = "ostium"
}
