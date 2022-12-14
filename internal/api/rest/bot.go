package rest

import (
	"fmt"
	"net/http"
	"scraping_service/internal/service"
	"scraping_service/pkg/models"

	"github.com/labstack/echo/v4"
    "github.com/rs/zerolog"
)

type RestAPI struct {
	bs *service.BotService
    logger *zerolog.Logger
}

func CreateRestAPI(service *service.BotService) *RestAPI {
    res := &RestAPI{bs: service}
    return res
}

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

/*
func (api *RestAPI) GetBotFile(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")
	fileName := c.Param("file")

	return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("/data/%s/%s", botName, fileName))
}
*/

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

func (api *RestAPI) Login(c echo.Context) error {
    type User struct {
        Name     string `json:"name"`
        Password string `json:"password"`
    }

    var u User

    // NOTE(miha): Bind/parse json body into variable 'u' of type 'User'.
    if err := c.Bind(&u); err != nil {
        return c.String(http.StatusBadRequest, "Error with parsing body.")
    }

    // TODO(miha): Bot service must have function Login that calls postgres database.
    /*
    var users []*struct {
        Id  int
        Name string
        PasswordHash []byte
    }
    err = api.bs.AuthDb.Query("SELECT * FROM users").Rows(&users)
    if err != nil {
    }
    for _, user := range users {
        fmt.Printf("id: %d, name: %s, hash: %s", user.Id, user.Name, user.PasswordHash)
    }
    */

    jwt, err := api.bs.Login(u.Name, u.Password)
    if err != nil {

    }
    _ = jwt

    fmt.Println("user:", u, err)

    //return c.String(http.StatusOK, fmt.Sprintf("we got %s %s", u.Name, u.Password))
    	return c.JSON(http.StatusOK, echo.Map{
		"token": jwt,
	})
}


// TODO(miha): Accept query params
// @Summary Returns array of all avaiable bots.
// @Description We are returning list of type '[]*models.Bot' for all the bots we have created.
// @ID get_bots 
// @Tags bots
// @Produce json
// @Success 200 {object} models.JSONData
// @Failure 500 {object} models.JSONError
// @Router /bots [get]
func (api *RestAPI) GetBots(c echo.Context) error {
    qp := c.QueryParams()

	bots, err := api.bs.GetBots(qp)
	if err != nil {
        api.bs.Logger.Error().Err(err).Msg(err.Error())
		return c.JSONPretty(http.StatusInternalServerError, models.JSONError{
			Error: models.JSONErrorInfo{
				Code:    http.StatusInternalServerError,
				Message: err.Error()},
            }, "  ")
	}

    data := struct {
        Bots []*models.Bot `json:"bots"`
    }{
        Bots: bots,
    }

    return c.JSONPretty(http.StatusOK, models.JSONData{Data: data}, "  ")
}

// @Summary Returns array of the filenames that contain scraped data.
// @Description We are returning array '[]models.File' that contains all the filenames with scraped data. All queries are on the day basis - meaning we only care about year, month, and day.
// @ID get_files
// @Tags bots
// @Produce json
// @Success 200 {object} models.JSONData
// @Failure 500 {object} models.JSONError
// @Param id      query []string false "return files where given date is the same as the file id"
// @Param date    query []string false "return files where given date is the same as the file date"
// @Param date.gt query string   false "return files where given date is greater than the file date"
// @Param date.lt query string   false "return files where given date is less than the file date"
// @Param name    query []string false "return files where given name is the same as the bot name"
// @Param limit   query int      false "define how many results will be returned"
// @Param sort    query []string false "define on which field we sort result (prefix with '-' for reversed sort)"
// @Param field   query []string false "select which fields to include (field=field_name) or exclude (field=-field_name) from the query. Note that you can't include and exclude fields at the same time."
// @Router /bots/files [get]
func (api *RestAPI) GetFiles(c echo.Context) error {
    qp := c.QueryParams()
    
	files, err := api.bs.GetFiles(qp)
	if err != nil {
        api.bs.Logger.Error().Err(err).Msg(err.Error())
		return c.JSONPretty(http.StatusInternalServerError,
            models.JSONError{
                Error: models.JSONErrorInfo{
                    Code: http.StatusInternalServerError, 
                    Message: err.Error()},
                }, "  ")
	}

    data := struct {
        Files []models.File `json:"files"`
    }{
        Files: files,
    }

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: data}, "  ")
}

