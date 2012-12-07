package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(`Usage: sudoku <puzzle file>
    where puzzle file contains one or more puzzles in 81 character format.
Each puzzle must be on a separate line and must have 81 characters.
Lines less than 81 characters are ignored.  Any character other than 1-9
(such as 0, space, or .) can be used to represent a blank square.`)
		return
	}
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	var solved, failed int
	for _, puzzle := range strings.Split(string(b), "\n") {
		if len(puzzle) < 81 {
			continue
		}
		puzzle = puzzle[:81] // ignore anything else on line
		printGrid("\nPuzzle:", puzzle)
		if s := solve(puzzle); s == "" {
			fmt.Println("no solution")
			failed++
		} else {
			printGrid("Solved:", s)
			solved++
		}
	}
	fmt.Println("\nPuzzles solved:", solved)
	fmt.Println("Failed to solve: ", failed)
}

// print grid (with title) from 81 character string
func printGrid(title, s string) {
	fmt.Println(title)
	for r, i := 0, 0; r < 9; r, i = r+1, i+9 {
		fmt.Printf("%c %c %c | %c %c %c | %c %c %c\n", s[i], s[i+1], s[i+2],
			s[i+3], s[i+4], s[i+5], s[i+6], s[i+7], s[i+8])
		if r == 2 || r == 5 {
			fmt.Println("------+-------+------")
		}
	}
}

// extracts 81 character sudoku string
func text(s [][]int) string {
	b := make([]byte, len(s))
	for _, r := range s {
		b[r[0]] = byte(r[1]%9) + '1'
	}
	return string(b)
}

// solve puzzle in 81 character string format.
// if solved, result is 81 character string.
// if not solved, result is the empty string.
func solve(u string) string {
	// construct an dlx object with 324 constraint columns.
	// other than the number 324, this is not specific to sudoku.
	d := New(324)
	// now add constraints that define sudoku rules.
	for r, i := 0, 0; r < 9; r++ {
		for c := 0; c < 9; c, i = c+1, i+1 {
			b := r/3*3 + c/3
			n := int(u[i] - '1')
			if n >= 0 && n < 9 {
				d.AddRow([]int{i, 81 + r*9 + n, 162 + c*9 + n,
					243 + b*9 + n})
			} else {
				for n = 0; n < 9; n++ {
					d.AddRow([]int{i, 81 + r*9 + n, 162 + c*9 + n,
						243 + b*9 + n})
				}
			}
		}
	}
	// run dlx.  not sudoku specific.
	// extract the sudoku-specific 81 character result from the dlx solution.
	return text(d.Search())
}
