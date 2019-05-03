package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/google/logger"
	"github.com/google/uuid"
	stocktwits "github.com/khbrendle/stocktwits/streams"
	_ "github.com/lib/pq"
)

func main() {
	configFile := flag.String("config-file", fmt.Sprintf("%s/src/config.json", os.Getenv("PWD")), "the path to the config file")
	configFile := flag.String("log-file", "/var/log/stockdata.log", "the path to the config file")
	var a App
	a.init()
	defer a.DB.Close()
	defer a.LogFile.Close()
	defer a.Logger.Close()

	var err error
	if err = a.Config.Read(configFile); err != nil {
		a.Logger.Fatal(err)
	}

	var id string
	var ct int
	ticker := time.NewTicker(time.Hour)
	for _ := range ticker.C {
		ct = time.Now().Hour()
		if (ct > 9) && (ct < 5) {
			for _, symbol := range config.Stocks {
				id = uuid.New().String()
				a.GetAndRecord(id, symbol)
			}
		}
	}

}

type Config struct {
	Stocks []StockMeta
}

func (c *Config) Read(path string) error {
	var err error
	// Open our jsonFile
	var jsonFile *os.File
	if jsonFile, err = os.Open(path); err != nil {
		return err
	}
	var byteValue []byte
	if byteValue, err = ioutil.ReadAll(jsonFile); err != nil {
		return err
	}
	defer jsonFile.Close()
	var config Config
	if err = json.Unmarshal(byteValue, &config); err != nil {
		return nil
	}
	*c = config
	return nil
}

type App struct {
	Config
	DB      *sql.DB
	LogFile *os.File
	Logger  *logger.Logger
}

func (a *App) init(logFile string) {
	a.initLogger(logFile)
	a.initDB()
}

func (a *App) initLogger(logPath string) {
	var err error
	if a.LogFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660); err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	a.Logger = logger.Init("StocksData", false, true, a.LogFile)
	logger.SetFlags(11)
}

func (a *App) initDB() {
	config := map[string]string{
		"dbhost": "localhost",
		"dbport": "15432",
		"dbuser": "webapp",
		"dbpass": "webapp",
		"dbname": "stocks"}
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config["dbhost"], config["dbport"], config["dbuser"], config["dbpass"], config["dbname"])

	if a.DB, err = sql.Open("postgres", psqlInfo); err != nil {
		a.Logger.Fatal(err)
	}
	if err = a.DB.Ping(); err != nil {
		a.Logger.Fatal(err)
	}
	a.Logger.Info("Successfully connected to db!")
}

type SymbolData struct {
	Id       string
	Datetime string
	Symbol   string
	Response string
	Error    string
}

func (a *App) GetAndRecord(id string, symbol string) {
	var x SymbolData
	x.Id = id
	x.Datetime = time.Now().Format("2006-01-02 15:04:05-07")
	a.Logger.Infof("%s - Getting data for %s", id, symbol)
	x.Symbol = symbol
	var res *http.Response
	var err, ReqErr error
	if res, err = stocktwits.GetStreamSymbol(symbol); err != nil {
		// return err
		a.Logger.Errorf("%s - %v", id, err)
	}
	a.Logger.Infof("%s - Reading response body", id)
	var body []byte
	if body, ReqErr = ioutil.ReadAll(res.Body); ReqErr != nil {
		x.Error = err.Error()
	} else {
		x.Response = string(body)
	}
	defer res.Body.Close()
	a.Logger.Infof("%s - Inserting data to DB", id)
	if _, err = a.DB.Exec("INSERT INTO symbolstream.data VALUES ($1, $2, $3, $4, $5)", x.Id, x.Datetime, x.Symbol, x.Response, x.Error); err != nil {
		// return err
		a.Logger.Errorf("%s - %v", id, err)
	}
	a.Logger.Infof("%s - Returning successfully", id)
	// return nil
}
