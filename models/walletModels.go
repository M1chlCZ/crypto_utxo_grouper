package models

import "fmt"

type GetInfo struct {
	Version         int     `json:"version"`
	Protocolversion int     `json:"protocolversion"`
	Walletversion   int     `json:"walletversion"`
	Balance         float64 `json:"balance"`
	Blocks          int     `json:"blocks"`
	Timeoffset      int     `json:"timeoffset"`
	Connections     int     `json:"connections"`
	Proxy           string  `json:"proxy"`
	Difficulty      float64 `json:"difficulty"`
	Testnet         bool    `json:"testnet"`
	Moneysupply     float64 `json:"moneysupply"`
	Keypoololdest   int     `json:"keypoololdest"`
	Keypoolsize     int     `json:"keypoolsize"`
	UnlockedUntil   int     `json:"unlocked_until"`
	Paytxfee        float64 `json:"paytxfee"`
	Relayfee        float64 `json:"relayfee"`
	StakingStatus   string  `json:"staking status"`
	Errors          string  `json:"errs"`
}

type GetInfoXDN struct {
	Version         string  `json:"version"`
	Protocolversion int     `json:"protocolversion"`
	Walletversion   int     `json:"walletversion"`
	Balance         float64 `json:"balance"`
	Newmint         float64 `json:"newmint"`
	Stake           float64 `json:"stake"`
	Blocks          int     `json:"blocks"`
	Timeoffset      int     `json:"timeoffset"`
	Moneysupply     float64 `json:"moneysupply"`
	Connections     int     `json:"connections"`
	Proxy           string  `json:"proxy"`
	IP              string  `json:"ip"`
	Difficulty      struct {
		ProofOfWork  float64 `json:"proof-of-work"`
		ProofOfStake float64 `json:"proof-of-stake"`
	} `json:"difficulty"`
	Testnet       bool    `json:"testnet"`
	Keypoololdest int     `json:"keypoololdest"`
	Keypoolsize   int     `json:"keypoolsize"`
	Paytxfee      float64 `json:"paytxfee"`
	Mininput      float64 `json:"mininput"`
	Errors        string  `json:"errors"`
}

type ListTransactions []struct {
	Account         string        `json:"account"`
	Address         string        `json:"address"`
	Category        string        `json:"category"`
	Amount          float64       `json:"amount"`
	Vout            int           `json:"vout"`
	Confirmations   int           `json:"confirmations"`
	Bcconfirmations int           `json:"bcconfirmations"`
	Blockhash       string        `json:"blockhash"`
	Blockindex      int           `json:"blockindex"`
	Blocktime       int           `json:"blocktime"`
	Txid            string        `json:"txid"`
	Walletconflicts []interface{} `json:"walletconflicts"`
	Time            int           `json:"time"`
	Timereceived    int           `json:"timereceived"`
}

type GetTransaction struct {
	Coinid          string        `json:"coinID"`
	Nodeid          int           `json:"nodeID"`
	Amount          float64       `json:"amount"`
	Fee             float64       `json:"fee"`
	Confirmations   int           `json:"confirmations"`
	Bcconfirmations int           `json:"bcconfirmations"`
	Generated       bool          `json:"generated"`
	Blockhash       string        `json:"blockhash"`
	Blockindex      int           `json:"blockindex"`
	Blocktime       int           `json:"blocktime"`
	Txid            string        `json:"txid"`
	Walletconflicts []interface{} `json:"walletconflicts"`
	Time            int           `json:"time"`
	Timereceived    int           `json:"timereceived"`
	Details         []struct {
		Account  string  `json:"account"`
		Category string  `json:"category"`
		Amount   float64 `json:"amount"`
		Vout     int     `json:"vout"`
		Fee      float64 `json:"fee,omitempty"`
		Address  string  `json:"address,omitempty"`
	} `json:"details"`
	Hex string `json:"hex"`
}

func (d *GetTransaction) ToString() string {
	return fmt.Sprintf("[%s, %s, %v, %f]", d.Coinid, d.Txid, d.Generated, d.Amount)
}

type GetTransactionXDN struct {
	Txid     string `json:"txid"`
	Version  int    `json:"version"`
	Time     int    `json:"time"`
	Locktime int    `json:"locktime"`
	Vin      []struct {
		Coinbase string `json:"coinbase"`
		Sequence int    `json:"sequence"`
	} `json:"vin"`
	Vout []struct {
		Value        float64 `json:"value"`
		N            int     `json:"n"`
		ScriptPubKey struct {
			Asm       string   `json:"asm"`
			Hex       string   `json:"hex"`
			ReqSigs   int      `json:"reqSigs"`
			Type      string   `json:"type"`
			Addresses []string `json:"addresses"`
		} `json:"scriptPubKey"`
	} `json:"vout"`
	Amount          float64       `json:"amount"`
	Confirmations   int           `json:"confirmations"`
	Bcconfirmations int           `json:"bcconfirmations"`
	Generated       bool          `json:"generated"`
	Blockhash       string        `json:"blockhash"`
	Blockindex      int           `json:"blockindex"`
	Blocktime       int           `json:"blocktime"`
	Walletconflicts []interface{} `json:"walletconflicts"`
	Timereceived    int           `json:"timereceived"`
	Details         []struct {
		Account  string  `json:"account"`
		Address  string  `json:"address"`
		Category string  `json:"category"`
		Amount   float64 `json:"amount"`
	} `json:"details"`
}

type ListUnspent []struct {
	Txid          string  `json:"txid"`
	Vout          int     `json:"vout"`
	Address       string  `json:"address"`
	Account       string  `json:"account"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	Confirmations int     `json:"confirmations"`
	Spendable     bool    `json:"spendable"`
}

type MasternodeList []struct {
	Rank       int    `json:"rank"`
	Network    string `json:"network"`
	Txhash     string `json:"txhash"`
	Outidx     int    `json:"outidx"`
	Status     string `json:"status"`
	Addr       string `json:"addr"`
	Version    int    `json:"version"`
	Lastseen   int    `json:"lastseen"`
	Activetime int    `json:"activetime"`
	Lastpaid   int    `json:"lastpaid"`
}

type MasternodeStatus struct {
	Txhash    string `json:"txhash"`
	Outputidx int    `json:"outputidx"`
	Netaddr   string `json:"netaddr"`
	Addr      string `json:"addr"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
}

type MasternodeStatusXDN struct {
	Vin              string `json:"vin"`
	Service          string `json:"service"`
	Status           int    `json:"status"`
	Pubkey           string `json:"pubkey"`
	NotCapableReason string `json:"notCapableReason"`
}

type MasternodeOutputs []struct {
	Txhash    string `json:"txhash"`
	Outputidx int    `json:"outputidx"`
}

type RawTxArray struct {
	Txid string `json:"txid"`
	Vout int    `json:"vout"`
}

type SignRawTransaction struct {
	Hex      string `json:"hex"`
	Complete bool   `json:"complete"`
}

type SendToAddress struct {
	Txid string `json:"txid"`
}
