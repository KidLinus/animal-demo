package model

type Dataset struct {
	ID      int    `json:"id"`
	Removed bool   `json:"removed"`
	Name    string `json:"name"`
}
