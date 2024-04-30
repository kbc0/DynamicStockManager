package main

import (
    "context"
    "fmt"
    "log"
	"os"

    "github.com/rs/zerolog"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/kbc0/DynamicStockManager/server" // Adjust this import path to your actual path
)

func main() {
    // MongoDB connection setup
    serverAPI := options.ServerAPI(options.ServerAPIVersion1)
    opts := options.Client().ApplyURI("mongodb+srv://kbc0:123456admin@dynamicstockmanagement.jdmc40a.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
    client, err := mongo.Connect(context.TODO(), opts)
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err = client.Disconnect(context.TODO()); err != nil {
            log.Fatal(err)
        }
    }()

    // Check MongoDB connection
    if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Successfully connected to MongoDB!")

    // Logger setup
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

    // Initialize the server with the database and logger
    srv := server.NewServer(client.Database("Users"), &logger) // Adjust "Users" to your actual user database name if different

    // Start the server
    log.Fatal(srv.App.Listen(":8080"))
}
