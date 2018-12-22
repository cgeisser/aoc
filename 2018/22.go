package main

import (
	"fmt"
)

const depth = 9465
const tx = 13
const ty = 704

type rtype int
const (
	ROCK = iota
	WET
	NARROW
)

type region struct {
	rt int
	el int
}

func main() {
	var b [tx+1][ty+1]region

	risk := 0
	for x, c := range b {
		for y, _ := range c {
			gi := 0
			if y == 0 {
				gi = x*16807
			} else 	if x == 0 {
				gi = y*48271
			} else if x == tx && y == ty {
				gi = 0
			} else {
				gi = b[x-1][y].el * b[x][y-1].el
			}
			b[x][y].el = (gi + depth) % 20183
			b[x][y].rt = b[x][y].el % 3
			risk+= b[x][y].rt
		}
	}
	fmt.Println("risk: ", risk)
}

	
			
			
