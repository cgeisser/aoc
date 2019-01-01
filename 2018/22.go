package main

import (
	"container/heap"
	"fmt"
)

const depth = 9465
const tx = 13
const ty = 704

type rtype int

const (
	ROCK rtype = iota
	WET
	NARROW
)

type region struct {
	rt rtype
	el int
}

type coord struct {
	x, y int
}

type grid map[coord]region

func (g *grid) getType(c coord) (rtype, int) {
	if v, ok := (*g)[c]; ok {
		return v.rt, v.el
	}
	gi := 0
	if c.y == 0 {
		gi = c.x * 16807
	} else if c.x == 0 {
		gi = c.y * 48271
	} else if c.x == tx && c.y == ty {
		gi = 0
	} else {
		_, lval := g.getType(coord{c.x - 1, c.y})
		_, uval := g.getType(coord{c.x, c.y - 1})
		gi = lval * uval
	}
	el := (gi + depth) % 20183
	rt := rtype(el % 3)
	(*g)[c] = region{rt, el}

	return rt, el
}

func main() {
	board := make(grid)

	risk := 0
	for x := 0; x <= tx; x++ {
		for y := 0; y <= ty; y++ {
			rt, _ := board.getType(coord{x, y})
			risk += int(rt)
		}
	}
	fmt.Println("risk: ", risk)
	fmt.Println("shortest: ", board.shortest(coord{0, 0}, coord{tx, ty}))
}

type tool int

const (
	TORCH tool = iota
	NONE
	CLIMB
)

type node struct {
	c coord
	t tool
}

type nodecost struct {
	node
	cost, est int
}

type edgelist []nodecost

func (el edgelist) Less(i, j int) bool { return el[i].cost+el[i].est < el[j].cost+el[j].est }
func (el edgelist) Len() int           { return len(el) }
func (el edgelist) Swap(i, j int)      { el[i], el[j] = el[j], el[i] }

func (el *edgelist) Push(i interface{}) {
	*el = append(*el, i.(nodecost))
}

func (el *edgelist) Pop() interface{} {
	ret := (*el)[len(*el)-1]
	*el = (*el)[:len(*el)-1]
	return ret
}

func dist(s, t coord) int {
	dx := s.x - t.x
	dy := s.y - t.y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

func toolcombos(a, b rtype) map[tool]bool {
	allowed := make(map[tool]bool)
	switch a {
	case ROCK:

		allowed[CLIMB] = true
		allowed[TORCH] = true

	case WET:

		allowed[CLIMB] = true
		allowed[NONE] = true

	case NARROW:

		allowed[TORCH] = true
		allowed[NONE] = true

	}
	switch b {
	case ROCK:
		allowed[NONE] = false
	case WET:
		allowed[TORCH] = false
	case NARROW:
		allowed[CLIMB] = false
	}
	return allowed
}

func (g grid) neighbors(c coord) []node {
	dests := make([]coord, 0)

	if c.x > 0 {
		dests = append(dests, coord{c.x - 1, c.y})
	}
	if c.y > 0 {
		dests = append(dests, coord{c.x, c.y - 1})
	}
	dests = append(dests, coord{c.x + 1, c.y})
	dests = append(dests, coord{c.x, c.y + 1})

	options := make([]node, 0)
	ctype, _ := g.getType(c)
	for _, d := range dests {
		dtype, _ := g.getType(d)
		tools := toolcombos(ctype, dtype)
		for t, v := range tools {
			if v {
				options = append(options, node{d, t})
			}
		}
	}
	return options
}

func (g grid) shortest(start, target coord) int {
	explore := make(edgelist, 0)
	heap.Push(&explore, nodecost{node{start, TORCH}, 0, dist(start, target)})

	beenthere := make(map[node]bool)
	for len(explore) > 0 {
		current := heap.Pop(&explore).(nodecost)
		beenthere[current.node] = true
		if current.c == target {
			return current.cost
		}
		next := g.neighbors(current.c)
		for _, n := range next {
			if !beenthere[n] {
				cost := 1
				if n.t != current.t {
					cost = 8
				}
				heap.Push(&explore, nodecost{n, current.cost + cost, dist(n.c, target)})
			}
		}

	}
	return 0
}
