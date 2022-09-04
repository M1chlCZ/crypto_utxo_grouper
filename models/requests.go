package models

type UnstakeClientReq struct {
	CoinID  int     `json:"idCoin"`
	Amount  float64 `json:"amount"`
	Deposit string  `json:"deposit"`
}
