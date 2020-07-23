package solver

import "fmt"

// Print prints sudoku
func (s Sudoku) Print() {
	for i, col := range s {
		fmt.Println(i, col)
	}
}

//MakeCopy creates a deep copy using array of array of int for copying each element
//of a multi-dimensional slice == used race detector flag of golang
func (s Sudoku) MakeCopy() Sudoku {
	sudoku := make(Sudoku, 0)
	//Create channel to make copy
	copyFinished := make(chan struct{})

	go func() {
		for _, rowFromSudokuBoardS := range s {
			row := make(Row, 0)
			for _, colFromSudokuBoardS := range rowFromSudokuBoardS {
				row = append(row, colFromSudokuBoardS)
			}
			sudoku = append(sudoku, row)
		}
		copyFinished <- struct{}{}
	}()
	<-copyFinished
	return sudoku
}

func (r Row) mapValidNum() ValidNumbers {
	//map with all numbers from 1 to 9
	validNumMap := make(ValidNumbers)
	for i := 1; i <= 9; i++ {
		validNumMap[i] = true
	}
	return validNumMap
}

//isRowFilledWithNum validates whether a row contains all numbers without 0
func (r Row) isRowFilledWithNum() bool {
	validNum := r.mapValidNum()
	for _, val := range validNum {
		if val {
			return false
		}
	}
	return true
}

//Validates that row contains non repeating numbers
func (r Row) rowHasNonRepeatingNum() bool {
	validMap := make(map[int]int)

	for _, column := range r {
		validMap[column] = validMap[column] + 1

		if validMap[column] > 1 {
			return false
		}
	}
	return true
}

// Get the index inside the 3 x 3 square box. Divided based on rows and then columns
func getSquareBoxIndex(rowIndex, colIndex int) int {
	squareIndex := 0
	switch {
	case rowIndex < 3:
		switch {
		case colIndex < 3:
			squareIndex = 0
			break
		case colIndex < 6:
			squareIndex = 1
			break
		case colIndex < 9:
			squareIndex = 2
			break
		}
		break

	case rowIndex < 6:
		switch {
		case colIndex < 3:
			squareIndex = 3
			break
		case colIndex < 6:
			squareIndex = 4
			break
		case colIndex < 9:
			squareIndex = 5
			break
		}
		break

	case rowIndex < 9:
		switch {
		case colIndex < 3:
			squareIndex = 6
			break
		case colIndex < 6:
			squareIndex = 7
			break
		case colIndex < 9:
			squareIndex = 8
			break
		}
		break
	}
	return squareIndex
}

func squareBoxes(rowIndex, columnIndex int) (int, int, int, int) {
	var lowerRow, upperRow, lowerColumn, upperColumn int

	switch {
	case rowIndex <= 2 && columnIndex <= 2:
		lowerRow = 0
		upperRow = 2
		lowerColumn = 0
		upperColumn = 2
		break
	case rowIndex <= 5 && columnIndex <= 2:
		lowerRow = 3
		upperRow = 5
		lowerColumn = 0
		upperColumn = 2
		break
	case rowIndex <= 8 && columnIndex <= 2:
		lowerRow = 6
		upperRow = 8
		lowerColumn = 0
		upperColumn = 2
		break
	case rowIndex <= 2 && columnIndex <= 5:
		lowerRow = 0
		upperRow = 2
		lowerColumn = 3
		upperColumn = 5
		break
	case rowIndex <= 5 && columnIndex <= 5:
		lowerRow = 3
		upperRow = 5
		lowerColumn = 3
		upperColumn = 5
		break
	case rowIndex <= 8 && columnIndex <= 5:
		lowerRow = 6
		upperRow = 8
		lowerColumn = 3
		upperColumn = 5
		break
	case rowIndex <= 2 && columnIndex <= 8:
		lowerRow = 0
		upperRow = 2
		lowerColumn = 6
		upperColumn = 8
		break
	case rowIndex <= 5 && columnIndex <= 8:
		lowerRow = 3
		upperRow = 5
		lowerColumn = 6
		upperColumn = 8
		break
	case rowIndex <= 8 && columnIndex <= 8:
		lowerRow = 6
		upperRow = 8
		lowerColumn = 6
		upperColumn = 8
		break
	}
	return lowerRow, upperRow, lowerColumn, upperColumn
}

