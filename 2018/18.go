package main

import (
	"bufio"
	"bytes"
	"os"
	"fmt"
)

const boardmax = 50

type board [boardmax][boardmax]string
type nmap map[string]int

func (b board) String() string {
	var buf bytes.Buffer
	buf.WriteByte('\n')
	for _, col := range b {
		for _, v := range col {
			buf.WriteString(v)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (b board) countNeighbors(xs, ys int) nmap {
	m := make(nmap)
	for x := xs - 1; x <= xs+1; x++ {
		for y := ys - 1; y <= ys+1; y++ {
			if (x == xs && y == ys) ||
				x < 0 || x == boardmax ||
				y < 0 || y == boardmax {
				continue
			}
			m[b[x][y]]++
		}
	}
	return m
}

func (b board) getResource() int {
	cnt := make(nmap)
	for _, a := range b {
		for _, v := range a {
			cnt[v]++
		}
	}
	return cnt["|"] * cnt["#"]
}
	

func (b board) iterate(nb *board) {
	for y := 0; y < boardmax; y++ {
		for x := 0; x < boardmax; x++ {
			nm := b.countNeighbors(y, x)
			var changed bool
			switch b[y][x] {
			case ".":
				if nm["|"] >= 3 {
					nb[y][x] = "|"
					changed = true
				}
			case "|":
				if nm["#"] >= 3 {
					nb[y][x] = "#"
				changed = true}

			case "#":
				if nm["#"] >= 1 && nm["|"] >= 1 {
					nb[y][x] = "#"
				} else {
					nb[y][x] = "."
					changed = true
				}
			default:
				panic("bad acre")
			}
			if !changed {
				nb[y][x] = b[y][x]
			}
		}
	}
}

func main() {
	f, _ := os.Open("18.data")
	input := bufio.NewScanner(f)

	var a, b board
	y := 0
	for input.Scan() {
		for x, v := range input.Text() {
			a[y][x] = string(v)
		}
		y++
	}
	b = a
	active, prev := &a, &b

	fmt.Println("prev: ", *prev)
	for i := 1; i <= 10; i++ {	
		prev.iterate(active)
		fmt.Println("active: ", i, *active)
		prev, active = active, prev
	}
	fmt.Println("resources ", prev.getResource())
}
