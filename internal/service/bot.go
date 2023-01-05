package service

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	//"path"
    "net/url"
	"runtime"
	"scraping_service/internal/database"
	"scraping_service/pkg/common"
	"scraping_service/pkg/models"
	"strings"
    "github.com/rs/zerolog"
)

// Errors that can occur in the bot service.
var (
	ErrDirectoryNotFound = errors.New("Bot directory not found")
	ErrDirectoryIsEmpty  = errors.New("Bot directory is empty")
	ErrBotBadFilename    = errors.New("Bot filename does not contain a dot")
	ErrFilesEmpty        = errors.New("Returned files are empty")
    ErrBotAlreadyRunning = errors.New("Bot is alredy running")
    ErrBotIsNotRunning   = errors.New("Bot is not running")
    ErrNoPython          = errors.New("Python couldn't be found on the system")
    ErrCantStartBot      = errors.New("Bot process can't be started")
    ErrCantKillBot       = errors.New("Bot process can't be killed")
)

// 'BotService' struct holds information about database connection and running
// bot processes.
//
// If we can find value in the 'botPID' map for given bot name, the bot is
// currently running - value contains information about the process.
type BotService struct {
    // Database connection
	db     *database.Database
    // Map for holding bot processes if any exists.
	botPID map[string]*exec.Cmd
    Logger *zerolog.Logger
}

// Returns new 'BotService' struct with the given database 'd' connection.
func CreateBotService(d *database.Database, logger *zerolog.Logger) *BotService {
    return &BotService{db: d, botPID: make(map[string]*exec.Cmd), Logger: logger}
}

// Function 'BotNames' returns array '[]string' of all avaiable bot names.
//
// Our bots are saved in the "./scrapy_grocery_store/scrapy_grocery_stores/spiders"
// directory (we use python framework Scrapy for scraping websites). Bots are
// named after grocery store that they scrape and have extension ".py". So to
// gather all the bots we use 'os' package to read all the filenames in mentioned
// directory (we skip default python files that starts with "_").
//
// Errors:
//  - ErrDirectoryNotFound: Given directory can't be opened.
//  - ErrDirectoryIsEmpty: Can't read the filenames in the given directory.
func (bs *BotService) BotNames() ([]string, error) {
    // NOTE(miha): Open directory with python bots.
	dir, err := os.Open("./scrapy_grocery_stores/scrapy_grocery_stores/spiders")
	if err != nil {
        bs.Logger.Error().Err(err).Msg(err.Error())
		return nil, ErrDirectoryNotFound
	}

    // NOTE(miha): Read all the filenames in the directory 'dir'. 
	files, err := dir.Readdirnames(0)
	if err != nil {
        bs.Logger.Error().Err(err).Msg(err.Error())
		return nil, ErrDirectoryIsEmpty
	}

	bots := []string{}

    // NOTE(miha): Iterate over files and ignore those which starts with "_",
    // also bot names in 'bots' don't have ".py" extension so trim it.
	for _, file := range files {
		// NOTE(miha): Parse spiders directory (ignore python files)
		if file[0] == '_' {
			continue
		}
		index := strings.Index(file, ".")
		if index == -1 {
			return nil, ErrBotBadFilename
		}
		trimmedFile := file[0:strings.Index(file, ".")]

		bots = append(bots, trimmedFile)
	}

	return bots, nil
}

// Returns array of type '[]*models.Bot' which contains information for all the
// bots.
//
// Params:
//  - qp: Query parameters that are parsed from URL.
//
// Errors:
//  - errors from 'internal.service.bot.go.BotNames()' 
//  - errors from 'internal.database.model.go'
func (bs *BotService) GetBots(qp url.Values) ([]*models.Bot, error) {
    // NOTE(miha): Get all bots.
    botNames, err := bs.BotNames()
    if err != nil {
        bs.Logger.Error().Err(err).Msg(err.Error())
        return nil, err
    }

	bots := []*models.Bot{}

    // NOTE(miha): Iterate bot names and get information about last scrape and
    // logs count for each bot - we make retrive information from database. We
    // also check and update the bot's status - if the bot is currently
    // scraping.
	for _, name := range botNames {
        // NOTE(miha): Get bot info
        // TODO(miha): We cannot have qp for just one bot (can't filter by name with qp)...
        bot, err := bs.db.GetBot(name, qp)
        if err != nil {
            bs.Logger.Error().Err(err).Msg(err.Error())
            return nil, err
        }

        // NOTE(miha): Check bot's status
        status := bs.BotCmdStatus(name)
        bot.Status = status

        bots = append(bots, bot)
	}

	return bots, nil
}

