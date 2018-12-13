package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

const graphsegs string = " /\\-|+"

type coord [2]int

func main() {
	f, _ := os.Open("13.data")
	input := bufio.NewScanner(f)

	grid := make([][]byte, 0)
	carts := make(map[coord]Cart)
	y := 0
	for input.Scan() {
		line := input.Text()
		row := make([]byte, len(line))
		for x, _ := range line {
			b := line[x]
			if strings.Contains(graphsegs, string(b)) {
				row[x] = b
			} else {
				dir := cartDir(b)
				carts[coord{x, y}] = Cart{dir, 0}
				switch dir {
				case Up, Down:
					row[x] = '|'
				case Left, Right:
					row[x] = '-'
				default:
					panic(fmt.Sprintf("bad direction: %v", dir))
				}
			}

		}
		y++
		grid = append(grid, row)
	}
	f.Close()

	//crash := false
	for len(carts) > 1 {
		//printGrid(grid, carts)
		fmt.Println("carts", carts)
		_, carts = tickAll(grid, carts)
	}
	fmt.Println("last cart: ", carts)

}

func printGrid(grid [][]byte, carts map[coord]Cart) {
	for y, _ := range grid {
		for x, b := range grid[y] {
			if c, ok := carts[coord{x, y}]; ok {
				fmt.Printf("%v", string(dirToChar(c.dir)))
			} else {
				fmt.Printf("%c", b)
			}
		}
		fmt.Println()
	}
}

func dirToChar(d Dir) byte {
	switch d {
	case Right:
		return '>'
	case Left:
		return '<'
	case Up:
		return '^'
	case Down:
		return 'v'
	default:
		panic(fmt.Sprintf("bad direction: %v", d))
	}
}

type CoordSlice []coord

func (x CoordSlice) Len() int { return len(x) }
func (x CoordSlice) Less(i, j int) bool {
	if x[i][1] == x[j][1] {
		return x[i][0] < x[j][0]
	} else if x[i][1] < x[j][1] {
		return true
	}
	return false
}

func (x CoordSlice) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func tickAll(grid [][]byte, carts map[coord]Cart) (bool, map[coord]Cart) {
	var coords []coord
	for c := range carts {
		coords = append(coords, c)
	}

	sort.Sort(CoordSlice(coords))
	for _, c := range coords {
		cur, ok := carts[c]
		if ok {
			cart, nc := tickCart(grid, cur, c)
			if _, exists := carts[nc]; exists {
				fmt.Println("Crash at: ", nc)
				delete(carts, c)
				delete(carts, nc)
			} else {
				delete(carts, c)
				carts[nc] = cart
			}
		}
	}
	return false, carts
}

func tickCart(grid [][]byte, c Cart, xy coord) (Cart, coord) {
	//fmt.Println("tick cart: ", c, xy)
	switch c.dir {
	case Right:
		xy[0]++
	case Left:
		xy[0]--
	case Up:
		xy[1]--
	case Down:
		xy[1]++
	default:
		panic(fmt.Sprintf("bad direction:%v", c.dir))
	}
	//fmt.Println("new coord: ", xy)

	switch grid[xy[1]][xy[0]] {
	case '+':
		switch turnseq[c.nextturn] {
		case 'L':
			c.dir = (c.dir + 3) % 4
		case 'R':
			c.dir = (c.dir + 1) % 4
		}
		c.nextturn = (c.nextturn + 1) % 3
	default:
		c.dir = turn(grid[xy[1]][xy[0]], c.dir)
	}
	return c, xy
}

type Dir int

func turn(t byte, d Dir) Dir {
	switch t {
	case '\\':
		return -((d - 3) % 4)
	case '/':
		return -((d - 5) % 4)
	default:
		return d
	}
}

const (
	Up Dir = iota
	Right
	Down
	Left
)

const turnseq string = "LSR"

type Cart struct {
	dir      Dir
	nextturn int
}

func cartDir(b byte) Dir {
	switch b {
	case '>':
		return Right
	case '<':
		return Left
	case '^':
		return Up
	case 'v':
		return Down
	default:
		panic(fmt.Sprintf("bad cart: %v", string(b)))
	}
}
