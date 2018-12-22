package main

import (
	"bufio"
	"fmt"
	"os"
)

type room struct {
	xy   coord
	dist int
	ar   *allrooms
	v    bool
}

type dir int

const (
	N = iota
	E
	S
	W
)

func index(d byte) dir {
	switch d {
	case 'N':
		return N
	case 'E':
		return E
	case 'S':
		return S
	case 'W':
		return W
	default:
		panic(fmt.Sprintf("no direction %c", d))
	}
}

type coord [2]int
type allrooms map[coord]*room

func (d dir) moveCoord(xy coord) coord {
	switch d {
	case N:
		return coord{xy[0], xy[1] - 1}
	case E:
		return coord{xy[0] + 1, xy[1]}
	case S:
		return coord{xy[0], xy[1] + 1}
	case W:
		return coord{xy[0] - 1, xy[1]}
	default:
		panic("bad direction")
	}
}

type parsetree struct {
	s string
	children []*parsetree
}

func build(re string) (*parsetree, int) {
	clist := make([]*parsetree)
	for i, c := range re {
		if c == '|' || c=='(' {
			subtree, eaten := parsetree(re[i+1:])
			clist = append(clist, subtree)
			break
		}
	}
	pt := &parsetree{ re[:i], clist, 0)}
	
}

func (cr *room) build(re string, closeblock) int {
	fmt.Println("  parsing: ", re, closeblock)
	if re == "" {
		return 0
	}
	if re[0] == '|' {
		cr.build(re[closeblock:], -1)
		return 0
	}
	if re[0] == ')' {
		return cr.build(re[1:], -1) + 1
	}

	if re[0] == '(' {
		consumed := 1
		pcount := 0
		closeparen := 0
		for i, c := range re {
			if c == '(' {
				pcount++
			}
			if c == ')' {
				pcount--
			}
			if c == ')' && pcount == 0 {
				closeparen = i
				break
			}
		}
		options := 0
		for ; consumed < closeparen; consumed++ {
			fmt.Println(" option", options, re[consumed:])
			consumed += cr.build(re[consumed:], closeparen-consumed)
			fmt.Println(" end option", options)
options++
		}
		return consumed
	} else {
		d := index(re[0])
		movec := d.moveCoord(cr.xy)

		nextroom := (*cr.ar)[movec]
		if nextroom == nil {
			nextroom = &room{movec, cr.dist + 1, cr.ar, false}
			//fmt.Println("A new room!", nextroom)
			(*cr.ar)[movec] = nextroom
		} else if cr.dist+1 < nextroom.dist {
			fmt.Println("shorter path: ", cr.dist+1)
			nextroom.dist = cr.dist + 1
		}
		return nextroom.build(re[1:], closeblock-1)
	}
}

func (r *room) maxdist() int {
	max := 0
	for _, n := range *r.ar {
		if n.dist > max {
			max = n.dist
		}
	}
	return max
}

func main() {
	f, _ := os.Open("20.data")
	input := bufio.NewScanner(f)
	for input.Scan() {
		if input.Text()[0] == '#' {
			continue
		}
		if input.Text()[0] != '^' || input.Text()[len(input.Text())-1] != '$' {
			panic("bad row")
		}

		r := &room{}
		r.ar = &allrooms{ r.xy:r }
		fmt.Println("running ",input.Text())
		r.build(input.Text()[1 : len(input.Text())-1], -1)
		for x, v := range *r.ar {
			fmt.Println(x, v)
		}
		fmt.Println(" max dist:", r.maxdist())
		//	break
	}

}