// @Summary Retruns array of the logs that contain scraped data information (items scraped,...) 
// @Description We are returning array '[]models.FileLog' that contains all the information of runs of the scraped data. All queries are on the day basis - meaning we only care about year, month, and day.
// @ID get_logs
// @Tags bots
// @Produce json
// @Success 200 {object} models.JSONData
// @Failure 500 {object} models.JSONError
// @Param id                query []string false "return files where given id is the same as the file id"
// @Param start_time        query []string false "return files where given date is the same as the file date"
// @Param start_time.gt     query int      false "return files where given date is greater than the file date"
// @Param start_time.lt     query int      false "return files where given date is less than the file date"
// @Param request_count     query []int    false "return files where given request count is the same as the file request count"
// @Param request_count.gt  query int      false "return files where given request count is greater than the file request count"
// @Param request_count.lt  query int      false "return files where given request count is less than the file request count"
// @Param response_count    query []int    false "return files where given response count is the same as the file response count"
// @Param response_count.gt query int      false "return files where given response count is greater than the file response count"
// @Param response_count.lt query int      false "return files where given response count is less than the file response count"
// @Param 404               query []int    false "return files where given 404 count is the same as the file 404 count"
// @Param 404.gt            query int      false "return files where given 404 count is greater than the file 404 count"
// @Param 404.lt            query int      false "return files where given 404 count is less than the file 404 count"
// @Param item_scraped      query []int    false "return files where given item scraped is the same as the file item scraped"
// @Param item_scraped.gt   query int      false "return files where given item scraped is greater than the file item scraped"
// @Param item_scraped.lt   query int      false "return files where given item scraped is less than the file item scraped"
// @Param name              query []string false "return files where given name is the same as the bot name"
// @Param limit             query int      false "define how many results will be returned"
// @Param sort              query []string false "define on which field we sort result (prefix with '-' for reversed sort)"
// @Param field             query []string false "select which fields to include (field=field_name) or exclude (field=-field_name) from the query. Note that you can't include and exclude fields at the same time."
// @Router /bots/logs [get]
func (api *RestAPI) GetLogs(c echo.Context) error {
    qp := c.QueryParams()

	logs, err := api.bs.GetLogs(qp)
	if err != nil {
        api.bs.Logger.Error().Err(err).Msg(err.Error())
		return c.JSONPretty(http.StatusInternalServerError,
            models.JSONError{
                Error: models.JSONErrorInfo{
                    Code: http.StatusInternalServerError, 
                    Message: err.Error()},
                }, "  ")
	}

    data := struct {
        Logs []models.FileLog `json:"scrape_logs"`
    }{
        Logs: logs,
    }

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: data}, "  ")
}

// @Summary Return all avaiable commands.
// @Description Returns all avaiable commands that can be run on the /bots/cmd/{cmd_name} endpoint with POST request.
// @ID get_cmd
// @Tags bots
// @Produce json
// @Success 200 {object} models.JSONData
// @Failure 500 {object} models.JSONError
// @Router /bots/cmd [get]
func (api *RestAPI) GetCmds(c echo.Context) error {
    type Command struct {
        Name string `json:"name"`
        Description string `json:"description"`
    }

    cmds := []Command{
        {"scrape", "start scraping all avaiable bots"},
        {"stop", "stop scraping all running bots"},
    }

    data := struct {
        Commands []Command `json:"commands"`
    }{
        Commands: cmds,
    }

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: data}, "  ")
}

// @Summary Send command to start scraping all the bots
// @ID post_cmd_scrape
// @Tags bots
// @Produce json
// @Success 200 {object} models.JSONData
// @Failure 500 {object} models.JSONError
// @Router /cmd/scrape [post]
func (api *RestAPI) PostCmdScrape(c echo.Context) error {
	err := api.bs.PostCmdScrape()
	if err != nil {
        api.bs.Logger.Error().Err(err).Msg(err.Error())
		return c.JSONPretty(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 500, Message: err.Error()}}, "  ")
	}

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: "command send"}, "  ")
}

// @Summary Send command to stop all the bots
// @ID post_cmd_stop
// @Tags bots
// @Produce json
// @Success 200 {object} models.JSONData
// @Failure 500 {object} models.JSONError
// @Router /cmd/stop [post]
func (api *RestAPI) PostCmdStop(c echo.Context) error {
	err := api.bs.PostCmdStop()
	if err != nil {
        api.bs.Logger.Error().Err(err).Msg(err.Error())
		return c.JSONPretty(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 500, Message: err.Error()}}, "  ")
	}

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: "command send"}, "  ")
}

