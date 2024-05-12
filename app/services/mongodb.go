package services

import (
	"context"
	"log"
	"sync"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"testproject/global"
	"fmt"
)

// type DataBaseName;
// type DataBaseCollection;
var DataBaseCollection string = "log" 
var DataBaseName string = "gin" 
type mongoDBService struct {
	client *mongo.Client  
	once   sync.Once  
}

// type mongoDBSet struct {
// 	DataBaseCollection string 
// 	DataBaseName   string  
// }


var MongoDBService = new(mongoDBService)

// var MongoDBSet = new(mongoDBSet)
func (mongoDBService *mongoDBService) Connect() (*mongo.Client,error) {
	// fmt.Println(mongoDBService);
	fmt.Println("測試鏈接");
	mongoDBService.once.Do(func() {
		opts := options.Client().ApplyURI("mongodb://"+global.App.Config.MongoDB.User+":"+global.App.Config.MongoDB.Password+"@"+global.App.Config.MongoDB.Host+"")
		client, err := mongo.Connect(context.TODO(), opts)
		fmt.Println(err);
		fmt.Println("test");
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connected to MongoDB!")  
		mongoDBService.client = client  
	})



	return mongoDBService.client, nil
}

//查詢數據
func (mongoDBService *mongoDBService) Find(filter bson.M) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := mongoDBService.client.Database(DataBaseName).Collection(DataBaseCollection)
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var results []bson.M
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	// fmt.Println(results);
	return results, nil
}

//插入一條
func (mongoDBService *mongoDBService) InsertOne(doc bson.M) error {
	fmt.Println("測試添加");
	fmt.Println(mongoDBService.client);
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := mongoDBService.client.Database(DataBaseName).Collection(DataBaseCollection)
	_, err := collection.InsertOne(ctx, doc)
	return err
}

//批量插入
func (mongoDBService *mongoDBService) InsertMany(doc []interface{}) error {
	fmt.Println("測試批量添加");
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := mongoDBService.client.Database(DataBaseName).Collection(DataBaseCollection)
	_, err := collection.InsertMany(ctx, doc)
	return err
}

//修改一條
func (mongoDBService *mongoDBService) UpdateOne(filter bson.M, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := mongoDBService.client.Database(DataBaseName).Collection(DataBaseCollection)
	_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	return err
}


