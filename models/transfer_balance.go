package models

type TransferBalanceRequestModel struct {
	ToUsername string  `json:"ToUsername"`
	Amount     float64 `json:"Amount"`
}
