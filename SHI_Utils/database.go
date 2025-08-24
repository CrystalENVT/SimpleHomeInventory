package shiutils

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	// global access Mongo Client
	Mongo_Client         *mongo.Client
	mongo_off_collection *mongo.Collection
)

func InitDB() {
	uri := os.Getenv("MONGO_URI")
	println("uri:", uri)

	var db_connect_err error
	Mongo_Client, db_connect_err = mongo.Connect(options.Client().
		ApplyURI(uri))
	if db_connect_err != nil {
		println(db_connect_err.Error())
		panic(db_connect_err)
	}

	mongo_off_collection = Mongo_Client.Database("shi").Collection("OpenFoodFacts_Data")
	mongo_off_collection.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{Keys: bson.D{{Key: "code", Value: 1}}, Options: options.Index().SetUnique(true)},
	)
	mongo_off_collection.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{Keys: bson.D{{Key: "product_name", Value: 1}}},
	)
}

func upcLookupViaDB(upc_string string) (found bool, results MinimalProduct, err error) {
	fmt.Println("DB Lookup UPC String:", upc_string)

	filter := bson.D{{Key: "code", Value: upc_string}}

	fmt.Println("Filter is: ", filter)

	err = mongo_off_collection.FindOne(context.TODO(), filter).Decode(&results)

	fmt.Println("result:", prettyPrint(results))

	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No documents found in MongoDB UPC Cache")
			return false, results, nil
		}
		return false, results, err
	}
	return true, results, nil
}

func writeUPCtoDB(product MinimalProduct) (err error) {
	println("Attempting to write document to DB")
	result, err := mongo_off_collection.InsertOne(context.TODO(), product)
	if err == nil {
		fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	} else if mongo.IsDuplicateKeyError(err) {
		// Not an error for the program, as each UPC is unique.
		//   However, log this as we shouldn't have reached this point
		fmt.Println("Unexpected Duplicate Key / UPC Code -", product.Code, ":", err)
		return nil
	}
	return err
}