// ValidateSudoku checks whether a sudoku is solved or not. Rules for solving sudoku:
// 1. No cells left with number 0
// 2. All cells have numbers 1 to 9
// 3. No repeating numbers within a row, column or 3 x 3 box
func (s Sudoku) ValidateSudoku() bool {
	//make map of columns
	columns := make(map[int]Row)

	//make map of 3 x 3 squares
	threeByThreeSq := make(map[int]Row)

	// Go through sudoku board (s) to find out rows, columns and 3 x 3 squares
	for rowIndex, row := range s {
		if !(row.isRowFilledWithNum() && row.rowHasNonRepeatingNum()) {
			return false
		}

		for colIndex, colVal := range row {
			// collect column values belonging to the same column Index in a separate Row
			columns[colIndex] = append(columns[colIndex], colVal)

			// collect column values belonging to the same three by three square into a separate Row
			squareBoxIndex := getSquareBoxIndex(rowIndex, colIndex)
			threeByThreeSq[squareBoxIndex] = append(threeByThreeSq[squareBoxIndex], colVal)
		}
	}

	if len(columns) > 0 {
		for _, row := range columns {
			if !(row.isRowFilledWithNum() && row.rowHasNonRepeatingNum()) {
				return false
			}
		}
	}

	if len(threeByThreeSq) > 0 {
		for _, row := range threeByThreeSq {
			if !(row.isRowFilledWithNum() && row.rowHasNonRepeatingNum()) {
				return false
			}
		}
	}
	return true
}

//UnfilledCells gets the total number of unfilled cells in a sudoku board
func (s Sudoku) UnfilledCells() int {
	unfilledCount := 0
	finished := make(chan struct{})

	go func() {
		for _, row := range s {
			for _, column := range row {
				if column == 0 {
					unfilledCount++
				}
			}
		}
		finished <- struct{}{}
	}()
	<-finished
	return unfilledCount
}

//GetValidNumMap gets a map of valid numbers for 3 x 3 positions on the sudoku map
func (s Sudoku) GetValidNumMap(rowIndex, colIndex int) Cell {
	var column, squareBox Row

	// make a standard map with all numbers from 1 to 9 as eligible
	valMap := make(ValidNumbers)
	for i := 1; i <= 9; i++ {
		valMap[i] = true
	}

	// first get the Row, column and square box corresponding to the position
	//Row
	row := s[rowIndex]

	//column ID
	for _, Row := range s {
		for cID, col := range Row {
			if colIndex == cID {
				column = append(column, col)
			}
		}
	}

	//square box
	lowerRow, upperRow, lowerColumn, upperColumn := squareBoxes(rowIndex, colIndex)

	for rowID, Row := range s {
		for colID, col := range Row {
			if (rowID >= lowerRow && rowID <= upperRow) && (colID >= lowerColumn && colID <= upperColumn) {
				squareBox = append(squareBox, col)
			}
		}
	}

	//Check for already present numbers and eliminate them
	for _, col := range row {
		if col != 0 {
			valMap[col] = false
		}
	}

	for _, col := range column {
		if col != 0 {
			valMap[col] = false
		}
	}

	for _, col := range squareBox {
		if col != 0 {
			valMap[col] = false
		}
	}

	return Cell{
		RowIndex:     rowIndex,
		ColIndex:     colIndex,
		ValidNumbers: valMap,
	}
}

//UpdateCellWithValidNum = pdates a specific Cell if the Cell contains only one eligible number.
/*
	Rules
	If no numbers are eligible, -1 is returned.
	If more than one number is elgible, 0 is returned.
	If only one number is elgible, the cell is filled with the number and this number is returned.
*/
func (s Sudoku) UpdateCellWithValidNum(v Cell) int {
	valMap := v.ValidNumbers
	rowIndex := v.RowIndex
	colIndex := v.ColIndex

	validNumberCount := 0
	incorrectNumCount := 0
	valNum := 0

	for key, val := range valMap {
		if val {
			validNumberCount++
			valNum = key
		} else {
			incorrectNumCount++
		}
	}

	// If exactly one number is eligible, send the corresponding number
	if validNumberCount > 1 {
		valNum = 0
	}

	// If no numbers are eligible, then send a different signal -1
	if incorrectNumCount == 9 {
		valNum = -1
	}

	if valNum >= 1 && valNum <= 9 {
		s.WriteToCell(rowIndex, colIndex, valNum)
	}

	return valNum
}

//WriteToCell fills up a specific cell- rowIndex, colIndex with the value
func (s Sudoku) WriteToCell(rowIndex, colIndex, val int) {
	completed := make(chan struct{})

	go func(s Sudoku) {
		s[rowIndex][colIndex] = val
		completed <- struct{}{}
	}(s)
	<-completed
}

//ListValidNumbers creates a list of valid number map
func (v ValidNumbers) ListValidNumbers() []int {
	var valNumArray []int
	for key, val := range v {
		if val {
			valNumArray = append(valNumArray, key)
		}
	}
	return valNumArray
}
