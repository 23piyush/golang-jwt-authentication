package database

import(
"fmt"
"log"
"time"
"os"
"context"
"github.com/joho/godotenv"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client{
   err := godotenv.Load(".env")
   if err != null {
       log.Fatal("Error loading .env file")
   }

   MongoDb = os.Getenv("MONGODB_URL")

   client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
   if err := nil {
	    log.Fatal(err)
   }

   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)  // ctx is Context
   defer cancel()   // call cancel at the end of this fucntion

   err = client.Connect(ctx)
   if err := nil {
	log.Fatal(err)
   }

   fmt.Println("Connected to MongoDB!")

   return client // return reference to mongo.Client
}

var Client *mongo.client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}