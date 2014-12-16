package goslack

import (
	"fmt"
	"log"
	"os"
)

var debugLog *log.Logger

func InitLogger() {

	file, err := os.OpenFile("gopher.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Cannot open log file. ERR: %v", err)
	}

	debugLog = log.New(file, "goslack DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}
