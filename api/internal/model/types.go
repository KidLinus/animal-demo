package model

type Gender int

const (
	Male   Gender = 0
	Female Gender = 1
)

func (g Gender) String() string {
	if g == Male {
		return "male"
	}
	return "female"
}
