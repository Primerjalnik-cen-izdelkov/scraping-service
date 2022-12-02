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
// @ID get_all_bots
// @Produce json
// @Success 200 {string} string
// @Router /bots [get]
/*
func (api *RestAPI) GetAllBots(c echo.Context) error {
	bots, err := api.bs.GetAllBots()
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, models.JSONError{
			Error: models.JSONErrorInfo{
				Code:    http.StatusInternalServerError,
				Message: err.Error()},
		}, "  ")
	}

	return c.JSONPretty(http.StatusOK, bots, "  ")
}
*/

// @Summary Send command to start scraping all the bots.
// @Description Scrape all bots.
// @ID scrape_all_bots
// @Produce json
// @Success 200 {string} data
// @Router /scrapeall [post]
/*
func (api *RestAPI) ScrapeAllBots(c echo.Context) error {
	err := api.bs.ScrapeAllBots()
	fmt.Println("scrapeall err:", err)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 500, Message: err.Error()}})
	}

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: "command send"}, "  ")
}
*/

/*
func (api *RestAPI) AllBotsStatus(c echo.Context) error {
	return nil
}
*/

// @Summary Send command to scrape specific bot.
// @Description Scrape bot.
// @ID scrape_bot
// @Produce json
// @Success 200 {string} data
// @Param bot_name path string true "Bot name"
// @Router /:bot_name/scrape [post]
/*
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
*/

/*
func (api *RestAPI) BotStatus(c echo.Context) error {
	return nil
}
*/

// @Summary List all files of the specific bot.
// @Description Get bot files.
// @ID bot_files
// @Produce json
// @Success 200 {string} data
// @Param bot_name path string true "Bot name"
// @Router /:bot_name/files [get]
/*
func (api *RestAPI) GetAllBotFiles(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")

	files, err := api.bs.GetBotFileNames(botName)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 404, Message: "Bot files err"}},
			"  ")
	}

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: files}, "  ")

}
*/

// @Summary List specific file of the specific bot.
// @Description Get specific bot file.
// @ID bot_file
// @Produce json
// @Success 200 {string} data
// @Param bot_name path string true "Bot name"
// @Param file path string true "File name"
// @Router /:bot_name/:file [get]
/*
func (api *RestAPI) GetBotFile(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")
	fileName := c.Param("file")

	return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("/data/%s/%s", botName, fileName))
}
*/

// @Summary
// @Description
// @ID
// @Produce json
// @Success 200 {string} data
// @Param bot_name path string true "Bot name"
// @Param file path string true "File name"
// @Router /:bot_name/:file/statistic [get]
/*
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
*/

/*
func (api *RestAPI) CmdStatus(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")

	api.bs.CmdStatus(botName)

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: "hehe"}, "  ")
}
*/

//////////////////////////////////////////////////////
// NOTE(miha): Here is the beggingig of our routing //
//////////////////////////////////////////////////////

// @Summary
// @Description
// @ID 
// @Produce
// @Success
// @Router
func (api *RestAPI) GetBots(c echo.Context) error {
	bots, err := api.bs.GetBots()
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, models.JSONError{
			Error: models.JSONErrorInfo{
				Code:    http.StatusInternalServerError,
				Message: err.Error()},
		}, "  ")
	}

    data := struct {
        Bots []*models.Bot
    }{
        Bots: bots,
    }

    return c.JSONPretty(http.StatusOK, models.JSONData{Data: data}, "  ")
}

func (api *RestAPI) GetFiles(c echo.Context) error {
	files, err := api.bs.GetFiles()
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 404, Message: "Bot files err"}},
			"  ")
	}

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: files}, "  ")
}

func (api *RestAPI) GetLogs(c echo.Context) error {
	return nil
}

func (api *RestAPI) GetCmds(c echo.Context) error {
	return nil
}

func (api *RestAPI) PostCmdScrape(c echo.Context) error {
	err := api.bs.PostCmdScrape()
	fmt.Println("scrapeall err:", err)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 500, Message: err.Error()}})
	}

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: "command send"}, "  ")
}

func (api *RestAPI) PostCmdStop(c echo.Context) error {
	return nil
}

func (api *RestAPI) PostCmdStatus(c echo.Context) error {
    return nil
}

func (api *RestAPI) GetBotCmds(c echo.Context) error {
    return nil
}

func (api *RestAPI) GetBotFiles(c echo.Context) error {
    return nil
}

// @Summary
// @Description
// @ID
// @Produce json
// @Success 200 {string} data
// @Param bot_name path string true "Bot name"
// @Router /:bot_name/statistic [get]
// TODO(miha): This function need to return bot logs from the scrapes
func (api *RestAPI) GetBotLogs(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")

	// QUERY params
    qm := map[string]string{}
	qm["querry"] = c.QueryParam("q")
	qm["projection"] = c.QueryParam("p")
	qm["sort"] = c.QueryParam("s")
	qm["fields"] = c.QueryParam("fields")
	qm["ignoreFields"] = c.QueryParam("ignore_fields")
	qm["timeLT"] = c.QueryParam("time.lt")
	qm["timeGT"] = c.QueryParam("time.gt")
	qm["timeSort"] = c.QueryParam("time.sort")
	qm["itemsScrapedLT"] = c.QueryParam("items_scraped.lt")
	qm["itemsScrapedGT"] = c.QueryParam("items_scraped.gt")
	qm["itemsScrapedSort"] = c.QueryParam("items_scraped.sort")

    fmt.Println("do we get timeLT:", qm["timeLT"], c.QueryParam("time.lt"))
    fmt.Println("do we get fields:", qm["fields"], c.QueryParam("fields"))

	logs, err := api.bs.GetBotLogs(botName, qm)
	if err != nil {
		fmt.Println(err)
		// TODO
		return c.JSONPretty(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 404, Message: "Bot files err"}},
			"  ")
	}

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: logs}, "  ")

}

func (api *RestAPI) GetBotLog(c echo.Context) error {
    return nil
}

func (api *RestAPI) GetBotFile(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")
    fileName := c.Param("file_name")

	return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("/data/%s/%s", botName, fileName))
}

func (api *RestAPI) PostBotCmdScrape(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")

    pid, err := api.bs.BotCmdScrape(botName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 500, Message: err.Error()}})
	}

    result := fmt.Sprintf("command send, bot pid: %d", pid)
    return c.JSONPretty(http.StatusOK, models.JSONData{Data: result}, "  ")
}

func (api *RestAPI) PostBotCmdStop(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")

    status := api.bs.BotCmdStop(botName)

    return c.JSONPretty(http.StatusOK, models.JSONData{Data: status}, "  ")
}

func (api *RestAPI) PostBotCmdStatus(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")

    status := api.bs.BotCmdStatus(botName)

    return c.JSONPretty(http.StatusOK, models.JSONData{Data: status}, "  ")
}
