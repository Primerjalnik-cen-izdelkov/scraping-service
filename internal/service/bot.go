package service

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"scraping_service/internal/database"
	"scraping_service/pkg/models"
	"strings"
	"sync"
	"time"
)

var (
	ErrDirectoryNotFound = errors.New("Directory not found")
	ErrDirectoryIsEmpty  = errors.New("Directory is empty")
	ErrBotBadFilename    = errors.New("Bot filename does not contain a dot")
	ErrFilesEmpty        = errors.New("Returned files are empty")
)

type BotService struct {
	db *database.Database
}

func CreateBotService(d *database.Database) *BotService {
	return &BotService{db: d}
}

func (bs *BotService) GetAllBots() ([]models.Bot, error) {
	files, err := os.Open("./scrapy_grocery_stores/scrapy_grocery_stores/spiders")
	if err != nil {
		return nil, ErrDirectoryNotFound
	}

	dirs, err := files.Readdirnames(0)
	if err != nil {
		return nil, ErrDirectoryIsEmpty
	}

	bots := []models.Bot{}

	for _, dir := range dirs {
		if dir[0] == '_' {
			continue
		}
		err := strings.Index(dir, ".")
		if err == -1 {
			return nil, ErrBotBadFilename
		}
		trimmedDir := dir[0:strings.Index(dir, ".")]
		bots = append(bots, models.Bot{Name: trimmedDir})
	}

	return bots, nil
}

func (bs *BotService) ScrapeAllBots() error {
	bots, err := bs.GetAllBots()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(bots))

	if runtime.GOOS == "windows" {
		_, err := exec.Command("cmd", "/C", "python", "-V").Output()
		if err != nil {
			return err
		}

		for _, b := range bots {
			go func(b *models.Bot, wg *sync.WaitGroup) {
				cmd := exec.Command("cmd", "/C", "scrapy", "crawl", b.Name)
				cmd.Dir = "./scrapy_grocery_stores"
				_, err := cmd.Output()
				if err != nil {
					// TODO(miha): We can't handle errors - we just need to
					// check back 5min after bot run to see if we scraped all
					// data.
				}

				wg.Done()
			}(&b, wg)
		}

	} else {
		return errors.New("Linux ScrapeAll() is not suported yet")
	}

	wg.Wait()
	return nil
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

func (bs *BotService) GetBotFile(file string) (string, error) {
	return "", nil
}

func (bs *BotService) GetBotFiles(files []string) (string, error) {
	return "", nil
}

func (bs *BotService) GetBotFileNames(botName string) ([]models.File, error) {
	files, err := bs.db.GetBotFileNames(botName)
	if err != nil {
		return nil, err
	}
	if files == nil {
		return nil, ErrFilesEmpty
	}

	for i, file := range files {
		name := fmt.Sprintf("%s_%d-%02d-%02dT%02d-%02d-%02d.csv", botName,
			file.StartTime.Year(), file.StartTime.Month(), file.StartTime.Day(),
			file.StartTime.Hour(), file.StartTime.Minute(), file.StartTime.Second())
		files[i].Name = name
	}

	return files, nil
}
