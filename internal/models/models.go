package models

type IPLocation struct {
	IP      string `json:"query"`
	Country string `json:"country"`
	City    string `json:"city"`
}
