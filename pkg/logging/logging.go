package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var BlackboxLogger *log.Logger
var BlackboxCLILogger *log.Logger

var CronLogger *log.Logger

func init() {
	t := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("logs/blacboxlog_" + t + ".log")
	logFile, err := openLogFile(fileName)
	cronlogFileName := fmt.Sprintf("logs/cronlog_" + t + ".log")
	cronLogFile, err := openLogFile(cronlogFileName)
	if err != nil {
		log.Fatal(err)
	}
	// close log file in warm shutdown
	BlackboxLogger = log.New(logFile, "[Blackbox]", log.Ldate|log.Ltime)
	BlackboxCLILogger = log.New(os.Stdout, "[Blackbox]", log.Ldate|log.Ltime)
	CronLogger = log.New(cronLogFile, "[Cron]", log.Ldate|log.Ltime)
}

func openLogFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return file, nil
}
