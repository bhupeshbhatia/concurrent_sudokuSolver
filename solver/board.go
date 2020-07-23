package solver

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func convertToRow(stringSlice []string) Row {
	var row Row

	for _, currString := range stringSlice {
		currIntVal, err := strconv.Atoi(currString)
		if err != nil {
			currIntVal = 0
		}

		row = append(row, currIntVal)
	}

	return row
}

// CreateBoard = Get the sudoku puzzle from filename and input data into the board
func CreateBoard(filename string) (Sudoku, error) {
	var sudoku Sudoku

	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	for scanner.Scan() {
		replaceStrings := strings.Replace(scanner.Text(), ",", " ", -1)
		stringSlice := strings.Fields(replaceStrings)
		row := convertToRow(stringSlice)
		sudoku = append(sudoku, row)
	}

	return sudoku, nil
}
