package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bhupeshbhatia/concurrent_sudokuSolver/solver"
)

func main() {
	fmt.Println("======Concurrent Sudoku=======")

	sudoku, err := solver.CreateBoard("puzzles/simple-1.txt")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	sudoku.Print()
	startTime := time.Now()
	sudoku, isSolved, iterations, err := solver.SolveSudoku(sudoku)
	elapsed := time.Since(startTime)
	fmt.Println("solved:", isSolved)
	fmt.Println("error:", err)
	fmt.Println("total iterations: ", iterations)
	fmt.Println("elapsed time: ", elapsed)
	sudoku.Print()
}
