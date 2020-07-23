package solver

// Sudoku is a slice of Rows. Used for representing a sudoku puzzle, solved or otherwise
type Sudoku []Row

// ValidNumbers - A hashmap to keep track of which numbers are valid to be filled in a column.
// The hashmap ALWAYS contains numbers 1 to 9 as the key. If a number is valid to be filled in a cell,
// the corresponding value is marked as true, else false.
type ValidNumbers map[int]bool

// Row is a slice of integers. Used for representing a row, column or a bounded box of a sudoku of length 9
type Row []int

// RowLength is a constant that represents the length of a sudoku row
const RowLength int = 9

// ColLength is a constant that represents the length of a sudoku column
const ColLength int = 9

// Cell is a structure that contains a sudoku position (RowID, ColID) and the eligible numbers that can be filled in it
type Cell struct {
	RowIndex     int
	ColIndex     int
	ValidNumbers ValidNumbers
}

// Channel is a structure that is used to communicate the results of a solve run that is run concurrently
type Communication struct {
	PotentialSudoku Sudoku
	Solved          bool
	Err             error
}
