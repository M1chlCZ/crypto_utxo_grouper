package utils

import (
	"github.com/joho/godotenv"
	"io"
	_ "io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func GetENV(key string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		WrapErrorLog("Error loading .env file")
		return "", err
	}
	return os.Getenv(key), nil
}

func logToFile(message string) {
	f, err := os.OpenFile("api.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v\n", err)
	}
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	log.Println(message)
	_ = f.Close()
}

func WrapErrorLog(message string) {
	if !strings.Contains(message, "tx_id_UNIQUE") {
		logToFile("//// - ERROR - ////")
		logToFile(message)
		logToFile("////===========////")
		logToFile("")
	}
}

func ReportMessage(message string) {
	logToFile(message)
	logToFile("")
}

func StringToFloat64(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// handle error
		return 0.0, err
	}
	return f, nil
}
