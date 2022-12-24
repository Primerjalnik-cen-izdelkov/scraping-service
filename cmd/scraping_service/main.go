package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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

// TODO(miha): Put more things into ENV variables
// TODO(miha): Create auth mechanism (check for echos framework website if they
// already have something) and use elephant postgres database to store
// credentials.

// TODO(miha): Logging
//  - change gelfs source code so zerolog don't short write
//  - check if graylog is online with ping
//  - create new zerolog multiple logger for: stdout, file, graylog
//  - add logs through service
//  - setup echo to use zerolog

// TODO(miha): Metrics
//  - use prometheus
//  - great tutorial on echos framework to combine with prometeus

// @Host localhost:1323
// TODO(miha): Put this into env variable
// @BasePath /v1
func main() {
	e := echo.New()
	e.Use(middleware.Logger())
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"http://localhost:5173"},
        // TODO(miha): What are allowHeaders? dig into this...
        AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
    }))

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
