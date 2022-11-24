package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bot struct {
	Name       string `json:"name"`
	FilesCount int    `json:"files_count"`
}

type Boter interface {
	ListAll() ([]Bot, error)
	ScrapeAll() error
	Scrape(botName string) error
	Stats(botName string) (string, error)
}

type File struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	StartTime time.Time          `json:"start_time,omitempty" bson:"start_time"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
}

type FileLog struct {
	ID                               primitive.ObjectID `bson:"_id"`
	LogCountInfo                     int                `bson:"log_count/INFO"`
	LogCountDebug                    int                `bson:"log_count/DEBUG"`
	StartTime                        time.Time          `bson:"start_time"`
	SchedulerEnqueuedMemory          int                `bson:"scheduler/enqueued/memory"`
	SchedulerEnqueued                int                `bson:"scheduler/enqueued"`
	SchedulerDenqueuedMemory         int                `bson:"scheduler/dequeued/memory"`
	SchedulerDenqueued               int                `bson:"scheduler/dequeued"`
	DownloaderRequestCount           int                `bson:"downloader/request_count"`
	DownloaderRequestMethodCount     int                `bson:"downloader/request_method_count/GET"`
	DownloaderRequestBytes           int                `bson:"downloader/request_bytes"`
	RobotstxtRequestCount            int                `bson:"robotstxt/request_count"`
	DownloaderResponseCount          int                `bson:"downloader/response_count"`
	DownloaderResponseStatusCount404 int                `bson:"downloader/response_status_count/404"`
	DownloaderResponseBytes          int                `bson:"downloader/response_bytes"`
	HttpCompressionResponseBytes     int                `bson:"httpcompression/response_bytes"`
	HttpCompressionResponseCount     int                `bson:"httpcompression/response_count"`
	ResponseReceivedCount            int                `bson:"response_received_count"`
	RobotstxtResponseCount           int                `bson:"robotstxt/response_count"`
	RobotstxtResponseStatusCount404  int                `bson:"robotstxt/response_status_count/404"`
	DownloaderResponseStatusCount200 int                `bson:"downloader/response_status_count/200"`
	ItemScrapedCount                 int                `bson:"item_scraped_count"`
	RequestDepthMax                  int                `bson:"request_depth_max"`
}
