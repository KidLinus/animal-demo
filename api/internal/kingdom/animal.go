package kingdom

import (
	"context"
	"time"
)

type Animal struct {
	ID        string              `json:"id"`
	IDs       map[string][]string `json:"ids"`
	Created   time.Time           `json:"created"`
	Updated   time.Time           `json:"updated"`
	Removed   bool                `json:"removed"`
	Name      string              `json:"name"`
	Nicknames []string            `json:"nicknames"`
	Gender    Gender              `json:"gender"`
	Country   *string             `json:"country"`
	COI       *float64            `json:"coi"`
	Born      *time.Time          `json:"born"`
	Deceased  *time.Time          `json:"deceased"`
	Parents   map[Gender]string   `json:"parents"`
}

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type AnimalDatabase interface {
	Get(context.Context, string) (*Animal, error)
	List(context.Context, AnimalFilter) (AnimalIterator, error)
	Store(context.Context, Animal) (bool, error)
	Update(context.Context, string, AnimalUpdate)
	Remove(context.Context, string) (bool, error)
}

type AnimalUpdate struct {
	Created   *time.Time
	Updated   *time.Time
	Removed   *bool
	Name      *string
	Nicknames *AnimalUpdateNicknames
	Gender    *Gender
	Country   *AnimalUpdateCountry
	COI       *AnimalUpdateCOI
	Born      *AnimalUpdateBorn
	Deceased  *AnimalUpdateDeceased
	Parents   *AnimalUpdateParents
}

type AnimalUpdateNicknames struct {
	Set    []string
	Ensure []string
	Remove []string
}

type AnimalUpdateCountry struct {
	Set   *string
	Clear bool
}

type AnimalUpdateCOI struct {
	Set   *float64
	Clear bool
}

type AnimalUpdateBorn struct {
	Set   time.Time
	Clear bool
}

type AnimalUpdateDeceased struct {
	Set   time.Time
	Clear bool
}

type AnimalUpdateParents struct {
	Set    map[Gender]string
	Ensure map[Gender]string
	Clear  []Gender
}

type AnimalFilter struct {
	Query *string
	IDs   []string
	Limit *int
}

type AnimalIterator interface{ Next(*Animal) error }
