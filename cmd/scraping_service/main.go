package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	_ "scraping_service/docs"
	"scraping_service/internal/api/rest"
	"scraping_service/internal/database"
	"scraping_service/internal/service"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Scraping service
// @version 1.0.1
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
	return c.String(http.StatusOK, "scraping_service, 0.0.1")
}

// @host localhost:1323
// @BasePath /
func main() {
	e := echo.New()

	mongoDB, err := database.CreateDatabase("MongoDB")
	if err != nil {
		fmt.Println("mongoErr: ", err)
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

	fs := http.FileServer(http.Dir("./scrapy_grocery_stores/data"))
	e.GET("/data/*", echo.WrapHandler(http.StripPrefix("/data/", fs)))

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/ping", Ping)
	e.GET("/bots", rest.GetAllBots)
	// botName := c.Param("bot_name")
	e.POST("/scrapeall", rest.ScrapeAllBots)
	botName := e.Group(":bot_name")
	{
		botName.GET("/filenames", rest.GetBotFiles)
		botName.GET("/files", rest.GetBotFiles)
		botName.GET("/:file", rest.GetBotFile)
		botName.POST("/scrape", rest.ScrapeBot)
		botName.GET("/statistic", rest.GetBotStats)
		botName.GET("/:file/statistic", rest.GetFileStats)
		// TODO(miha): Get some health-check for the bots - ie. if the bot run was succesfull.
	}

	e.Logger.Fatal(e.Start(":1323"))
}
