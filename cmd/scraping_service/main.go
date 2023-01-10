package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/golang-jwt/jwt/v4"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
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
    fmt.Println("HMMMMMM")
    log.Info().Str("ping", "just log works?")
	return c.String(http.StatusOK, os.Getenv("VERSION"))
}

type jwtClaims struct {
    Name string `json:"name"`
    jwt.RegisteredClaims
}

// TODO(miha): Put more things into ENV variables
// TODO(miha): Create auth mechanism (check for echos framework website if they
// already have something) and use elephant postgres database to store
// credentials.
// DONE(miha): Add correlation IDs
// TODO(miha): Add healthchecks in kubernetes

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
    //multi := zerolog.MultiLevelWriter(os.Stdout, os.Stderr)
    // NOTE(miha): Create our logger
    logger := zerolog.New(os.Stdout)//.With().Timestamp().Caller().Logger()

	e := echo.New()
	e.Use(middleware.Recover())
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"*"}, //[]string{os.Getenv("CORS_URI")},
        // TODO(miha): What are allowHeaders? dig into this...
        AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
    }))
    // NOTE(miha): Setup echo to use zerolog.
    /*
DefaultLoggerConfig = LoggerConfig{
  Skipper: DefaultSkipper,
  Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
    `"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
    `"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
    `,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
  CustomTimeFormat: "2006-01-02 15:04:05.00000",
}
    */
    e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
        Generator: func() string {
            return xid.New().String()
        },
        TargetHeader: "X-Request-Id",
    }))
    e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
        LogURI:           true,
        LogStatus:        true,
        LogRemoteIP:      true,
        LogHost:          true,
        LogMethod:        true,
        LogUserAgent:     true,
        LogLatency:       true,
        LogRequestID:     true,
        LogError:         true,
        LogProtocol:      true,
        LogURIPath:       true,
        LogRoutePath:     true,
        LogReferer:       true,
        LogContentLength: true,
        LogResponseSize:  true,
        LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
            logger.Info().
                Str("id", v.RequestID).
                Str("remote_ip", v.RemoteIP).
                Str("host", v.Host).
                Str("method", v.Method).
                Str("uri", v.URI).
                Str("user_agent", v.UserAgent).
                Int("status", v.Status).
                Str("latency", v.Latency.String()).
                Msg("logging middleware")

            return nil
        },
    }))

    // NOTE(miha): Setup prometheus metrics.
    p := prometheus.NewPrometheus("echo", nil)
    p.Use(e)

    // NOTE(miha): Set swagger's version of the API.
    swaggerDocs.SwaggerInfo.Version = fmt.Sprintf("%s", os.Getenv("VERSION"))

	fmt.Println("Scraping service started, running on version: ", os.Getenv("VERSION"))

	mongoDB, err := database.CreateDatabase("MongoDB", &logger)
	if err != nil {
		fmt.Println("mongoErr: ", err)
	}
	err = mongoDB.Ping()
	if err != nil {
		fmt.Println("mongoErr ping: ", err)
	}

    authDB, err := database.CreateAuthDatabase("AuthPostgresDB")
	if err != nil {
		fmt.Println("postgres auth err: ", err)
	}
    logger = zerolog.New(os.Stdout)//.With().Timestamp().Caller().Logger()
	bs := service.CreateBotService(mongoDB, &logger, authDB)
    bs.Logger.Info().Str("jupi", "we are here")
    log.Info().Str("ahdfsa", "just log works?")
    fmt.Println("Ce to dela....")

	rest := rest.CreateRestAPI(bs)

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
    versionGroup.POST("/login", rest.Login)

    jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtClaims)
		},
        ErrorHandler: func(c echo.Context, err error) error {
            fmt.Println("error handler:", err)
            return err
        },
		SigningKey: []byte("secret"),
	}
    _ = jwtConfig
    //versionGroup.Use(echojwt.WithConfig(jwtConfig))

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
