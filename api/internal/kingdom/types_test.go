package kingdom

import (
	"encoding/json"
	"os"
	"testing"
)

var groupSimple = &Group{Animals: map[string]*Animal{
	"1":  {ID: "1", Name: "1", Gender: Male, Parents: map[Gender]string{}},
	"2":  {ID: "2", Name: "2", Gender: Female, Parents: map[Gender]string{}},
	"3":  {ID: "3", Name: "3", Gender: Male, Parents: map[Gender]string{}},
	"4":  {ID: "4", Name: "4", Gender: Female, Parents: map[Gender]string{}},
	"5":  {ID: "5", Name: "5", Gender: Male, Parents: map[Gender]string{Male: "1", Female: "2"}},
	"6":  {ID: "6", Name: "6", Gender: Female, Parents: map[Gender]string{Male: "3", Female: "4"}},
	"7":  {ID: "7", Name: "7", Gender: Male, Parents: map[Gender]string{Male: "5", Female: "6"}},
	"8":  {ID: "8", Name: "8", Gender: Female, COI: 0.25, Parents: map[Gender]string{Male: "5", Female: "7"}},
	"9":  {ID: "9", Name: "9", Gender: Male, Parents: map[Gender]string{Male: "7", Female: "8"}},
	"10": {ID: "10", Name: "10", Gender: Male, Parents: map[Gender]string{Male: "9", Female: "8"}},
}}

var groupSiblingInbreed = &Group{Animals: map[string]*Animal{
	"X": {ID: "X", Name: "X", Gender: Male, Parents: map[Gender]string{Male: "A", Female: "B"}},
	"A": {ID: "A", Name: "A", Gender: Male, Parents: map[Gender]string{Male: "C", Female: "D"}},
	"B": {ID: "B", Name: "B", Gender: Female, Parents: map[Gender]string{Male: "C", Female: "D"}},
	"C": {ID: "C", Name: "C", Gender: Male, Parents: map[Gender]string{Male: "E", Female: "F"}},
	"D": {ID: "D", Name: "D", Gender: Female, Parents: map[Gender]string{Male: "E", Female: "F"}},
	"E": {ID: "E", Name: "E", Gender: Male, Parents: map[Gender]string{}},
	"F": {ID: "F", Name: "F", Gender: Female, Parents: map[Gender]string{}},
}}

var groupDistantCommon = &Group{Animals: map[string]*Animal{
	"X": {ID: "X", Name: "X", Gender: Male, Parents: map[Gender]string{Male: "A", Female: "B"}},
	"A": {ID: "A", Name: "A", Gender: Male, Parents: map[Gender]string{Male: "C", Female: "D"}},
	"B": {ID: "B", Name: "B", Gender: Female, Parents: map[Gender]string{Male: "E", Female: "F"}},
	"C": {ID: "C", Name: "C", Gender: Male, Parents: map[Gender]string{Male: "G"}},
	"D": {ID: "D", Name: "D", Gender: Female, Parents: map[Gender]string{}},
	"E": {ID: "E", Name: "E", Gender: Male, Parents: map[Gender]string{Male: "G"}},
	"F": {ID: "F", Name: "F", Gender: Female, Parents: map[Gender]string{}},
	"G": {ID: "G", Name: "G", Gender: Male, Parents: map[Gender]string{}},
}}

func TestGroupFilter(t *testing.T) {
	filtered, err := groupSimple.FilterAnimalParents("9", 10)
	if err != nil {
		t.Fatal("filter fail", err)
	}
	js, _ := json.MarshalIndent(filtered.Animals, "", "  ")
	f, _ := os.Create("out.json")
	defer f.Close()
	f.Write(js)
}

func TestAnimalInbreedingInbreedingCoefficientBasic(t *testing.T) {
	coi, err := groupSimple.AnimalInbreedingCoefficient("9", 10)
	if err != nil {
		t.Fatal("inbreed calc fail", err)
	}
	t.Log("COI", coi)
}

func TestAnimalInbreedingInbreedingCoefficientInbreeding(t *testing.T) {
	coi, err := groupSiblingInbreed.AnimalInbreedingCoefficient("X", 10)
	if err != nil {
		t.Fatal("inbreed calc fail", err)
	}
	t.Log("COI", coi)
}
func TestAnimalInbreedingInbreedingCoefficientDistantCommon(t *testing.T) {
	coi, err := groupDistantCommon.AnimalInbreedingCoefficient("X", 10)
	if err != nil {
		t.Fatal("inbreed calc fail", err)
	}
	t.Log("COI", coi)
}

// func TestGroupTree2(t *testing.T) {
// 	tree, err := groupSimple.AnimalTree("9", 10)
// 	if err != nil {
// 		t.Fatal("tree build fail", err)
// 	}
// 	js, _ := json.MarshalIndent(tree.Root, "", "  ")
// 	f, _ := os.Create("out.json")
// 	defer f.Close()
// 	f.Write(js)
// }
