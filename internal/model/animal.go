package model

import "time"

type Gender int

const (
	Male   Gender = 0
	Female Gender = 1
)

type Animal struct {
	ID      uint64
	Born    time.Time
	Died    *time.Time
	Name    string
	Gender  Gender
	Parents map[Gender]uint64
}
