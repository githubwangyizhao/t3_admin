package models

type PlayerFunction struct {
	PlayerId   int `json:"player_id"`
	FunctionId int `json:"function_id"`
	State      int `json:"state"`
	GetState   int `json:"get_state"`
	Time       int `json:"time"`
}
