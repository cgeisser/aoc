package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

const boardmax = 50

type board [boardmax][boardmax]string
type nmap [3]int

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

func (n *nmap) inc(s string) {
	var key int
	switch s {
	case ".":
		key = CLEAR
	case "|":
		key = TREE
	case "#":
		key = YARD
	default:
		panic("bad key")
	}
	n[key]++
}

func (b board) countNeighbors(xs, ys int) nmap {
	var m nmap
	for x := xs - 1; x <= xs+1; x++ {
		for y := ys - 1; y <= ys+1; y++ {
			if (x == xs && y == ys) ||
				x < 0 || x == boardmax ||
				y < 0 || y == boardmax {
				continue
			}
			m.inc(b[x][y])
		}
	}
	return m
}

func (b board) getResource() int {
	var cnt nmap
	for _, a := range b {
		for _, v := range a {
			cnt.inc(v)
		}
	}
	return cnt[TREE] * cnt[YARD]
}

const (
	CLEAR = iota
	TREE
	YARD
)

func (b board) iterate(nb *board) {
	for y := 0; y < boardmax; y++ {
		for x := 0; x < boardmax; x++ {
			nm := b.countNeighbors(y, x)
			var changed bool
			switch b[y][x] {
			case ".":
				if nm[TREE] >= 3 {
					nb[y][x] = "|"
					changed = true
				}
			case "|":
				if nm[YARD] >= 3 {
					nb[y][x] = "#"
					changed = true
				}

			case "#":
				if nm[YARD] >= 1 && nm[TREE] >= 1 {
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
	for i := 1; i <= 100000; i++ {
		prev.iterate(active)
		//		fmt.Println("active: ", i, *active)
		prev, active = active, prev
		if i%1000 == 0 {
			fmt.Println("resources ", i, prev.getResource())
		}
	}

}