// @Summary Get bot commands for the given bot
// @ID get_bot_cmd
// @Tags bots
// @Produce json
// @Success 200 {object} models.JSONData
// @Failure 500 {object} models.JSONError
// @Param botName path string true "define for which bot we get commands"
// @Router /bots/{bot_name}/cmds [get]
func (api *RestAPI) GetBotCmds(c echo.Context) error {
    // PATH params
	botName := c.Param("bot_name")

    type Command struct {
        Name string `json:"name"`
        Description string `json:"description"`
    }

    cmds := []Command{
        {"scrape", "start scraping bot: " + botName},
        {"stop", "stop scraping bot: " + botName},
    }

    data := struct {
        Commands []Command `json:"commands"`
    }{
        Commands: cmds,
    }

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: data}, "  ")
}

// @Summary Get bot files for the given bot
// @ID get_bot_files
// @Tags bots
// @Produce json
// @Success 200 {object} models.JSONData
// @Failure 500 {object} models.JSONError
// @Param botName path string true "define for which bot we get files"
// @Router /bots/{bot_name}/files [get]
func (api *RestAPI) GetBotFiles(c echo.Context) error {
    // PATH params
	botName := c.Param("bot_name")

    c.QueryParams().Set("name", botName)
    return api.GetFiles(c)
}

// @Summary Get bot logs for the given bot
// @ID get_bot_logs
// @Tags bots
// @Produce json
// @Success 200 {object} models.JSONData
// @Failure 500 {object} models.JSONError
// @Param botName path string true "define for which bot we get logs"
// @Router /bots/{bot_name}/logs [get]
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

    //fmt.Println("do we get timeLT:", qm["timeLT"], c.QueryParam("time.lt"))
    //fmt.Println("do we get fields:", qm["fields"], c.QueryParam("fields"))

	logs, err := api.bs.GetBotLogs(botName, qm)
	if err != nil {
        api.bs.Logger.Error().Err(err).Msg(err.Error())
		return c.JSONPretty(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 500, Message: err.Error()}}, "  ")
	}

	return c.JSONPretty(http.StatusOK, models.JSONData{Data: logs}, "  ")

}

func (api *RestAPI) GetBotLog(c echo.Context) error {
    return nil
}

// @Summary Get bot file for the given bot
// @ID get_bot_files_file
// @Tags bots
// @Produce json
// @Success 200 {object} models.JSONData
// @Failure 500 {object} models.JSONError
// @Param botName path string true "define for which bot we get logs"
// @Param fileName path string true "define which file we want to get"
// @Router /bots/{bot_name}/files/{file_name} [get]
func (api *RestAPI) GetBotFile(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")
    fileName := c.Param("file_name")

	return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("/scraping-service/data/%s/%s", botName, fileName))
}

// @Summary Start scraping given bot
// @ID get_bot_cmd_scrape
// @Tags bots
// @Produce json
// @Success 200 {object} models.JSONData
// @Failure 500 {object} models.JSONError
// @Param botName path string true "define for which bot we get logs"
// @Router /bots/{bot_name}/cmd/scrape [post]
func (api *RestAPI) PostBotCmdScrape(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")

    err := api.bs.PostBotCmdScrape(botName)
	if err != nil {
        api.bs.Logger.Error().Err(err).Msg(err.Error())
		return c.JSON(http.StatusInternalServerError,
			models.JSONError{Error: models.JSONErrorInfo{Code: 500, Message: err.Error()}})
	}

    result := fmt.Sprintf("command send")
    return c.JSONPretty(http.StatusOK, models.JSONData{Data: result}, "  ")
}

// @Summary Stop scraping given bot
// @ID get_bot_cmd_stop
// @Tags bots
// @Produce json
// @Success 200 {object} models.JSONData
// @Failure 500 {object} models.JSONError
// @Param botName path string true "define for which bot we get logs"
// @Router /bots/{bot_name}/cmd/stop [post]
func (api *RestAPI) PostBotCmdStop(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")

    status := api.bs.BotCmdStop(botName)

    return c.JSONPretty(http.StatusOK, models.JSONData{Data: status}, "  ")
}

// @Summary Get status of scraping for the given bot
// @ID get_bot_cmd_status
// @Tags bots
// @Produce json
// @Success 200 {object} models.JSONData
// @Failure 500 {object} models.JSONError
// @Param botName path string true "define for which bot we get logs"
// @Router /bots/{bot_name}/cmd/status [post]
func (api *RestAPI) PostBotCmdStatus(c echo.Context) error {
	// PATH params
	botName := c.Param("bot_name")

    status := api.bs.BotCmdStatus(botName)

    return c.JSONPretty(http.StatusOK, models.JSONData{Data: status}, "  ")
}
