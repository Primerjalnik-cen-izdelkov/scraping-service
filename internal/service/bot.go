package service

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"scraping_service/internal/database"
	"scraping_service/pkg/common"
	"scraping_service/pkg/models"
	"strings"
)

var (
	ErrDirectoryNotFound = errors.New("Directory not found")
	ErrDirectoryIsEmpty  = errors.New("Directory is empty")
	ErrBotBadFilename    = errors.New("Bot filename does not contain a dot")
	ErrFilesEmpty        = errors.New("Returned files are empty")
    ErrBotAlreadyRunning = errors.New("Bot is alredy running")
    ErrBotIsNotRunning   = errors.New("Bot is not running")
)

type BotService struct {
	db     *database.Database
	botPID map[string]*exec.Cmd
}

func CreateBotService(d *database.Database) *BotService {
	return &BotService{db: d, botPID: make(map[string]*exec.Cmd)}
}

func (bs *BotService) BotNames() ([]string, error) {
	files, err := os.Open("./scrapy_grocery_stores/scrapy_grocery_stores/spiders")
	if err != nil {
		return nil, ErrDirectoryNotFound
	}

	dirs, err := files.Readdirnames(0)
	if err != nil {
		return nil, ErrDirectoryIsEmpty
	}

	bots := []string{}

	for _, dir := range dirs {
		// NOTE(miha): Parse spiders directory (ignore python files)
		if dir[0] == '_' {
			continue
		}
		index := strings.Index(dir, ".")
		if index == -1 {
			return nil, ErrBotBadFilename
		}
		trimmedDir := dir[0:strings.Index(dir, ".")]

		bots = append(bots, trimmedDir)
	}

	return bots, nil
}


func (bs *BotService) GetBots() ([]*models.Bot, error) {
    botNames, err := bs.BotNames()
    fmt.Println("BotNames:", botNames)
    if err != nil {
        return nil, err
    }

	bots := []*models.Bot{}

	for _, name := range botNames {
        bot, err := bs.db.GetBot(name)
        if err != nil {

        }

        status := bs.BotCmdStatus(name)
        bot.Status = status

		bots = append(bots, bot)
	}

	return bots, nil
}

// TODO(miha): Call get bots first, and then call getBotFiles on each bot to
// compose data.
func (bs *BotService) GetFiles() ([]string, error) {
    botName := ""
	dir, err := os.Open(fmt.Sprintf("./scrapy_grocery_stores/data/%s/", botName))
	if err != nil {
		return nil, ErrDirectoryNotFound
	}

	files, err := dir.Readdirnames(0)
	if err != nil {
		return nil, ErrDirectoryIsEmpty
	}

	names := []string{}

	for _, name := range files {
		names = append(names, name)
	}

	return names, nil
}

func (bs *BotService) GetLogs() error {
    return nil
}

func (bs *BotService) GetCmds() error {
    return nil
}

func (bs *BotService) PostCmdScrape() error {
	bots, err := bs.BotNames()
	if err != nil {
		return err
	}

	// NOTE(miha): Check if python is installed on the system.
	_, err = exec.LookPath("python")
	if err != nil {
		return err
	}

	for _, bot := range bots {
		cmd := common.MultiOSCommand("scrapy", "crawl", bot)
		cmd.Dir = "./scrapy_grocery_stores"

		// TODO(miha): Set stdout and stderr files, so we can tail them to see
		// status.
		// NOTE(miha): Everything we write in python is send to stderr
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

		if val, ok := bs.botPID[bot]; ok {
			if val.ProcessState.ExitCode() == -1 {
				fmt.Println("scraping for bot ", bot, " already running", val.Process.Pid)
			}
		}

		bs.botPID[bot] = cmd

		// TODO(miha): Check if the bot is already running in bs.botPID

		err = cmd.Start()
		if err != nil {
			return err
		}

		fmt.Println("B, bot exit code: ", bot, bs.botPID[bot].ProcessState.ExitCode(), cmd.Process.Pid)

		// NOTE(miha): This go-routine triggers when process finished - can
		// return errors!
		go func(bot string) {
			err := cmd.Wait()
			if err != nil {

			}

            err = bs.db.UpdateBot(bot)
			if err != nil {

			}

			fmt.Println("process finished for bot: ", bot)
		}(bot)
	}

	return nil
}

func (bs *BotService) PostCmdStop() error {
    return nil
}

func (bs *BotService) PostCmdStatus() error {
    return nil
}

func (bs *BotService) GetBotCmds() error {
    return nil
}

func (bs *BotService) GetBotFiles() error {
    return nil
}

func (bs *BotService) GetBotLogs(botName string, qm map[string]string) ([]models.FileLog, error) {
	logs, err := bs.db.GetBotLogs(botName, qm)
	if err != nil {
		return nil, err
	}
	if logs == nil {
		return nil, ErrFilesEmpty
	}

	return logs, nil
}

func (bs *BotService) ScrapeBot(botName string) error {
	if runtime.GOOS == "windows" {
		_, err := exec.Command("cmd", "/C", "python", "-V").Output()
		if err != nil {
			return err
		}

		cmd := exec.Command("cmd", "/C", "scrapy", "crawl", botName)
		cmd.Dir = "./scrapy_grocery_stores"
		_, err = cmd.Output()
		if err != nil {
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
		return errors.New("Linux ScrapeAll() is not suported yet")
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
		return nil, ErrDirectoryNotFound
	}

	fileName, err := dir.Readdirnames(0)
	if err != nil {
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

func (bs *BotService) BotCmdScrape(botName string) (int, error) {
	// NOTE(miha): Check if python is installed on the system.
    _, err := exec.LookPath("python")
	if err != nil {
		return -1, err
	}

    // TODO(miha): Set log file?
    cmd := common.MultiOSCommand("scrapy", "crawl", botName)
    cmd.Dir = "./scrapy_grocery_stores"

    if val, ok := bs.botPID[botName]; ok {
        if val.ProcessState.ExitCode() == -1 {
            return -1, ErrBotAlreadyRunning
        }
    }

    bs.botPID[botName] = cmd

    err = cmd.Start()
    if err != nil {
        return -1, err
    }

    // NOTE(miha): This go-routine triggers when process finished - can
    // return errors!
    go func(botName string) {
        err := cmd.Wait()
        if err != nil {

        }

        delete(bs.botPID, botName)
        fmt.Println("process finished for bot: ", botName)
    }(botName)

	return cmd.Process.Pid, nil
}

func (bs *BotService) BotCmdStop(botName string) error {
    if process, ok := bs.botPID[botName]; ok {
        return process.Process.Kill()
    } else {
        return ErrBotIsNotRunning
    }
}

func (bs *BotService) BotCmdStatus(botName string) *models.BotStatus {
    if _, ok := bs.botPID[botName]; ok {
        return &models.BotStatus{Running: true}
    } else {
        return &models.BotStatus{Running: false}
    }
}
