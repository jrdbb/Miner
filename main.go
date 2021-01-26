package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/CommonProsperity/Miner/crawler"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
)

var (
	influxAddr = flag.String("influxAddr", "127.0.0.1:8080", "The influx db address in the format of host:port")
	influxToken = flag.String("influxToken", "", "The influx db token")
	sdate = flag.String("sdate", "2020-01-01", "start date")
	edate = flag.String("edate", "2020-02-01", "end date")
	logFile    = flag.String("logFile", "", "The log file of GoDock")
	logLevel   = flag.String("logLevel", "Info", "The log level. (Debug/Info/Warn/Error)")
)

func initLog() {
	log.SetFormatter(&log.JSONFormatter{})
	if *logFile == "" {
		log.SetOutput(os.Stdout)
	} else {
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("error opening logFile(%s): %v", *logFile, err))
		}
		log.SetOutput(f)
	}

	if *logLevel == "Debug" {
		log.SetLevel(log.DebugLevel)
	} else if *logLevel == "Info" {
		log.SetLevel(log.InfoLevel)
	} else if *logLevel == "Warn" {
		log.SetLevel(log.WarnLevel)
	} else if *logLevel == "Error" {
		log.SetLevel(log.ErrorLevel)
	} else {
		panic(fmt.Sprintf("Unknown logLevel: %s", *logLevel))
	}
}

func init() {
	flag.Parse()
	initLog()
}

type callback struct {
	crl crawler.FundCrawler
	fundDefAPI api.WriteAPI
	fundValueAPI api.WriteAPI
}

func (cb *callback) OnBasicFund(funds []*crawler.BasicFund){
	for _, fund := range funds {
		p := influxdb2.NewPoint(
			"Fund",
			map[string]string{"id": fund.ID},
			map[string]interface{}{"id": fund.ID, "name_cn": fund.NameCN, "name_py": fund.NamePY},
			time.Now(),
		)
		cb.fundDefAPI.WritePoint(p)
	}

	for _, fund := range funds {
		cb.crl.GetHistoryValue(false, fund.ID, 1, *sdate, *edate)
	}
}

func (cb *callback) OnHistoryValue(apiData *crawler.ApiData){
	for _, v := range apiData.Content {
		p := influxdb2.NewPoint(
			"FundValue",
			map[string]string{"id": apiData.Code},
			map[string]interface{}{"value": v.NetAssetValue},
			v.Date,
		)
		cb.fundDefAPI.WritePoint(p)
	}

	if apiData.CurrPage < apiData.Pages {
		cb.crl.GetHistoryValue(false, apiData.Code, int(apiData.CurrPage) + 1, *sdate, *edate)
	}
}

func main() {
	client := influxdb2.NewClient(*influxAddr, *influxToken)
	defer client.Close()

	crl := crawler.NewFundCrawler()
	crl.SetCallBack(&callback{
		crl: crl,
		fundDefAPI: client.WriteAPI("CommonProsperity", "fund_def"),
		fundValueAPI: client.WriteAPI("CommonProsperity", "fund_history_value"),
	})

	crl.GetAllBasicFund(false)

	fmt.Println("stop by enter")
    input := bufio.NewScanner(os.Stdin)
    input.Scan()
    fmt.Println("stoping")
}
