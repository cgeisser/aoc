package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	//	"strings"
)

type coord [2]int

func main() {
	f, _ := os.Open(os.Args[1])
	input := bufio.NewScanner(f)

	bo := new(Board)
	bo.grid = make([][]string, 0)
	bo.players = make(map[coord]*Player)
	y := 0
	for input.Scan() {
		line := input.Text()
		row := make([]string, len(line))
		for x, _ := range line {
			b := line[x]
			if b == 'G' || b == 'E' {
				row[x] = "."
				bo.players[coord{x, y}] = &Player{string(b), coord{x, y}, 3, 200}
			} else {
				row[x] = string(b)
			}

		}
		y++
		bo.grid = append(bo.grid, row)
	}
	f.Close()

	bo.printGrid()
	outcome := []int{1, 1, 0}
	for r := 1; outcome[0] > 0 && outcome[1] > 0; r++ {
		fmt.Println("round: ", r)
		playorder := make(CoordSlice, 0)
		//snap_players := make(map[coord]*Player)
		for c := range bo.players {
			playorder = append(playorder, c)
		}

		sort.Sort(playorder)

		for _, c := range playorder {
			cur, ok := bo.players[c]
			if !ok {
				continue
			}
			fmt.Println("ready player: ", c, cur)
			outcome = bo.getOutcome()
			if outcome[0] == 0 || outcome[1] == 0 {
				fmt.Println("done before round ended")
				break
			}
			if bad := bo.badGuys(c); len(bad) > 0 {
				bo.attack(c)
			} else {
				dests := make(map[coord]bool)
				for b, badguy := range bo.players {
					d := bo.avail(b)
					if cur.t != badguy.t {
						//fmt.Println("available from: ", b, d)
						for _, v := range d {
							dests[v] = true
						}
					}
				}
				move := bo.multiDyk(c, dests)
				fmt.Println(" move: ", move)
				delete(bo.players, c)
				bo.players[move] = cur
				bo.attack(move)
			}
		}
		outcome = bo.getOutcome()
		bo.printGrid()
		fmt.Println(r, outcome)
	}
}

type Board struct {
	 grid [][]string
	 players map[coord]*Player
}

func (b Board) getOutcome() []int {
	outcome := []int{0, 0, 0}
	for _, p := range b.players {
		if p.t == "E" {
			outcome[0]++
		}
		if p.t == "G" {
			outcome[1]++
		}
		outcome[2] += p.hp
	}
	return outcome
}

func (b *Board) attack(c coord) bool {
	if bad := b.badGuys(c); len(bad) > 0 {
		fmt.Println(" targets in range: ", bad)
		minhits := 0
		attackorder := make([]coord, 0, len(bad))
		for x, b := range bad {
			if minhits == 0 || b.hp < minhits {
				minhits = b.hp
			}
			attackorder = append(attackorder, x)
		}

		sort.Sort(CoordSlice(attackorder))

		for _, x := range attackorder {
			if bad[x].hp == minhits {
				target := bad[x]
				fmt.Println(" selected target: ", x, target)
				target.hp -= b.players[c].p
				if target.hp <= 0 {
					delete(b.players, x)
					fmt.Println(" target destroyed!!")
				}
				return true				
			}
		}
	}
	return false
}

func (b Board) printGrid() {
	for y, _ := range b.grid {
		for x, r := range b.grid[y] {
			if c, ok := b.players[coord{x, y}]; ok {
				fmt.Printf("%v", c.t)
			} else {
				fmt.Printf("%v", r)
			}
		}
		fmt.Println()
	}
	fmt.Println(b.players)
}

type CoordSlice []coord

func (b Board) badGuys(s coord) map[coord]*Player {
	out := make(map[coord]*Player)
	for _, n := range neighbors(s) {
		if c, ok := b.players[n]; ok && b.players[s].t != c.t {
			//fmt.Println("adding badguy ", n, c)
			out[n] = c
		}
	}
	return out
}

func (b Board) empty(s coord) bool {
	if _, ok := b.players[s]; ok {
		return false
	}
	if b.grid[s[1]][s[0]] == "." {
		return true
	}
	return false
}

func neighbors(s coord) CoordSlice {
	avail := make([]coord, 0, 4)
	avail = append(avail, coord{s[0], s[1] - 1})
	avail = append(avail, coord{s[0] - 1, s[1]})
	avail = append(avail, coord{s[0] + 1, s[1]})
	avail = append(avail, coord{s[0], s[1] + 1})
	return avail
}

func (b Board) avail(s coord) []coord {
	avail := make([]coord, 0, 4)
	for _, n := range neighbors(s) {
		if b.empty(n) {
			avail = append(avail, n)
		}
	}
	return avail
}

func (b Board) multiDyk(s coord, goal map[coord]bool) coord {
	if _, ok := goal[s]; ok {
		return s
	}
	if len(goal) == 0 {
		return s
	}
	v := make(map[coord]bool)
	movelist := make([][]coord, 0, 10)
	for _, a := range b.avail(s) {		
		movelist = append(movelist, []coord{ a })
	}

	//fmt.Printf("search: %v for %v\n", s, goal)
	//	fmt.Println("m: ", movelist)
	for len(movelist) > 0 && len(movelist[0]) > 0 {
//		fmt.Println("m: ", movelist)
		cur := movelist[0]
		lastspot := cur[len(cur) - 1]
		if _, found := goal[lastspot]; found {
			//fmt.Println(cur)
			return cur[0]
		}
		v[lastspot] = true
		nextmoves := b.avail(lastspot)
		for _, a := range nextmoves {
			if _, visited := v[a]; !visited {
				search := []coord{ cur[0], a }
				//fmt.Println("  s: ", search)
				movelist = append(movelist, search)
				v[a] = true
			}
		}
		movelist = movelist[1:]
	}
	return s
}

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

type Player struct {
	t  string
	xy coord
	p  int
	hp int
}
