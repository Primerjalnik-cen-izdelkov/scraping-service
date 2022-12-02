package repo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"scraping_service/pkg/common"
	"scraping_service/pkg/models"
	"strconv"
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
	// uri := "mongodb://localhost:27017"
	//uri := "mongodb://mongo_server:27018"
	uri := os.Getenv("MONGODB_URI")
	fmt.Println("my uri is: ", uri)
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
		fmt.Println("mongo cursor err:", err)
		return nil, ErrMongoCursor
	}

	var results []models.File
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, ErrMongoDecoding
	}

	return results, nil
}

func (db MongoDB) GetBotLogs(botName string, qm map[string]string) ([]models.FileLog, error) {
    // TODO(miha): We need to rename collection to logs, rename also in python
    // file!
	botCollection := db.Client.Database("stats").Collection(botName)

    q := []bson.D{bson.D{{}}}
    p := []bson.D{}
    s := []bson.D{}

    // NOTE(miha): Append querry parameters
    if len(qm["querry"]) == 0 {
        // NOTE(miha): time.lt querry parameter.
        if len(qm["timeLT"]) != 0 {
            lt, err := common.QuerryParamParseTime(qm["timeLT"])
            if err != nil {
                // TODO
            }
            q = append(q, bson.D{{"start_time", bson.D{{"$lt", lt}}}})
        }
        if len(qm["timeGT"]) != 0 {
            gt, err := common.QuerryParamParseTime(qm["timeGT"])
            if err != nil {
                // TODO
            }
            q = append(q, bson.D{{"start_time", bson.D{{"$gte", gt}}}})
        }
        if len(qm["itemsScrapedLT"]) != 0 {
            lt, err := strconv.ParseInt(qm["itemsScrapedLT"], 10, 64)
            if err != nil {
                return nil, err
            }
            q = append(q, bson.D{{"items_scraped_count", bson.D{{"$lt", lt}}}})
        }
        if len(qm["itemsScrapedGT"]) != 0 {
            gt, err := strconv.ParseInt(qm["itemsScrapedGT"], 10, 64)
            if err != nil {
                return nil, err
            }
            q = append(q, bson.D{{"items_scraped_count", bson.D{{"$gte", gt}}}})
        }
    } else {
        q = append(q, bson.D{{"", qm["querry"]}})
    }

    if len(qm["projection"]) == 0 {

    } else {
        p = append(p, bson.D{{"", qm["projection"]}})
    }

    if len(qm["sort"]) == 0 {

    } else {
        s = append(s, bson.D{{"", qm["sort"]}})
    }

    fmt.Println("mongodb querry:", bson.M{"$and": q})

	//projection := bson.D{{"start_time", 1}}
    sort := bson.D{{"start_time", -1}}
    cursor, err := botCollection.Find(db.Context,  bson.M{"$and": q}, options.Find().SetSort(sort)) //SetProjection(projection).SetSort(sort))
	if err != nil {
		fmt.Println("mongo cursor err:", err)
		return nil, ErrMongoCursor
	}

	var results []models.FileLog
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, ErrMongoDecoding
	}

	return results, nil
}

func (db MongoDB) GetFileLogs(botName string, unixTime int64) (*models.FileLog, error) {
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

	var result models.FileLog
	err := cursor.Decode(&result)
	if err != nil {
		fmt.Println("err:", err)
		return nil, ErrMongoDecoding
	}

	return &result, nil
}


func (db MongoDB) UpdateBot(botName string) error {
    coll := db.Client.Database("dev").Collection("bots")
    filter := bson.D{{"bot_name", botName}}
    update := bson.M{"$set": bson.D{{"last_run", time.Now()}},
                     "$inc": bson.D{{"logs_count", 1}}}
    uo := options.Update().SetUpsert(true)
    _, err := coll.UpdateOne(db.Context, filter, update, uo)
    if err != nil {
        return err
    }

    return nil
}

func (db MongoDB) GetBot(botName string) (*models.Bot, error) {
    var bot models.Bot

    coll := db.Client.Database("dev").Collection("bots")
	cursor := coll.FindOne(db.Context, bson.D{{"bot_name", botName}}, options.FindOne())

	err := cursor.Decode(&bot)
	if err != nil {
		fmt.Println("err:", err)
		return nil, ErrMongoDecoding
	}

    fmt.Println("GetBot bot:", bot)

    return &bot, nil
}
