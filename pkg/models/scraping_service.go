package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bot struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
    Name       string             `json:"name" bson:"bot_name"`
    Status     *BotStatus         `json:"status" bson:"status"`
    LastRun    time.Time          `json:"last_run" bson:"last_run"`
    LogsCount  int                `json:"logs_count" bson:"logs_count"`
}

type Boter interface {
	ListAll() ([]Bot, error)
	ScrapeAll() error
	Scrape(botName string) error
	Stats(botName string) (string, error)
}

type BotStatus struct {
    Running bool `json:"running"`
}

type File struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Date      time.Time          `json:"date,omitempty" bson:"date,omitempty"`
	BotName   string             `json:"bot_name,omitempty" bson:"bot_name,omitempty"`
	FileName  string             `json:"file_name,omitempty" bson:"file_name,omitempty"`
}

type FileLog struct {
	ID                               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
    BotNamee                         string             `json:"bot_name,omitempty" bson:"bot_name,omitempty"`
	//LogCountInfo                     int                `json:"log_count_info,omitempty" bson:"log_count_info,omitempty"`
	//LogCountDebug                    int                `json:"log_count_debug,omitempty" bson:"log_count_debug,omitempty"`
	StartTime                        time.Time          `json:"start_time,omitempty" bson:"start_time,omitempty"`
	//SchedulerEnqueuedMemory          int                `json:"scheduler_enqueued_memory,omitempty" bson:"scheduler_enqueued_memory,omitempty"`
	//SchedulerEnqueued                int                `json:"scheduler_enqueued,omitempty" bson:"scheduler_enqueued,omitempty"`
	//SchedulerDequeuedMemory          int                `json:"scheduler_dequeued_memory,omitempty" bson:"scheduler_dequeued_memory,omitempty"`
	//SchedulerDequeued                int                `json:"scheduler_dequeued,omitempty" bson:"scheduler_dequeued,omitempty"`
	DownloaderRequestCount           int                `json:"downloader_request_count,omitempty" bson:"downloader_request_count,omitempty"`
	//DownloaderRequestMethodCount     int                `json:"downloader_request_method_count,omitempty" bson:"downloader_request_method_count,omitempty"`
	//DownloaderRequestBytes           int                `json:"downloader_request_bytes,omitempty" bson:"downloader_request_bytes,omitempty"`
	//RobotstxtRequestCount            int                `json:"robotstxt_request_count,omitempty" bson:"robotstxt_request_count,omitempty"`
	DownloaderResponseCount          int                `json:"downloader_response_count,omitempty" bson:"downloader_response_count,omitempty"`
	DownloaderResponseStatusCount404 int                `json:"downloader_response_status_count_404,omitempty" bson:"downloader_response_status_count_404,omitempty"`
	//DownloaderResponseBytes          int                `json:"downloader_response_bytes,omitempty" bson:"downloader_response_bytes,omitempty"`
	//HttpCompressionResponseBytes     int                `json:"httpcompression_response_bytes,omitempty" bson:"httpcompression_response_bytes,omitempty"`
	//HttpCompressionResponseCount     int                `json:"httpcompression_response_count,omitempty" bson:"httpcompression_response_count,omitempty"`
	//ResponseReceivedCount            int                `json:"response_received_count,omitempty" bson:"response_received_count,omitempty"`
	//RobotstxtResponseCount           int                `json:"robotstxt_response_count,omitempty" bson:"robotstxt_response_count,omitempty"`
	//RobotstxtResponseStatusCount404  int                `json:"robotstxt_response_status_count_404,omitempty" bson:"robotstxt_response_status_count_404,omitempty"`
	//DownloaderResponseStatusCount200 int                `json:"downloader_response_status_count_200,omitempty" bson:"downloader_response_status_count_200,omitempty"`
	ItemScrapedCount                 int                `json:"item_scraped_count,omitempty" bson:"item_scraped_count,omitempty"`
	//RequestDepthMax                  int                `json:"request_depth_max,omitempty" bson:"request_depth_max,omitempty"`
}

type User struct {
    Id           int    `json:"id"`
    Name         string `json:"name"`
    PasswordHash []byte `json:"password_hash"`
}
