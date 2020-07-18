# Sudoku Solver in Golang using concurrency

## Format

The format of the puzzles:

##### 

    5 _ _ 9 _ _ _ 7 _
    _ 6 _ _ _ _ 9 _ 4
    8 _ _ _ _ _ _ _ 5
    7 5 1 _ _ _ _ _ 8
    6 _ _ 2 _ _ 5 _ _
    _ 8 _ _ _ _ _ _ 1
    9 _ _ _ _ _ 3 _ _
    _ _ _ _ 4 _ _ _ _
    _ _ _ 5 _ 1 _ _ _ 

It reads the file and creates the sudoku board. Assuming you've saved the Sudoku in the same folder as puzzles you can run:

    go run main.go

It will provide the original sudoku and the solved solution:

    The Sudoku was solved successfully:
    
##### 

    [5 1 4 9 6 2 8 7 3]
    [3 6 7 1 8 5 9 2 4]
    [8 2 9 3 7 4 1 6 5]
    [7 5 1 4 9 6 2 3 8]
    [6 9 3 2 1 8 5 4 7]
    [4 8 2 7 5 3 6 9 1]
    [9 4 5 8 2 7 3 1 6]
    [1 3 8 6 4 9 7 5 2]
    [2 7 6 5 3 1 4 8 9]


## How does it work?

The solver uses backtracking and recursion to find a solution for the Sudoku puzzle. 

The solver (SolveSudoku function) tries to find the first empty cell (those with a `_`) and fills it with a number between 1-9 so long as it doesn't already occur in the corresponding row, column or 3x3 section. Then it moves to the next empty cell and inserts another number and so on.

If there is an issue (such as number already exists) and all possible values for the current cell has been tried then we move back to the previous one (which is called as backtracking). 

The same process is repeated until the board is finally solved.
