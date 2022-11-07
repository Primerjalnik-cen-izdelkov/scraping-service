package repo

import (
	"context"
	"errors"
	"fmt"
	"scraping_service/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ErrMongoCursor   = errors.New("MongoDB cursor returned an error.")
	ErrMongoDecoding = errors.New("MongoDB couldn't decode results.")
)

type MongoDB struct {
	Client  *mongo.Client
	Context context.Context
}

func CreateMongoDB() *MongoDB {
	uri := "mongodb://localhost:27017"
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("client err:", err)
	}

	return &MongoDB{Client: client, Context: ctx}
}

func (db MongoDB) Ping() error {
	err := db.Client.Ping(db.Context, readpref.Primary())
	if err != nil {
		return err
	}

	return nil
}

func (db MongoDB) GetBotFileNames(botName string) ([]models.File, error) {
	botCollection := db.Client.Database("stats").Collection(botName)

	projection := bson.D{{"start_time", 1}}
	sort := bson.D{{"start_time", -1}}
	cursor, err := botCollection.Find(db.Context, bson.D{}, options.Find().SetProjection(projection).SetSort(sort))
	if err != nil {
		return nil, ErrMongoCursor
	}

	var results []models.File
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, ErrMongoDecoding
	}

	return results, nil
}

func (db MongoDB) GetBotStats(botName, query string) ([]models.FileStat, error) {
	botCollection := db.Client.Database("stats").Collection(botName)

	//projection := bson.D{{"start_time", 1}}
	sort := bson.D{{"start_time", -1}}
	cursor, err := botCollection.Find(db.Context, bson.D{}, options.Find().SetSort(sort)) //SetProjection(projection).SetSort(sort))
	if err != nil {
		return nil, ErrMongoCursor
	}

	var results []models.FileStat
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, ErrMongoDecoding
	}

	return results, nil
}

func (db MongoDB) GetFileStats(botName string, unixTime int64) (*models.FileStat, error) {
	botCollection := db.Client.Database("stats").Collection(botName)

	/*
		{"start_time":
		{
		    "$gte" : ISODate(1667572867000),
		    "$lt" : ISODate(1667572868000)
		}}
	*/

	gte := primitive.NewDateTimeFromTime(time.Unix(unixTime, 0))
	lt := primitive.NewDateTimeFromTime(time.Unix(unixTime+1, 0))

	cursor := botCollection.FindOne(db.Context, bson.D{{"start_time", bson.D{{"$gte", gte}, {"$lt", lt}}}}, options.FindOne()) //SetProjection(projection).SetSort(sort))

	var result models.FileStat
	err := cursor.Decode(&result)
	if err != nil {
		fmt.Println("err:", err)
		return nil, ErrMongoDecoding
	}

	return &result, nil
}
