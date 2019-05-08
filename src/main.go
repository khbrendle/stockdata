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
	logFile := flag.String("log-file", "/var/log/stockdata.log", "the path to the config file")
	var a App
	var err error
	if err = a.init(*configFile, *logFile); err != nil {
		a.Logger.Fatal(err)
	}
	defer a.DB.Close()
	defer a.LogFile.Close()
	defer a.Logger.Close()

	var id string
	// var ct int
	// ticker := time.NewTicker(time.Hour)
	// for _ := range ticker.C {
	// 	ct = time.Now().Hour()
	// 	if (ct > 9) && (ct < 5) {
	for i, symbol := range a.Config.Stocks {
		id = uuid.New().String()
		fmt.Printf("%s, running symbol `%s`\n", id, symbol.Symbol)
		a.Config.Stocks[i] = a.GetAndRecord(id, symbol)
	}
	a.Config.Write(*configFile)
	// 	}
	// }

}

type Config struct {
	Stocks []StockInfo
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

func (c *Config) Write(path string) error {
	var err error
	var dat []byte
	if dat, err = json.Marshal(c); err != nil {
		return err
	}
	if err = ioutil.WriteFile(path, dat, 0644); err != nil {
		return err
	}
	return nil
}

type App struct {
	Config
	DB      *sql.DB
	LogFile *os.File
	Logger  *logger.Logger
}

func (a *App) init(configFile string, logFile string) error {
	a.initLogger(logFile)
	a.initDB()

	if err := a.Config.Read(configFile); err != nil {
		return err
	}
	return nil
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
		a.Logger.Fatal(fmt.Sprintf("error opening database connection: %s", err.Error()))
	}
	if err = a.DB.Ping(); err != nil {
		a.Logger.Fatal(fmt.Sprintf("error pinging database: %s", err.Error()))
	}
	a.Logger.Info("Successfully connected to database!")
}

type StockInfo struct {
	Symbol string
	MaxId  int
	MinId  int
}

type SymbolData struct {
	Id       string
	Datetime string
	Symbol   string
	Response string
	Error    string
}

func (a *App) GetAndRecord(id string, symbol StockInfo) StockInfo {
	var x SymbolData
	x.Id = id
	x.Datetime = time.Now().Format("2006-01-02 15:04:05-07")
	a.Logger.Infof("%s - Getting data for `%s`", id, symbol.Symbol)
	x.Symbol = symbol.Symbol
	var res *http.Response
	var err, ReqErr error
	params := map[string]interface{}{"since": symbol.MaxId}
	if res, err = stocktwits.GetStreamSymbol(x.Symbol, params); err != nil {
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
	a.Logger.Infof("%s - Inserting data to database", id)
	if _, err = a.DB.Exec("INSERT INTO symbolstream.response_raw VALUES ($1, $2, $3, $4, $5)", x.Id, x.Datetime, x.Symbol, x.Response, x.Error); err != nil {
		a.Logger.Errorf("%s - %v", id, err)
	}
	var tmp stocktwits.StreamSymbol
	if err = json.Unmarshal(body, &tmp); err != nil {
		a.Logger.Errorf("%s - %v", id, err)
	} else {
		symbol.MaxId = tmp.GetMaxId()
		symbol.MinId = tmp.GetMinId()
	}
	a.Logger.Infof("%s - Returning successfully", id)
	return symbol
}
