package rest

import (
	"fmt"
	"net/http"
	"scraping_service/internal/service"
	"scraping_service/pkg/models"

	"github.com/labstack/echo/v4"
)

type RestAPI struct {
	bs *service.BotService
}

func CreateRestAPI(service *service.BotService) *RestAPI {
	return &RestAPI{bs: service}
}

// @Summary Returns list of all bots.
// @Description Get list of all bots.
// @ID list_all_bots
// @Produce json
// @Success 200 {string} string
// @Router /bots [get]
func (api *RestAPI) GetAllBots(c echo.Context) error {
	bots, err := api.bs.GetAllBots()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSONPretty(http.StatusOK, bots, "  ")
}

// @Summary Send commands to start scraping all the bots.
// @Description Scrape all bots.
// @ID scrape_all_bots
// @Produce json
// @Success 200 {string} data
// @Router /scrapeall [post]
func (api *RestAPI) ScrapeAllBots(c echo.Context) error {
	err := api.bs.ScrapeAllBots()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 500, Message: err.Error()}})
	}

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: "command send"}, "  ")
}

// @Summary Send command to scrape specific bot.
// @Description Scrape bot.
// @ID scrape_bot
// @Produce json
// @Success 200 {string} data
// @Param bot_name path string true "Bot name"
// @Router /:bot_name/scrape [post]
func (api *RestAPI) ScrapeBot(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")

	err := api.bs.ScrapeBot(botName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 500, Message: err.Error()}})
	}

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: "command send"}, "  ")
}

// @Summary List all files of the specific bot.
// @Description Get bot files.
// @ID bot_files
// @Produce json
// @Success 200 {string} data
// @Param bot_name path string true "Bot name"
// @Router /:bot_name/files [get]
func (api *RestAPI) GetBotFiles(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")

	files, err := api.bs.GetBotFileNames(botName)
	if err != nil {
		fmt.Println(err)
		// TODO
		return c.JSONPretty(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 404, Message: "Bot files err"}},
			"  ")
	}

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: files}, "  ")

}

// @Summary List specific file of the specific bot.
// @Description Get specific bot file.
// @ID bot_file
// @Produce json
// @Success 200 {string} data
// @Param bot_name path string true "Bot name"
// @Param file path string true "File name"
// @Router /:bot_name/:file [get]
func (api *RestAPI) GetBotFile(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")
	fileName := c.Param("file")

	return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("/data/%s/%s", botName, fileName))
}

// @Summary
// @Description
// @ID
// @Produce json
// @Success 200 {string} data
// @Param bot_name path string true "Bot name"
// @Router /:bot_name/statistic [get]
func (api *RestAPI) GetBotStats(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")

	// QUERY params
	// TODO(miha):
	// 		- construct string for the querry (we should pass []map[string]string param instead)
	querry := c.QueryParam("q")
	projection := c.QueryParam("p")
	sort := c.QueryParam("s")
	fields := c.QueryParam("fields")
	ignoreFields := c.QueryParam("ignore_fields")
	timeLT := c.QueryParam("time.lt")
	timeGT := c.QueryParam("time.gt")
	timeSort := c.QueryParam("time.sort")
	itemsScrapedLT := c.QueryParam("items_scraped.lt")
	itemsScrapedGT := c.QueryParam("items_scraped.gt")
	itemsScrapedSort := c.QueryParam("items_scraped.sort")

	fmt.Println("QuerryParas: ", querry, projection, sort, fields, ignoreFields, timeLT, timeGT, timeSort, itemsScrapedLT, itemsScrapedGT, itemsScrapedSort)
	urlQuerry := ""
	files, err := api.bs.GetBotStats(botName, urlQuerry)
	if err != nil {
		fmt.Println(err)
		// TODO
		return c.JSONPretty(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 404, Message: "Bot files err"}},
			"  ")
	}

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: files}, "  ")

}

// @Summary
// @Description
// @ID
// @Produce json
// @Success 200 {string} data
// @Param bot_name path string true "Bot name"
// @Param file path string true "File name"
// @Router /:bot_name/:file/statistic [get]
func (api *RestAPI) GetFileStats(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")
	fileName := c.Param("file")

	// QUERY params
	fields := c.QueryParam("fields")
	ignoreFields := c.QueryParam("ignore_fields")

	fmt.Println("QuerryParas: ", fields, ignoreFields)

	file, err := api.bs.GetFileStats(botName, fileName)
	if err != nil {
		fmt.Println(err)
		// TODO
		return c.JSONPretty(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 404, Message: "Bot files err"}},
			"  ")
	}

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: file}, "  ")
}
