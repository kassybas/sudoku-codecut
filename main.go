package main

import (
	"fmt"
	"sudoku/genetic"
)

func main() {
	//field:=[]byte{3,0,1,0,0,0,0,0,2}
	field := []byte{3, 0, 1, 0, 0, 0, 0, 0, 2, 13, 0, 11, 0, 0, 0, 0, 0, 12, 14, 15, 0, 0, 0, 0, 0, 20}
	fmt.Println("Sudoku solver")
	genetic.Solve(field)

}
