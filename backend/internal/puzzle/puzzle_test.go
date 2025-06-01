package puzzle_test

import (
	"github.com/28Pollux28/webcross3d/internal/puzzle"
	"testing"
)

func TestPuzzle_GetClues(t *testing.T) {
	p := &puzzle.Puzzle{
		SizeX: 2,
		SizeY: 2,
		SizeZ: 2,
		Grid: [][][]bool{
			// Z=0
			{
				// Y=0
				{true, false},
				// Y=1
				{false, true},
			},
			// Z=1
			{
				// Y=0
				{false, true},
				// Y=1
				{true, false},
			},
		},
	}

	clues := p.GetClues()

	if len(clues) == 0 {
		t.Errorf("expected clues, got none")
	}

	// Example: check for a specific clue
	found := false
	for _, clue := range clues {
		if clue.Axis == "X" && clue.Coord1 == 0 && clue.Coord2 == 0 && clue.Count != nil && *clue.Count == 1 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected clue for X axis at (0,0) with count 1 not found")
	}
}
