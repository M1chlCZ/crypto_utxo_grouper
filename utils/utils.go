package utils

import (
	"cryptoGroup/coind"
	"cryptoGroup/models"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	_ "io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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

func WrapDaemon(daemon models.Daemon, maxTries int, method string, params ...interface{}) ([]byte, error) {
	var wg sync.WaitGroup
	c := make(chan []byte)
	e := make(chan error)
	var res []byte
	var err error
	go callDaemon(c, e, &wg, &daemon, maxTries, method, params)
	wg.Wait()
	select {
	case data := <-c:
		res = data
	case er := <-e:
		WrapErrorLog(er.Error())
		err = er
	}
	close(c)
	close(e)

	if err == nil {
		return res, nil
	} else {
		return nil, err
	}
}

func callDaemon(c chan []byte, e chan error, wg *sync.WaitGroup, daemon *models.Daemon, triesMax int, command string, params any) {
	defer wg.Done()
	wg.Add(1)
	var client *coind.Coind
	var errClient error
	tries := 0
	for {
		if tries != 0 {
			ReportMessage(fmt.Sprintf("Try %d of %d. CMD: %s Daemon: %s", tries, triesMax, command, daemon.Folder))
		}
		if tries >= triesMax {
			e <- errors.New("error getting RPC data")
			break
		}

		tries++
		if client == nil {
			client, errClient = coind.New("127.0.0.1", daemon.WalletPort, daemon.WalletUser, daemon.WalletPass, false, 45)
			if errClient != nil {
				ReportMessage(errClient.Error())
				time.Sleep(15 * time.Second)
				client = nil
				continue
			}
		}
		p, er := client.Call(command, params)

		if er != nil {
			ReportMessage("Daemon unreachable " + er.Error())
			time.Sleep(15 * time.Second)
			client = nil
			continue
		}
		if string(p) != "null" {
			if len(p) != 0 {
				c <- p
				client = nil
				break
			}
		}

		ReportMessage("Error, trying again")
		time.Sleep(15 * time.Second)
	}
}
