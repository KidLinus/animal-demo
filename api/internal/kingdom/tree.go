package kingdom

import "fmt"

type Tree struct {
	Root *TreeAnimal `json:"root"`
}

type TreeAnimal struct {
	Animal   *Animal       `json:"animal"`
	Children []*TreeAnimal `json:"children"`
}

func (group *Group) TreeAnimalParents(id string, maxDepth int, fillUnknown bool) (*Tree, error) {
	// Find the root node
	animal, ok := group.Animals[id]
	if !ok {
		return nil, fmt.Errorf("animal not found")
	}
	// Build the tree
	tree := &Tree{Root: &TreeAnimal{Animal: animal, Children: []*TreeAnimal{}}}
	search := []*TreeAnimal{tree.Root}
	for i := 0; i < maxDepth; i++ {
		if len(search) == 0 {
			break
		}
		var nextSearch []*TreeAnimal
		for _, animal := range search {
			if animal.Animal != nil {
				for _, parentID := range animal.Animal.Parents {
					if parent, ok := group.Animals[parentID]; ok {
						child := &TreeAnimal{Animal: parent, Children: []*TreeAnimal{}}
						animal.Children = append(animal.Children, child)
						nextSearch = append(nextSearch, child)
					}
				}
			}
			if fillUnknown {
				for idx, child := range animal.Children {
					if child == nil {
						child := &TreeAnimal{Children: []*TreeAnimal{}}
						animal.Children[idx] = child
						nextSearch = append(nextSearch, child)
					}
				}
			}
		}
		search = nextSearch
	}
	return tree, nil
}

func (group *Group) TreeAnimalChildren(id string, maxDepth int) (*Tree, error) {
	// Find the root node
	animal, ok := group.Animals[id]
	if !ok {
		return nil, fmt.Errorf("animal not found")
	}
	// Build a index of parent and their children
	parentChildren := map[string][]*Animal{}
	for _, animal := range group.Animals {
		for _, id := range animal.Parents {
			parentChildren[id] = append(parentChildren[id], animal)
		}
	}
	// fmt.Println("parentChildren", parentChildren)
	// Build the tree
	tree := &Tree{Root: &TreeAnimal{Animal: animal}}
	search := []*TreeAnimal{tree.Root}
	for i := 0; i < maxDepth; i++ {
		if len(search) == 0 {
			break
		}
		var nextSearch []*TreeAnimal
		for _, animal := range search {
			if animal.Animal != nil {
				for _, child := range parentChildren[animal.Animal.ID] {
					v := &TreeAnimal{Animal: child, Children: []*TreeAnimal{}}
					animal.Children = append(animal.Children, v)
					nextSearch = append(nextSearch, v)
				}
			}
		}
		search = nextSearch
	}
	return tree, nil
}
