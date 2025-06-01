package puzzle

var puzzleStore = map[string]*Puzzle{
	"1": {
		ID:     "1",
		Name:   "Intro Puzzle",
		Author: "Admin",
		Grid: [][][]bool{
			{{true, false}, {false, false}},
			{{false, false}, {true, true}},
		},
		Lives: 5,
		SizeX: 2,
		SizeY: 2,
		SizeZ: 2,
	},
}

func GetPuzzle(id string) (*Puzzle, bool) {
	p, ok := puzzleStore[id]
	return p, ok
}

func GetPuzzles() []*Puzzle {
	var result []*Puzzle
	for _, p := range puzzleStore {
		result = append(result, p)
	}
	return result
}
