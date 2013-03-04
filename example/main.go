package main

import (
	"puzzle"
	"fmt"
)

func init() {
	fmt.Println("Sudoku:\n")
}

func main() {
	src := puzzle.StaticSudoku()
//	src := puzzle.RandomSudoku(20)
	dst, err := src.Solve()
	switch {
	case err != nil: fmt.Printf(err.Error())
	default: fmt.Printf("source:\n%v\n\nsolved:\n%v\n", src, dst)
	}
}
