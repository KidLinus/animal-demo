package model

import (
	"time"
)

type Animal struct {
	ID      int            `json:"id"`
	Dataset int            `json:"dataset"`
	Removed bool           `json:"removed"`
	Born    *time.Time     `json:"born"`
	Died    *time.Time     `json:"died"`
	Name    string         `json:"name"`
	Gender  Gender         `json:"gender"`
	Parents map[Gender]int `json:"parents"`
}
