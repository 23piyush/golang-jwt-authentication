package database
// sets up a MongoDB client connection and provides functions to access a MongoDB collection

import(
"fmt"
"log"
"time"
"os"
"context"
"github.com/joho/godotenv"  //for loading environment variables from a .env file
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client{
   err := godotenv.Load(".env")
   if err != nil {
       log.Fatal("Error loading .env file")
   }

   MongoDb := os.Getenv("MONGODB_URL")  // MONGODB_URL is the connection URL for the MongoDB server

   client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb)) // craetes a mongoDB client for this MongoDB connection URL
   if err != nil {
	    log.Fatal(err)
   }
 
   // context is used to manage the connection and request deadlines
   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)  // ctx is Context
   defer cancel()   // call cancel at the end of this fucntion

   err = client.Connect(ctx) // connect client to the mongoDB server
   if err != nil {
	log.Fatal(err)
   }

   fmt.Println("Connected to MongoDB!")

   return client // return reference to mongo.Client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}