// Function 'GetFiles' returns array '[]models.File' of all the scraped files. They
// are saved in the next format '{bot_name}_{time_stamp}.csv'.
//
// Params:
//  - qp: Query parameters that are parsed from URL.
//
// Errors:
//  - errors from 'internal.database.model.go'
func (bs *BotService) GetFiles(qp url.Values) ([]models.File, error) {
    files, err := bs.db.GetFiles(qp)
    return files, err
}

// Return array '[]models.FileLog' that contains information logs on the bot runs. 
//
// Params:
//  - qp: Query parameters that are parsed from URL.
//
// Errors:
//  - errors from 'internal.database.model.go'
func (bs *BotService) GetLogs(qp url.Values) ([]models.FileLog, error) {
    logs, err := bs.db.GetLogs(qp)
    return logs, err
}

// Send command that start scraping to all bots. If bot is already running skip
// starting it.
//
// Errors:
//  - errors from 'internal.service.bot.go.BotNames()'
//  - ErrNoPython: System don't have python installed
//  - ErrCantStartBot: System couldn't start the bot process
func (bs *BotService) PostCmdScrape() error {
    // NOTE(miha): Get all the bot names.
	bots, err := bs.BotNames()
	if err != nil {
        bs.Logger.Error().Err(err).Msg(err.Error())
		return err
	}

	// NOTE(miha): Check if python is installed on the system.
	_, err = exec.LookPath("python")
	if err != nil {
        bs.Logger.Error().Err(err).Msg(err.Error())
		return ErrNoPython
	}

    // NOTE(miha): For each bot run a crawl command - starts scraping.
	for _, bot := range bots {
		cmd := common.MultiOSCommand("scrapy", "crawl", bot)
		cmd.Dir = "./scrapy_grocery_stores"

		// TODO(miha): Set stdout and stderr files, so we can tail them to see
		// status.
		// NOTE(miha): Everything we write in python is send to stderr
        /*
		stdout, err := os.OpenFile(path.Join(cmd.Dir, "/data", "/"+bot, "test_stdout.txt"),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		cmd.Stdout = stdout
		if err != nil {
		}

		stderr, err := os.OpenFile(path.Join(cmd.Dir, "/data", "/"+bot, "test_stderr.txt"),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		cmd.Stdout = stderr
		if err != nil {
		}
        */

        // NOTE(miha): Check if bot is alredy running. If it is skip it.
		if _, ok := bs.botPID[bot]; ok {
            continue
		}

		bs.botPID[bot] = cmd

		err = cmd.Start()
		if err != nil {
            bs.Logger.Error().Err(err).Msg(err.Error())
			return ErrCantStartBot
		}

		// NOTE(miha): This go-routine triggers when process finished - can
		// return errors!
        // TODO(miha): What to do about errors?
		go func(bot string) {
			err := cmd.Wait()
			if err != nil {
                bs.Logger.Error().Err(err).Msg(err.Error())
			}

            err = bs.db.UpdateBot(bot)
			if err != nil {
                bs.Logger.Error().Err(err).Msg(err.Error())
			}
            delete(bs.botPID, bot)
		}(bot)
	}

	return nil
}

// Send command that stop all the bots.
//
// Errors:
//  - ErrCantKillBot: System couldn't kill the bot process
func (bs *BotService) PostCmdStop() error {
    for _, v := range bs.botPID {
        err := v.Process.Kill()
        if err != nil {
            bs.Logger.Error().Err(err).Msg(err.Error())
            return ErrCantKillBot
        }
    }
    return nil
}

func (bs *BotService) GetBotFiles() error {
    return nil
}

func (bs *BotService) GetBotLogs(botName string, qm map[string]string) ([]models.FileLog, error) {
	logs, err := bs.db.GetBotLogs(botName, qm)
	if err != nil {
        bs.Logger.Error().Err(err).Msg(err.Error())
		return nil, err
	}
	if logs == nil {
        bs.Logger.Error().Err(err).Msg(err.Error())
		return nil, ErrFilesEmpty
	}

	return logs, nil
}

