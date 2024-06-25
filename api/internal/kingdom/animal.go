package kingdom

import (
	"context"
	"time"
)

type Animal struct {
	ID        string
	IDs       map[string][]string
	Created   time.Time
	Updated   time.Time
	Removed   bool
	Name      string
	Nicknames []string
	Gender    Gender
	Country   *string
	COI       *float64
	Born      *time.Time
	Deceased  *time.Time
	Parents   map[Gender]string
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
