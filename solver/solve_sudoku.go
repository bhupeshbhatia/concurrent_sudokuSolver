package solver

import (
	"errors"
	"sync"
	"sync/atomic"
)

//https://gobyexample.com/atomic-counters - mechanism for managing state in go is communication over channels
//sync/atomic package for atomic counters accessed by multiple goroutines
var globalCounter = new(int32)

// SolveSudoku solves an unsolved sudoku. Returns sudoku, number of iterations it took to solve and error
func SolveSudoku(sudokuBoard Sudoku) (Sudoku, bool, int, error) {
	// (Sudoku, int, error)
	//Step 1: create a copy
	sudokuCopy := sudokuBoard.MakeCopy()
	// fmt.Println("copy: ", sudokuCopy)

	//Keeps track of unfilled column
	unfilledCount := 0

	for {
		if sudokuCopy.ValidateSudoku() {
			break
		}

		unfilledCount = sudokuCopy.UnfilledCells()

		//
		atomic.AddInt32(globalCounter, 1)

		if *globalCounter >= 10000000 {
			break
		}

		//Next step is to browse the cells and use map reduce to fill cells with one potential number
		for rowIndex, row := range sudokuCopy {
			for colIndex, column := range row {
				if column == 0 {
					validCell := sudokuCopy.GetValidNumMap(rowIndex, colIndex)
					result := sudokuCopy.UpdateCellWithValidNum(validCell)

					if result == -1 {
						return sudokuCopy, sudokuCopy.ValidateSudoku(), int(*globalCounter), errors.New("No solve")
					}
				}
			}
		}

		// If the Sudoku is solved, exit out of the routine
		if sudokuCopy.ValidateSudoku() {
			break
		}

		// If no cells have been reduced, do not repeat. Start from valid numbers for each cell.
		//We pick a cell with least valid numbers
		if sudokuCopy.UnfilledCells() >= unfilledCount {
			mapOfPotentialCells := make(map[int]Cell)
			for rowIndex, row := range sudokuCopy {
				for colIndex, col := range row {
					if col == 0 {
						validCellAgain := sudokuCopy.GetValidNumMap(rowIndex, colIndex)
						valNumArray := validCellAgain.ValidNumbers.ListValidNumbers()
						lengthOfValNumArray := len(valNumArray)
						mapOfPotentialCells[lengthOfValNumArray] = validCellAgain
					}
				}
			}

			//Find out from all the cells which cells will require least amount of valid numbers
			var findOutLeastNumCell Cell
			rangeOfNum := []int{2, 3, 4, 5, 6, 7, 8, 9}
			for _, potentialValue := range rangeOfNum {
				if _, ok := mapOfPotentialCells[potentialValue]; ok {
					findOutLeastNumCell = mapOfPotentialCells[potentialValue]
					break
				}
			}

			//Pick all valid numbers and see if they fit in. This is going to be sped up using concurrency
			communicationChannel := make(chan Communication)
			waitGroup := new(sync.WaitGroup)

			for _, vNum := range findOutLeastNumCell.ValidNumbers.ListValidNumbers() {

				//A WaitGroup will help us wait for all goroutines to finish their work.
				waitGroup.Add(1)

				//Now let's call SolveSudoku recursively. In this case we will use goroutines to do this asynchronously
				go func(sudokuBoard Sudoku, rowIndex int, colIndex int, validNum int, waitGroup *sync.WaitGroup, c *chan Communication) {
					defer waitGroup.Done()

					sudokuCopy := sudokuBoard.MakeCopy()
					sudokuCopy.WriteToCell(rowIndex, colIndex, validNum)

					sudokuPotential, isSolved, _, err := SolveSudoku(sudokuCopy)
					*c <- Communication{
						PotentialSudoku: sudokuPotential,
						Solved:          isSolved,
						Err:             err,
					}
				}(sudokuCopy, findOutLeastNumCell.RowIndex, findOutLeastNumCell.ColIndex, vNum, waitGroup, &communicationChannel)
			}

			//Wait for threads to finish and close channel after
			go func(waitGroup *sync.WaitGroup, c chan Communication) {
				waitGroup.Wait()
				close(c)
			}(waitGroup, communicationChannel)

			//Collect all results
			for val := range communicationChannel {
				potentialSudoku := val.PotentialSudoku
				isSolved := val.Solved
				err := val.Err

				if isSolved {
					return potentialSudoku, isSolved, int(*globalCounter), err
				}

				if err.Error() != "No solve" {
					// not solved, but the guess is correct. try from beginning
					sudokuCopy = sudokuBoard.MakeCopy()
					break
				}
			}
		}
	}
	return sudokuCopy, sudokuCopy.ValidateSudoku(), int(*globalCounter), errors.New("Done")
}
