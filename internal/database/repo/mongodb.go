package repo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"scraping_service/pkg/common"
	"scraping_service/pkg/models"
	"time"
    "net/url"

    //"github.com/pasztorpisti/qs"
    "github.com/ahmetcanozcan/fet"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ErrMongoCursor        = errors.New("MongoDB cursor returned an error.")
	ErrMongoDecoding      = errors.New("MongoDB couldn't decode results.")
	ErrQSUnmarshal        = errors.New("Library 'qs' couldn't unmarshal query parameters to the struct.")
    ErrObjectIDConversion = errors.New("Can't convert given string to the primitive.objectID")
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

/*
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
*/

func (db MongoDB) GetBotLogs(botName string, qm map[string]string) ([]models.FileLog, error) {
    return nil, nil
    /*
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
    */
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

// Get information on bot with name 'botName' from mongoDB.
func (db MongoDB) GetBot(botName string, qp url.Values) (*models.Bot, error) {
    type Query struct {
        ID            []string    `qs:"id,nil"`
        LastRun       []time.Time `qs:"last_run,nil"`
        LastRunGt     time.Time   `qs:"last_run.gt,nil"`
        LastRunLt     time.Time   `qs:"last_run.lt,nil"`
        LogsCount     []int       `qs:"logs_count,nil"`
        LogsCountGt   int         `qs:"logs_count.gt,nil"`
        LogsCountLt   int         `qs:"logs_count.lt,nil"`
        Name          []string    `qs:"name,nil"`
        Limit         int         `qs:"limit,nil"`
        Sort          []string    `qs:"sort,nil"`
        Field         []string    `qs:"field,nil"`
    }
    var q Query
    err := common.CustomUnmarshaler.UnmarshalValues(&q, qp)
    if err != nil {
        return nil, ErrQSUnmarshal
    }

    // NOTE(miha): Build mongoDB query with the help of the 'fet' library.
    u := []fet.Updater{}
    u, err = common.QueryID(u, q.ID, "_id")
    if err != nil {
        return nil, ErrObjectIDConversion
    }
    u = common.QueryStringSlice(u, q.Name, "bot_name")
    u = common.QuerySameDay(u, q.LastRun, "last_run")
    u = common.QueryDateGreater(u, q.LastRunGt, "last_run")
    u = common.QueryDateLess(u, q.LastRunLt, "last_run")
    u = common.QueryIntSlice(u, q.LogsCount, "logs_count")
    u = common.QueryIntGreater(u, q.LogsCountGt, "logs_count")
    u = common.QueryIntLess(u, q.LogsCountLt, "logs_count")
    filter := fet.Build(u...)
    _ = filter

    // NOTE(miha): Set limit and sort options if the query parameter exists.
    opts := options.FindOne()
    opts = common.QueryOneOptsProjection(opts, q.Field)

    // NOTE(miha): Select collection and execute query to find bot with
    // 'botName'.
    coll := db.Client.Database("dev").Collection("bots")
	cursor := coll.FindOne(db.Context, bson.D{{"bot_name", botName}}, opts)

    // NOTE(miha): Decode results into variable 'bot' or return error.
    var bot models.Bot
	err = cursor.Decode(&bot)
	if err != nil {
		return nil, ErrMongoDecoding
	}

    return &bot, nil
}

// Parse query parameters, build up query and execute it on mongoDB database
// and get files.
func (db MongoDB) GetFiles(qp url.Values) ([]models.File, error) {
    // NOTE(miha): Unmarshal query parameters into struct bellow.
    type Query struct {
        ID     []string    `qs:"id,nil"`
        Date   []time.Time `qs:"date,nil"`
        DateGt time.Time   `qs:"date.gt,nil"`
        DateLt time.Time   `qs:"date.lt,nil"`
        Name   []string    `qs:"name,nil"`
        Limit  int         `qs:"limit,nil"`
        Sort   []string    `qs:"sort,nil"`
        Field  []string    `qs:"field,nil"`
        FileName []string  `qs:"file_name,nil"`
    }
    var q Query
    err := common.CustomUnmarshaler.UnmarshalValues(&q, qp)
    if err != nil {
        return nil, ErrQSUnmarshal
    }

    // NOTE(miha): Build mongoDB query with the help of the 'fet' library.
    u := []fet.Updater{}
    u, err = common.QueryID(u, q.ID, "_id")
    if err != nil {
        return nil, ErrObjectIDConversion
    }
    u = common.QueryStringSlice(u, q.Name, "bot_name")
    u = common.QueryStringSlice(u, q.FileName, "file_name")
    u = common.QuerySameDay(u, q.Date, "date")
    u = common.QueryDateGreater(u, q.DateGt, "date")
    u = common.QueryDateLess(u, q.DateLt, "date")
    filter := fet.Build(u...)

    // NOTE(miha): Set limit and sort options if the query parameter exists.
    opts := options.Find()
    opts = common.QueryOptsLimit(opts, q.Limit)
    opts = common.QueryOptsSort(opts, q.Sort)
    opts = common.QueryOptsProjection(opts, q.Field)

    // NOTE(miha): Select collection in mongoDB and execute query with options.
    coll := db.Client.Database("files").Collection("files")
    cursor, err := coll.Find(db.Context, filter, opts)
    if err != nil {
		return nil, ErrMongoCursor
    }

    // NOTE(miha): Decode result from mongoDB to array of structs.
    var files []models.File
	if err = cursor.All(context.TODO(), &files); err != nil {
		return nil, ErrMongoDecoding
	}
    
    return files, nil
}

// Parse query parameters, build up query and execute it on mongoDB database
// and get logs.
func (db MongoDB) GetLogs(qp url.Values) ([]models.FileLog, error) {
    // NOTE(miha): Unmarshal query parameters into struct bellow.
    type Query struct {
        ID                                 []string    `qs:"id,nil"`
        StartTime                          []time.Time `qs:"start_time,nil"`
        StartTimeGt                        time.Time   `qs:"start_time.gt,nil"`
        StartTimeLt                        time.Time   `qs:"start_time.lt,nil"`
        DownloaderRequestCount             []int       `qs:"request_count,nil"`
        DownloaderRequestCountGt           int         `qs:"request_count.gt,nil"`
        DownloaderRequestCountLt           int         `qs:"request_count.lt,nil"`
        DownloaderResponseCount            []int       `qs:"response_count,nil"`
        DownloaderResponseCountGt          int         `qs:"response_count.gt,nil"`
        DownloaderResponseCountLt          int         `qs:"response_count.lt,nil"`
        DownloaderResponseStatusCount404   []int       `qs:"404,nil"`
        DownloaderResponseStatusCount404Gt int         `qs:"404.gt,nil"`
        DownloaderResponseStatusCount404Lt int         `qs:"404.lt,nil"`
        ItemScrapedCount                   []int       `qs:"item_scraped,nil"`
        ItemScrapedCountGt                 int         `qs:"item_scraped.gt,nil"`
        ItemScrapedCountLt                 int         `qs:"item_scraped.lt,nil"`
        Name                               []string    `qs:"name,nil"`
        Limit                              int         `qs:"limit,nil"`
        Sort                               []string    `qs:"sort,nil"`
        Field                              []string    `qs:"field,nil"`
    }
    var q Query
    err := common.CustomUnmarshaler.UnmarshalValues(&q, qp)
    if err != nil {
        return nil, ErrQSUnmarshal
    }

    // NOTE(miha): Build mongoDB query with the help of the 'fet' library.
    u := []fet.Updater{}
    u, err = common.QueryID(u, q.ID, "_id")
    if err != nil {
        return nil, ErrObjectIDConversion
    }
    u = common.QuerySameDay(u, q.StartTime, "start_time")
    u = common.QueryDateGreater(u, q.StartTimeGt, "start_time")
    u = common.QueryDateLess(u, q.StartTimeLt, "start_time")
    u = common.QueryIntSlice(u, q.DownloaderRequestCount, "downloader_request_count")
    u = common.QueryIntGreater(u, q.DownloaderRequestCountGt, "downloader_request_count")
    u = common.QueryIntLess(u, q.DownloaderRequestCountLt, "downloader_request_count")
    u = common.QueryIntSlice(u, q.DownloaderResponseCount, "downloader_response_count")
    u = common.QueryIntGreater(u, q.DownloaderResponseCountGt, "downloader_response_count")
    u = common.QueryIntLess(u, q.DownloaderResponseCountLt, "downloader_response_count")
    u = common.QueryIntSlice(u, q.DownloaderResponseStatusCount404, "downloader_response_status_count_404")
    u = common.QueryIntGreater(u, q.DownloaderResponseStatusCount404Gt, "downloader_response_status_count_404")
    u = common.QueryIntLess(u, q.DownloaderResponseStatusCount404Lt, "downloader_response_status_count_404")
    u = common.QueryIntSlice(u, q.ItemScrapedCount, "item_scraped_count")
    u = common.QueryIntGreater(u, q.ItemScrapedCountGt, "item_scraped_count")
    u = common.QueryIntLess(u, q.ItemScrapedCountLt, "item_scraped_count")
    u = common.QueryStringSlice(u, q.Name, "bot_name")
    filter := fet.Build(u...)

    // NOTE(miha): Set limit and sort options if the query parameter exists.
    opts := options.Find()
    opts = common.QueryOptsLimit(opts, q.Limit)
    opts = common.QueryOptsSort(opts, q.Sort)
    opts = common.QueryOptsProjection(opts, q.Field)

    // NOTE(miha): Select collection in mongoDB and execute query with options.
    // TODO(miha): Put db and collection into env variable.
    coll := db.Client.Database("stats").Collection("stats1")
    cursor, err := coll.Find(db.Context, filter, opts)
    if err != nil {
		return nil, ErrMongoCursor
    }

    // NOTE(miha): Decode result from mongoDB to array of structs.
    var logs []models.FileLog
	if err = cursor.All(context.TODO(), &logs); err != nil {
		return nil, ErrMongoDecoding
	}
    
    return logs, nil
}
