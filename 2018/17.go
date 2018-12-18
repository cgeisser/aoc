package main

import (
	"bufio"
	//"bytes"
	"fmt"
	"os"
	//"sort"
	//	"strings"
)

type coord [2]int

func main() {
	f, _ := os.Open(os.Args[1])
	input := bufio.NewScanner(f)

	bo := new(Board)
	bo.clay = make(map[coord]bool)
	for input.Scan() {
		var start, x, y, xlen, ylen, skip int
		fmt.Sscanf(input.Text(), "%c=%d, %c=%d..%d", &start, &x, &skip, &y, &ylen)
		xlen = 1
		ylen = ylen - y + 1
		if start == int('y') {
			x, y = y, x
			xlen, ylen = ylen, xlen
		}
		bo.addclay(coord{x, y}, xlen, ylen)
	}
	graph := bo.getGraph()
	bo.printGrid(graph)

	fmt.Println("now filled.")
	fillGraph(graph[coord{500, 0}])

	bo.printGrid(graph)
	fmt.Println("reachable ", graph.reachable())
	f.Close()
}

type Board struct {
	clay map[coord]bool

	topleft, bottomright coord
}

func (b *Board) addclay(xy coord, xlen, ylen int) {
	//fmt.Println("adding ", xy, xlen, ylen)
	if (b.topleft == coord{0, 0} && b.bottomright == coord{0, 0}) {
		b.topleft, b.bottomright = xy, xy
	}

	if xy[0]-1 < b.topleft[0] {
		b.topleft[0] = xy[0] - 1
	}
	if xy[0]+xlen > b.bottomright[0] {
		b.bottomright[0] = xy[0] + xlen
	}
	if xy[1]-1 < b.topleft[1] {
		b.topleft[1] = xy[1] - 1
	}
	if xy[1]+ylen > b.bottomright[1] {
		b.bottomright[1] = xy[1] + ylen
	}
	//fmt.Println("tl ", b.topleft)
	//fmt.Println("br ", b.bottomright)
	for x := 0; x < xlen; x++ {
		for y := 0; y < ylen; y++ {
			b.clay[coord{xy[0] + x, xy[1] + y}] = true
		}
	}

}

func (b Board) printGrid(n NodeGraph) {
	fmt.Println("tl ", b.topleft)
	fmt.Println("br ", b.bottomright)
	for y := b.topleft[1]; y <= b.bottomright[1]; y++ {
		for x := b.topleft[0]; x <= b.bottomright[0]; x++ {
			if b.clay[coord{x, y}] {
				fmt.Printf("X")
			} else if n[coord{x, y}].enc {
				fmt.Printf("~")
			} else if n[coord{x, y}].v {
				fmt.Printf("|")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}

}

func fillGraph(n *Node) bool {
	//fmt.Println(" filling: ", n)
	if n.v || n.bottom {
		return n.enc
	}

	enc := true
	if n.d != nil {
		//fmt.Println(" d", n.d)
		if !n.d.v {
			enc = fillGraph(n.d) && enc
		}
	}
	if n.d == nil || n.d.enc {
		if n.l != nil {
			//fmt.Println(" l", n.l)
			if !n.l.v {
				enc = fillGraph(n.l) && enc
			}

		}
		n.v = true
		if n.r != nil {
			//fmt.Println(" r", n.r)
			if !n.r.v {
				enc = fillGraph(n.r) && enc
			}
		}
	}
	n.enc = enc

	n.v = true
	return n.enc
}

func (ng NodeGraph) reachable() int {
	r := 0
	for _, x := range ng {
		if x.v || x.enc {
			r++
		}
	}
	return r
}

func (b Board) getGraph() NodeGraph {
	g := make(NodeGraph)

	for y := b.topleft[1]; y <= b.bottomright[1]; y++ {
		for x := b.topleft[0]; x <= b.bottomright[0]; x++ {
			if b.clay[coord{x, y}] {
				continue
			}
			n := &Node{}
			n.xy = coord{x, y}
			g[coord{x, y}] = n
			if x > b.topleft[0] {
				if x <= b.bottomright[0] {
					if !b.clay[coord{x - 1, y}] {
						n.l = g[coord{x - 1, y}]
						g[coord{x - 1, y}].r = n
					}
				}
			}
			if y > b.topleft[1] {
				if !b.clay[coord{x, y - 1}] {
					g[coord{x, y - 1}].d = n
				}
				if y == b.bottomright[1] {
					n.bottom = true
				}
			}
		}
	}
	return g
}

type NodeGraph map[coord]*Node

type Node struct {
	xy             coord
	d, l, r        *Node
	bottom, enc, v bool
}

func (n *Node) String() string {
	var l, r, d string
	l, r, d = "nil", "nil", "nil"
	if n.l != nil {
		l = fmt.Sprintf("%v", n.l.xy)
	}
	if n.r != nil {
		r = fmt.Sprintf("%v", n.r.xy)
	}
	if n.d != nil {
		d = fmt.Sprintf("%v", n.d.xy)
	}
	return fmt.Sprintf("%v: %v %v %v b:%v v:%v enc:%v", n.xy, l, r, d, n.bottom, n.v, n.enc)
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
