package kingdom

type Animal struct {
	ID      string
	Name    string
	Gender  Gender
	COI     float64
	Parents map[Gender]string
}

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)
