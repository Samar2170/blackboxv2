package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var BlackboxLogger *log.Logger
var BlackboxCLILogger *log.Logger

func init() {
	t := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("logs/blacboxlog_" + t + ".log")
	logFile, err := openLogFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	// close log file in warm shutdown
	BlackboxLogger = log.New(logFile, "[Blackbox]", log.Ldate|log.Ltime)
	BlackboxCLILogger = log.New(os.Stdout, "[Blackbox]", log.Ldate|log.Ltime)
}

func openLogFile(fileName string) (*os.File, error) {
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		file, err = os.Create(fileName)
		if err != nil {
			return file, err
		}
		if err := file.Chmod(0777); err != nil {
			return file, err
		}
	} else if err != nil {
		panic(err)
	}
	return file, nil
}
