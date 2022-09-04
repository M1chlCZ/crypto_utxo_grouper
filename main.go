package main

import (
	"cryptoGroup/coind"
	"cryptoGroup/models"
	"cryptoGroup/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

func main() {
	a, err := utils.GetENV("GROUP_ADDRESS")
	if err != nil || len(a) == 0 {
		utils.WrapErrorLog("Could not get GROUP_ADDRESS in env file")
		return
	}
	b, err := utils.GetENV("GROUP_AMOUNT")
	if err != nil || len(b) == 0 {
		utils.WrapErrorLog("Could not get GROUP_AMOUNT in env file")
		return
	}
	c, err := utils.GetENV("WALLET_USER")
	if err != nil || len(c) == 0 {
		utils.WrapErrorLog("Could not get WALLET_USER in env file")
		return
	}
	d, err := utils.GetENV("WALLET_PASS")
	if err != nil || len(d) == 0 {
		utils.WrapErrorLog("Could not get WALLET_PASS in env file")
		return
	}
	e, err := utils.GetENV("WALLET_PORT")
	if err != nil || len(e) == 0 {
		utils.WrapErrorLog("Could not get WALLET_PORT in env file")
		return
	}

	groupTX()
}

func groupTX() {
	utils.ReportMessage("Grouping TX")
	dm, err := getStakingDaemon()
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}

	address, _ := utils.GetENV("GROUP_ADDRESS")
	ga, _ := utils.GetENV("GROUP_AMOUNT")
	groupAmount, errNum := utils.StringToFloat64(ga)
	if errNum != nil {
		utils.WrapErrorLog("GROUP_AMOUNT is not a float, needs to be in format: 10.0")
		return
	}
	utils.ReportMessage(fmt.Sprintf("GROUP ADDRESS %s", address))
	utils.ReportMessage(fmt.Sprintf("GROUP AMOUNT %f", groupAmount))

outerLoop:
	for {
		amount := 0.0
		numberOfInputs := 0
		res, err := utils.WrapDaemon(dm, 10, "listunspent")
		if err != nil {
			utils.WrapErrorLog(err.Error())
			return
		}
		var ing models.ListUnspent
		errJson := json.Unmarshal(res, &ing)
		if errJson != nil {
			utils.WrapErrorLog(errJson.Error())
			return
		}
	innerLoop:
		for _, unspent := range ing {
			if numberOfInputs == 25 || amount > groupAmount {
				break innerLoop
			}
			if unspent.Spendable && (unspent.Amount+amount) < groupAmount {
				amount += unspent.Amount
				numberOfInputs++
			}
		}
		if numberOfInputs == 0 {
			break outerLoop
		}
		utils.ReportMessage(fmt.Sprintf("Amount %f, %d UTXO, deposit addr %s", amount, numberOfInputs, address))
		_, errSend := sendCoins(models.UnstakeClientReq{
			CoinID:  dm.CoinID,
			Amount:  amount,
			Deposit: address,
		}, dm)
		if errSend != nil {
			utils.WrapErrorLog("N0 spendable outputs, waiting 30 seconds before trying again")
			time.Sleep(time.Second * 30)
		}
	}

	utils.ReportMessage("GroupTX done")
}

func getStakingDaemon() (models.Daemon, error) {
	dm := models.Daemon{}
	walletUser, err := utils.GetENV("WALLET_USER")
	if err != nil {
		return dm, err
	}
	walletPass, err := utils.GetENV("WALLET_PASS")
	if err != nil {
		return dm, err
	}

	walletPort, err := utils.GetENV("WALLET_PORT")
	if err != nil {
		return dm, err
	}

	wp, errNum := strconv.Atoi(walletPort)
	if errNum != nil {
		utils.WrapErrorLog("WALLET_PORT is not a number")
		return dm, errors.New("WALLET_PORT is not a number")
	}

	var passPhrase sql.NullString
	walletPassPhrase, err := utils.GetENV("WALLET_PASSPHRASE")
	if err != nil || len(walletPassPhrase) == 0 {
		utils.ReportMessage("Wallet passphrase not found, not using passphrase")
		passPhrase = sql.NullString{
			Valid:  false,
			String: "",
		}
	} else {
		passPhrase = sql.NullString{
			Valid:  true,
			String: walletPassPhrase,
		}
	}

	dm.WalletUser = walletUser
	dm.WalletPass = walletPass
	dm.WalletPort = wp
	dm.PassPhrase = passPhrase

	return dm, nil

}

func sendCoins(request models.UnstakeClientReq, daemon models.Daemon) (string, error) {
	client, errClient := coind.New("127.0.0.1", daemon.WalletPort, daemon.WalletUser, daemon.WalletPass, false, 30)
	if errClient != nil {
		log.Println(errClient.Error())
		//utils.ReportError(w, "Wallet coin id is unreachable", http.StatusInternalServerError)
		return "", errClient
	}
	if daemon.PassPhrase.Valid {
		_, erUnlock := client.Call("walletpassphrase", daemon.PassPhrase.String, 1000)
		if erUnlock != nil {
			utils.WrapErrorLog("error unlock" + erUnlock.Error())
			return "", erUnlock
		}
		time.Sleep(500 * time.Millisecond)
	}

	utils.ReportMessage(fmt.Sprintf("Amount %f, deposit addr %s", request.Amount, request.Deposit))
	txid, er := client.Call("sendtoaddress", request.Deposit, request.Amount)
	if er != nil {
		log.Println(er.Error())
		//utils.ReportError(w, er.Error(), http.StatusInternalServerError)
		return "", er
	}
	if string(txid) == "null" {
		//utils.ReportError(w, utils.TrimQuotes(string(txid)), http.StatusConflict)
		return "", errors.New("null")
	}
	time.Sleep(500 * time.Millisecond)
	if daemon.PassPhrase.Valid {
		_, erLock := client.Call("walletlock")
		if erLock != nil {
			utils.WrapErrorLog(erLock.Error())
			log.Println("error unlock" + erLock.Error())
			//utils.ReportError(w, "Wallet coin id is unreachable", http.StatusInternalServerError)
			return "", erLock
		}
		time.Sleep(500 * time.Millisecond)
	}
	if daemon.PassPhrase.Valid {
		_, erLock := client.Call("walletpassphrase", daemon.PassPhrase.String, 999999999, true)
		if erLock != nil {
			utils.WrapErrorLog(erLock.Error())
			log.Println("error unlock" + erLock.Error())
			//utils.ReportError(w, "Wallet coin id is unreachable", http.StatusInternalServerError)
			return "", erLock
		}
	}
	client = nil
	return string(txid), nil
}
