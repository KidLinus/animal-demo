package kingdom

import (
	"fmt"
	"math"
)

type Group struct {
	Animals map[string]*Animal `json:"animal"`
}

func NewGroup(animals ...Animal) *Group {
	group := &Group{Animals: map[string]*Animal{}}
	for _, animal := range animals {
		group.Add(animal)
	}
	return group
}

func (group *Group) Add(v Animal) (existed bool) {
	_, existed = group.Animals[v.ID]
	group.Animals[v.ID] = &v
	return existed
}

func (group *Group) FilterAnimalParents(id string, maxDepth int) (*Group, error) {
	// Find the root node
	if _, ok := group.Animals[id]; !ok {
		return nil, fmt.Errorf("animal not found")
	}
	// Build the tree using only it's parents
	newGroup := NewGroup()
	search := []string{id}
	for i := 0; i < maxDepth; i++ {
		if len(search) == 0 {
			break
		}
		var searchNext []string
		for _, id := range search {
			if animal, ok := group.Animals[id]; ok {
				if _, ok := newGroup.Animals[id]; !ok {
					newGroup.Animals[id] = animal
					for _, parent := range animal.Parents {
						searchNext = append(searchNext, parent)
					}
				}
			}
		}
		search = searchNext
	}
	fmt.Println("group", len(group.Animals), "->", len(newGroup.Animals))
	return newGroup, nil
}

type AnimalInbreedingCoefficient struct {
	Result float64                           `json:"result"`
	Paths  []AnimalInbreedingCoefficientPath `json:"paths"`
}

type AnimalInbreedingCoefficientPath struct {
	Parent string   `json:"parent"`
	Path   []string `json:"path"`
	COI    float64  `json:"coi"`
	Result float64  `json:"result"`
}

func (groupOriginal *Group) AnimalInbreedingCoefficient(id string, maxDepth int) (*AnimalInbreedingCoefficient, error) {
	// Filter out only the nodes we need
	group, err := groupOriginal.FilterAnimalParents(id, maxDepth)
	if err != nil {
		return nil, err
	}
	// Make sure that the animal exists
	if _, ok := group.Animals[id]; !ok {
		return nil, fmt.Errorf("animal not found")
	}
	// Map animals parent dependencies
	dependencies := map[string][]string{}
	for _, animal := range group.Animals {
		for _, parent := range animal.Parents {
			dependencies[parent] = append(dependencies[parent], animal.ID)
		}
	}
	// log.Println("dependencies", dependencies)
	// Determine a execution order
	executionOrder, scheduled := []string{}, map[string]struct{}{}
	for {
		if len(scheduled) == len(group.Animals) { // All animals are scheduled
			break
		}
		var order []string
	animalLoop:
		for _, animal := range group.Animals {
			if _, ok := scheduled[animal.ID]; ok { // Skip if scheduled
				continue
			}
			for _, dependency := range dependencies[animal.ID] { // Check that all dependencies are resolved
				if _, scheduled := scheduled[dependency]; !scheduled {
					continue animalLoop
				}
			}
			order = append(order, animal.ID)
		}
		if len(order) == 0 { // Circular structure, can't calculate
			return nil, fmt.Errorf("impossible to resolve tree, circular dependencies")
		}
		executionOrder = append(executionOrder, order...)
		for _, id := range order {
			scheduled[id] = struct{}{}
		}
	}
	// log.Println("executionOrder", executionOrder)
	// Find all path intersections through the graph
	paths := map[string][][]string{id: {{}}}
	intersects := map[string][][]string{}
	for _, id := range executionOrder {
		animal := group.Animals[id]
		for _, parentID := range animal.Parents {
			if _, ok := group.Animals[parentID]; ok {
				// Add my paths to parent and check if it intersects something
				var newPaths [][]string
				for _, path := range paths[id] {
					for _, existingPath := range paths[parentID] { // Check intersects
						// fmt.Println("Intersect check!", parentID, ":", existingPath, "<->", path)
						if !sliceContainsAny(existingPath, path...) {
							intersectPath := existingPath
							for idx := range path {
								intersectPath = append(intersectPath, path[len(path)-idx-1])
							}
							intersects[parentID] = append(intersects[parentID], intersectPath)
							// fmt.Println("New intersect at id:", parentID, "path:", intersectPath)
						}
					}
					newPaths = append(newPaths, append(path, parentID))
				}
				paths[parentID] = append(paths[parentID], newPaths...) // Add new paths leading here
			}
		}
	}
	// log.Println("intersects", intersects)
	// Calculate COI
	res := &AnimalInbreedingCoefficient{Paths: []AnimalInbreedingCoefficientPath{}}
	for parentID, paths := range intersects {
		var parentCOI float64
		if parent, ok := group.Animals[parentID]; ok && parent.COI != nil {
			parentCOI = *parent.COI
		}
		for idx, path := range paths {
			coi := math.Pow(0.5, float64(len(path))*(1+parentCOI))
			// fmt.Println("Coefficient", path, "Fx", (1 + parentCOI), " = ", coi)
			res.Paths = append(res.Paths, AnimalInbreedingCoefficientPath{Parent: parentID, Path: paths[idx], COI: parentCOI, Result: coi})
			res.Result += coi
		}
	}
	// Build result
	return res, nil
}

func (groupOriginal *Group) Merge(groupOther *Group) *Group {
	group := NewGroup()
	for _, v := range groupOriginal.Animals {
		group.Add(*v)
	}
	for _, v := range groupOther.Animals {
		group.Add(*v)
	}
	return group
}
