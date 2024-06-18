package model

import (
	"time"
)

type Animal struct {
	ID      int
	Dataset int
	Removed bool
	Born    time.Time
	Died    *time.Time
	Name    string
	Gender  Gender
	Parents map[Gender]int
}