func (bs *BotService) ScrapeBot(botName string) error {
	if runtime.GOOS == "windows" {
		_, err := exec.Command("cmd", "/C", "python", "-V").Output()
		if err != nil {
            bs.Logger.Error().Err(err).Msg(err.Error())
			return err
		}

		cmd := exec.Command("cmd", "/C", "scrapy", "crawl", botName)
		cmd.Dir = "./scrapy_grocery_stores"
		_, err = cmd.Output()
		if err != nil {
            bs.Logger.Error().Err(err).Msg(err.Error())
			// TODO(miha): We can't handle errors - we just need to
			// check back 5min after bot run to see if we scraped all
			// data.
		}

		/*
			stdout, _ := cmd.StderrPipe()
			cmd.Start()
			scanner := bufio.NewScanner(stdout)
			scanner.Split(bufio.ScanWords)
			p := false
			for scanner.Scan() {
				m := scanner.Text()
				if m == "(finished)" {
					p = true
				}

				if p {
					fmt.Println(m)
				}
			}

			if err != nil {
				fmt.Println("err: ", err)
				// TODO
			}
			fmt.Println("bot: ", botName)
			//fmt.Println("out: ", string(out))
		*/

	} else {
        err := errors.New("Linux ScrapeAll() is not suported yet")
        bs.Logger.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}

/*
func (bs *BotService) GetBotStats(botName, querry string) ([]models.FileStat, error) {
	files, err := bs.db.GetBotStats(botName, querry)
	if err != nil {
		return nil, err
	}
	if files == nil {
		return nil, ErrFilesEmpty
	}

	return files, nil
}
*/

/*
func (bs *BotService) GetFileStats(botName, fileName string) (*models.FileStat, error) {
	fileName = fileName[strings.Index(fileName, "_")+1:]
	if strings.Index(fileName, ".") > 0 {
		fileName = fileName[0:strings.Index(fileName, ".")]
	}

	t, err := time.Parse("2006-01-02T15-04-05", fileName)
	if err != nil {
		return nil, err
	}
	_ = t

	file, err := bs.db.GetFileStats(botName, t.Unix())
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, ErrFilesEmpty
	}

	return file, nil
}
*/

func (bs *BotService) GetBotFile(file string) (string, error) {
	return "", nil
}

/*
func (bs *BotService) GetBotFiles(files []string) (string, error) {
	return "", nil
}
*/

func (bs *BotService) GetBotFileNames(botName string) ([]string, error) {
	dir, err := os.Open(fmt.Sprintf("./scrapy_grocery_stores/data/%s/", botName))
	if err != nil {
        bs.Logger.Error().Err(err).Msg(err.Error())
		return nil, ErrDirectoryNotFound
	}

	fileName, err := dir.Readdirnames(0)
	if err != nil {
        bs.Logger.Error().Err(err).Msg(err.Error())
		return nil, ErrDirectoryIsEmpty
	}

	names := []string{}

	for _, name := range fileName {
		names = append(names, name)
	}

	return names, nil
}

func (bs *BotService) CmdStatus(botName string) {
	fmt.Println("botPID: ", bs.botPID)
	fmt.Println("bot pid: ", bs.botPID[botName].Process.Pid)
	fmt.Println("bot pid: ", bs.botPID[botName].Process)
	fmt.Println("bot pid: ", bs.botPID[botName].ProcessState)
}

// Send command that start scraping 'botName' bot. If bot is already running
// skip starting it.
//
// Errors:
//  - ErrNoPython: System don't have python installed
//  - ErrCantStartBot: System couldn't start the bot process
func (bs *BotService) PostBotCmdScrape(botName string) error {
	// NOTE(miha): Check if python is installed on the system.
    _, err := exec.LookPath("python")
	if err != nil {
        bs.Logger.Error().Err(err).Msg(err.Error())
		return ErrNoPython
	}

    // TODO(miha): Set log file?
    cmd := common.MultiOSCommand("scrapy", "crawl", botName)
    cmd.Dir = "./scrapy_grocery_stores"

    // NOTE(miha): Bot is already running
    if _, ok := bs.botPID[botName]; ok {
        return nil
    }

    bs.botPID[botName] = cmd

    err = cmd.Start()
    if err != nil {
        bs.Logger.Error().Err(err).Msg(err.Error())
        return ErrCantStartBot
    }

    // NOTE(miha): This go-routine triggers when process finished - can
    // return errors!
    go func(bot string) {
        err := cmd.Wait()
        if err != nil {
            bs.Logger.Error().Err(err).Msg(err.Error())
        }
        err = bs.db.UpdateBot(bot)
        if err != nil {
            bs.Logger.Error().Err(err).Msg(err.Error())
        }
        delete(bs.botPID, botName)
    }(botName)

	return nil
}

// Send command that stop all the bots.
//
// Errors:
//  - ErrCantKillBot: System couldn't kill the bot process
func (bs *BotService) BotCmdStop(botName string) error {
    if process, ok := bs.botPID[botName]; ok {
        err := process.Process.Kill()
        if err != nil {
            bs.Logger.Error().Err(err).Msg(err.Error())
            return ErrCantKillBot
        }
    }
    return nil
}

// Check if bot named 'botName' is currently running
func (bs *BotService) BotCmdStatus(botName string) *models.BotStatus {
    if _, ok := bs.botPID[botName]; ok {
        return &models.BotStatus{Running: true}
    } else {
        return &models.BotStatus{Running: false}
    }
}
