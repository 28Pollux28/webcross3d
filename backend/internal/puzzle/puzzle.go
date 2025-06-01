package puzzle

import (
	"errors"
)

type Puzzle struct {
	ID     string
	Name   string
	Author string
	Grid   [][][]bool
	Lives  int
	SizeX  int
	SizeY  int
	SizeZ  int
}

type Clue struct {
	Axis   string
	Coord1 int
	Coord2 int
	Count  *int
	Split  SplitType
}

type SplitType string

var (
	NoSplit    SplitType = "NoSplit"
	Split2     SplitType = "Split2"
	Split3Plus SplitType = "Split3Plus"
)

var (
	ErrOutOfBounds    = errors.New("coordinates out of bounds")
	ErrIncorrectVoxel = errors.New("incorrect voxel destruction")
)

func (p *Puzzle) ValidateVoxel(x, y, z int) (correct bool, err error) {
	// out of bounds check
	if x < 0 || x >= p.SizeX || y < 0 || y >= p.SizeY || z < 0 || z >= p.SizeZ {
		return false, ErrOutOfBounds
	}

	// Check the voxel
	if p.Grid[z][y][x] {
		return false, ErrIncorrectVoxel
	}

	return true, nil
}

func (p *Puzzle) IsComplete(removed map[[3]int]bool) bool {
	for z := 0; z < p.SizeZ; z++ {
		for y := 0; y < p.SizeY; y++ {
			for x := 0; x < p.SizeX; x++ {
				if !p.Grid[z][y][x] && !removed[[3]int{x, y, z}] {
					return false
				}
			}
		}
	}
	return true
}

func (p *Puzzle) GetClues() []Clue {
	var clues []Clue

	// X-axis clues
	for y := 0; y < p.SizeY; y++ {
		for z := 0; z < p.SizeZ; z++ {
			count, split := countCluesOnAxis(p, "X", y, z)
			clues = append(clues, Clue{Axis: "X", Coord1: y, Coord2: z, Count: &count, Split: split})
		}
	}

	// Y-axis clues
	for x := 0; x < p.SizeX; x++ {
		for z := 0; z < p.SizeZ; z++ {
			count, split := countCluesOnAxis(p, "Y", x, z)
			clues = append(clues, Clue{Axis: "Y", Coord1: x, Coord2: z, Count: &count, Split: split})

		}
	}

	// Z-axis clues
	for x := 0; x < p.SizeX; x++ {
		for y := 0; y < p.SizeY; y++ {
			count, split := countCluesOnAxis(p, "Z", x, y)
			clues = append(clues, Clue{Axis: "Z", Coord1: x, Coord2: y, Count: &count, Split: split})
		}
	}
	return clues
}

func countCluesOnAxis(p *Puzzle, axis string, i1, i2 int) (int, SplitType) {
	count := 0
	split := NoSplit
	var length int
	var getVoxel func(idx int) bool

	switch axis {
	case "X":
		length = p.SizeX
		getVoxel = func(idx int) bool { return p.Grid[i2][i1][idx] }
	case "Y":
		length = p.SizeY
		getVoxel = func(idx int) bool { return p.Grid[i2][idx][i1] }
	case "Z":
		length = p.SizeZ
		getVoxel = func(idx int) bool { return p.Grid[idx][i2][i1] }
	default:
		return 0, NoSplit
	}

	for idx := 0; idx < length; idx++ {
		if !getVoxel(idx) {
			if count > 0 {
				switch split {
				case NoSplit:
					split = Split2
				case Split2:
					split = Split3Plus
				}
			}
		} else {
			count++
		}
	}
	return count, split
}
