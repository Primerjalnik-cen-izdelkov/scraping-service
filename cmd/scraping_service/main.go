package main

import (
	"fmt"
	"net/http"
	"os"
    "path"
    "context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

    "github.com/go-ping/ping"
    "github.com/rs/zerolog"
    "gopkg.in/Graylog2/go-gelf.v1/gelf"
    "github.com/rs/xid"
    "github.com/labstack/echo-contrib/prometheus"

	swaggerDocs "scraping_service/docs"
	"scraping_service/internal/api/rest"
	"scraping_service/internal/database"
	"scraping_service/internal/service"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Scraping service
// @description Service responsible for running scraping bots.

// @contact.name Miha
// @contact.email mf8974@student.uni-lj.si

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// Ping godoc
// @Summary Returns current service version.
// @Description Get current service version.
// @ID ping
// @Produce plain
// @Success 200 {string} string "service version"
// @Router /ping [get]
func Ping(c echo.Context) error {
	return c.String(http.StatusOK, os.Getenv("VERSION"))
}

var (
	HeaderCorrelationID = "X-Correlation-Id"
	KeyCorrelationID    = "correlationid"
)

func CorrelationIdMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			id := req.Header.Get(HeaderCorrelationID)
			if id == "" {
                id = xid.New().String()
			}
			c.Response().Header().Set(HeaderCorrelationID, id)
			newreq := req.WithContext(context.WithValue(req.Context(), KeyCorrelationID, id))
			c.SetRequest(newreq)
			return next(c)
		}
	}
}

// TODO(miha): Put more things into ENV variables
// TODO(miha): Create auth mechanism (check for echos framework website if they
// already have something) and use elephant postgres database to store
// credentials.
// DONE(miha): Add correlation IDs
// TODO(miha): Add healthchecks in docker (and kubernetes?)

// TODO(miha): Logging
//  - change gelfs source code so zerolog don't short write
//  - check if graylog is online with ping
//  - create new zerolog multiple logger for: stdout, file, graylog
    //  - add logs through service
//  - setup echo to use zerolog
//  - add correlation ID (rs/xid package)
    //  - add lumberjack package for rotating logs

// DONE(miha): Metrics
//  - use prometheus
//  - great tutorial on echos framework to combine with prometeus

