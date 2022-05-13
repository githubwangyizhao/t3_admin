package models

type PlayerClientData struct {
	PlayerId int    `json:"playerId"`
	Id       string `json:"id"`
	Value    string `json:"value"`
	Times    int    `json:"times"`
}
