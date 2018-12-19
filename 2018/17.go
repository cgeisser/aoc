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
		_, ok := fmt.Sscanf(input.Text(), "%c=%d, %c=%d..%d", &start, &x, &skip, &y, &ylen)
		if ok != nil {
			panic(fmt.Sprintf("%v %v", ok, input.Text()))
		}
		xlen = 1
		ylen = ylen - y + 1
		if start == int('y') {
			x, y = y, x
			xlen, ylen = ylen, xlen
		}
		bo.addclay(coord{x, y}, xlen, ylen)
	}
	//os.Exit(0)
	bo.getGraph()
	//bo.printGrid()

	fmt.Println("now filled.")
	bo.fillGraph(bo.ng[coord{500, 0}])

	bo.printGrid()
	r, w := bo.reachable()
	fmt.Println("reachable ", r, w)
	f.Close()
}

type Board struct {
	clay                 map[coord]bool
	ng                   NodeGraph
	topleft, bottomright coord
	miny                 int
}

func (b *Board) addclay(xy coord, xlen, ylen int) {
	//fmt.Println("adding ", xy, xlen, ylen)
	if b.miny == 0 || xy[1] < b.miny {
		b.miny = xy[1]
	}
	if (b.topleft == coord{0, 0} || b.bottomright == coord{0, 0}) {
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

func (b Board) printGrid() {
	fmt.Println("tl ", b.topleft)
	fmt.Println("br ", b.bottomright)
	fmt.Println("miny ", b.miny)
	for y := b.topleft[1]; y <= b.bottomright[1]; y++ {
		for x := b.topleft[0]; x <= b.bottomright[0]; x++ {
			if b.clay[coord{x, y}] {
				fmt.Printf("X")
			} else if b.ng[coord{x, y}].enc {
				fmt.Printf("~")
			} else if b.ng[coord{x, y}].v {
				fmt.Printf("|")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}

}

func (b Board) fillGraph(n *Node) {
	//fmt.Println(" filling: ", n)
	//b.printGrid()
	if n.v || n.bottom {
		return
	}
	if n.d != nil && !n.d.v {
		b.fillGraph(n.d)
	}
	if n.d == nil || n.d.enc {
		//	fmt.Println("  hit floor")
		if n.l != nil {
			b.fillGraph(n.l)
		}
		n.v = true
		if n.r != nil && !n.r.v {
			b.fillGraph(n.r)
		}

		if n.d == nil || n.d.enc {
			enc := true
			if n.l != nil {
				enc = n.l.enc
			}
			//	fmt.Println("   scanning right", enc)
			for sweep := n.r; sweep != nil && enc; sweep = sweep.r {
				if !(sweep.d == nil || sweep.d.enc) {
					enc = false
				}
				sweep.v = true
			}
			if enc {
				for sweep := n.r; sweep != nil; sweep = sweep.r {
					sweep.enc = true
				}
			}
			n.enc = enc
		}
	}
	n.v = true
}

func (b Board) reachable() (int, int) {
	r := 0
	w := 0
	for xy, x := range b.ng {
		if xy[1] >= b.topleft[1]+1 {
			if x.enc {
				w++
			}
			if x.v {
				r++
			}
		}
	}
	return r, w
}

func (b *Board) getGraph() {
	b.ng = make(NodeGraph)

	for y := 0; y <= b.bottomright[1]; y++ {
		for x := b.topleft[0]; x <= b.bottomright[0]; x++ {
			if b.clay[coord{x, y}] {
				continue
			}
			n := &Node{}
			n.xy = coord{x, y}
			b.ng[coord{x, y}] = n
			if x > b.topleft[0] {
				if x <= b.bottomright[0] {
					if !b.clay[coord{x - 1, y}] {
						n.l = b.ng[coord{x - 1, y}]
						b.ng[coord{x - 1, y}].r = n
					}
				}
			}
			if y > 0 {
				if !b.clay[coord{x, y - 1}] {
					b.ng[coord{x, y - 1}].d = n
				}
				if y == b.bottomright[1] {
					n.bottom = true
				}
			}
		}
	}
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