// @Host localhost:1323
// TODO(miha): Put this into env variable
// @BasePath /v1
func main() {
    // NOTE(miha): Init logger
    graylogAddress := os.Getenv("GRAYLOG_ADDR") 
    // TODO(miha): ENV to set global level of graylog
    gelfWriter, err := gelf.NewWriter(graylogAddress)
    useGraylog := true

    // NOTE(miha): Create a logging file, get name from the ENV variable
    // LOG_FILE.
    logFileName := os.Getenv("LOG_FILE")
    logFile, err := os.OpenFile(
        path.Join("/logs", logFileName),
        os.O_APPEND|os.O_CREATE|os.O_WRONLY,
        0664,
    )
    if err != nil {
        fmt.Printf("Cannot create file %s\n", logFileName)
        fmt.Println(err)
    }
    defer logFile.Close()

    // NOTE(miha): Try to ping graylog address. If it is not accessible, don't
    // add gelfWriter as a zerolog source.
    pinger, err := ping.NewPinger(graylogAddress)
    if err != nil {
        useGraylog = false
        fmt.Printf("Logging service graylog at address %s is not avaiable.\n", graylogAddress)
    }

    // NOTE(miha): Add sources to zerolog, ie. add or ignore gelfWriter.
    var multi zerolog.LevelWriter
    if useGraylog {
        multi = zerolog.MultiLevelWriter(os.Stdout, logFile, gelfWriter)
    } else {
        multi = zerolog.MultiLevelWriter(os.Stdout, logFile)
    }

    // NOTE(miha): Create our logger
    logger := zerolog.New(multi).With().Timestamp().Caller().Logger()

    _, _, _, _, _ = gelfWriter, pinger, logFile, multi, logger

	e := echo.New()
	//e.Use(middleware.Logger())
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"http://localhost:5173"},
        // TODO(miha): What are allowHeaders? dig into this...
        AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
    }))
    // NOTE(miha): Setup echo to use zerolog.
    e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
        LogURI:    true,
        LogStatus: true,
        LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
            logger.Info().
                Str("URI", v.URI).
                Int("status", v.Status).
                Msg("request")

            return nil
        },
    }))
    e.Use(CorrelationIdMiddleware())

    // NOTE(miha): Setup prometheus metrics.
    p := prometheus.NewPrometheus("echo", nil)
    p.Use(e)

    // NOTE(miha): Set swagger's version of the API.
    swaggerDocs.SwaggerInfo.Version = fmt.Sprintf("%s", os.Getenv("VERSION"))

	fmt.Println("Scraping service started, running on version: ", os.Getenv("VERSION"))

	mongoDB, err := database.CreateDatabase("MongoDB")
	if err != nil {
		fmt.Println("mongoErr: ", err)
	}
	err = mongoDB.Ping()
	if err != nil {
		fmt.Println("mongoErr ping: ", err)
	}
	bs := service.CreateBotService(mongoDB)
	rest := rest.CreateRestAPI(bs, &logger)

	/* NOTE(miha): How to check if Boter interface is implemented.
	var pp interface{} = bs
	if _, ok := pp.(models.Boter); ok {
		fmt.Println("bs implements boter")
	} else {
		fmt.Println("bs don't implements boter")
	}
	*/
	//e.Static("data", "./scrapy_grocery_stores")

	// /v0/bots/{mercator}/cmd/scrape
	// /v0/bots/{mercator}/files/
	// /v0/bots/{mercator}/files/mercator_T324
	// /v0/bots/{mercator}/stats
	// /v0/bots/scrape
	// /v0/bots/files
	// /v0/bots/stats

	fs := http.FileServer(http.Dir("./scrapy_grocery_stores/data"))
	e.GET("/data/*", echo.WrapHandler(http.StripPrefix("/data/", fs)))

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/ping", Ping)

	versionGroup := e.Group("/v1")

	// /bots/
	botsGroup := versionGroup.Group("/bots")
	{
		botsGroup.GET("", rest.GetBots)
		botsGroup.GET("/files", rest.GetFiles)
		botsGroup.GET("/logs", rest.GetLogs)
		botsGroup.GET("/cmd", rest.GetCmds)
		botsGroup.POST("/cmd/scrape", rest.PostCmdScrape)
		botsGroup.POST("/cmd/stop", rest.PostCmdStop)
		//botsGroup.POST("/cmd/status", rest.PostCmdStatus)

		// /bots/{bot_name}/
		botNameGroup := botsGroup.Group("/:bot_name")
		{
			botNameGroup.GET("/cmd", rest.GetBotCmds)
			botNameGroup.GET("/files", rest.GetBotFiles)
			botNameGroup.GET("/logs", rest.GetBotLogs)
            //botNameGroup.GET("/logs/:file_name", rest.GetBotLog)
			botNameGroup.GET("/files/:file_name", rest.GetBotFile)
			botNameGroup.POST("/cmd/scrape", rest.PostBotCmdScrape)
			botNameGroup.POST("/cmd/stop", rest.PostBotCmdStop)
			botNameGroup.POST("/cmd/status", rest.PostBotCmdStatus)
		}
	}
	/*
		// botName := c.Param("bot_name")
		e.POST("/scrapeall", rest.ScrapeAllBots)
		e.GET("/statusall", rest.AllBotsStatus)
		botGroup := e.Group(":bot_name")
		{
			botGroup.GET("/filenames", rest.GetBotFiles)
			botGroup.GET("/files", rest.GetBotFiles)
			botGroup.GET("/:file", rest.GetBotFile)
			botGroup.POST("/scrape", rest.ScrapeBot)
			botGroup.GET("/status", rest.BotStatus)
			botGroup.GET("/statistic", rest.GetBotStats)
			botGroup.GET("/:file/statistic", rest.GetFileStats)
			// TODO(miha): Get some health-check for the bots - ie. if the bot run was succesfull.

			cmdGroup := botGroup.Group("/cmd")
			{
				cmdGroup.GET("/status", rest.CmdStatus)
			}
		}
	*/

	e.Logger.Fatal(e.Start(":1323"))
}
