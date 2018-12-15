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
	f, _ := os.Open("15.data")
	input := bufio.NewScanner(f)

	grid := make([][]string, 0)
	players := make(map[coord]Player)
	y := 0
	for input.Scan() {
		line := input.Text()
		row := make([]string, len(line))
		for x, _ := range line {
			b := line[x]
			if b == 'G' || b == 'E' {
				row[x] = "."
				players[coord{x, y}] = Player{string(b), 3, 200}
			} else {
				row[x] = string(b)
			}

		}
		y++
		grid = append(grid, row)
	}
	f.Close()

	printGrid(grid, players)
	outcome := []int{1, 1, 0}
	for r := 0; outcome[0] > 0 && outcome[1] > 0; r++ {
		fmt.Println("round: ", r)
		playorder := make(CoordSlice, 0)
		for c := range players {
			playorder = append(playorder, c)
		}

		sort.Sort(playorder)

		for _, c := range playorder {
			cur, ok := players[c]
			if !ok {
				break
			}
			//fmt.Println("ready player: ", c, cur)
			if bad := badGuys(grid, c, players); len(bad) > 0 {
				players = attack(grid, c, players)
			} else {
				dests := make(map[coord]bool)
				for b, badguy := range players {
					d := avail(grid, b, players)
					if cur.t != badguy.t {
						//fmt.Println("available from: ", b, d)
						for _, v := range d {
							dests[v] = true
						}
					}
				}
				move := multiDyk(grid, c, players, dests)
				//fmt.Println(" move: ", move)
				delete(players, c)
				players[move] = cur
				players = attack(grid, move, players)
			}
		}
		outcome = getOutcome(players)
		//printGrid(grid, players)
		fmt.Println(r, outcome)
	}
}

func getOutcome(players map[coord]Player) []int {
	outcome := []int{0, 0, 0}
	for _, p := range players {
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

func attack(grid [][]string, c coord, players map[coord]Player) map[coord]Player {
	if bad := badGuys(grid, c, players); len(bad) > 0 {
		//fmt.Println(" targets in range: ", bad)
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
				//fmt.Println(" selected target: ", x, target)
				target.hp -= players[c].p
				if target.hp <= 0 {
					delete(players, x)
	
				} else {
					players[x] = target
				}
				return players				
			}
		}
	}
	return players
}

func printGrid(grid [][]string, players map[coord]Player) {
	for y, _ := range grid {
		for x, b := range grid[y] {
			if c, ok := players[coord{x, y}]; ok {
				fmt.Printf("%v", c.t)
			} else {
				fmt.Printf("%v", b)
			}
		}
		fmt.Println()
	}
	fmt.Println(players)
}

type CoordSlice []coord

func badGuys(grid [][]string, s coord, players map[coord]Player) map[coord]Player {
	out := make(map[coord]Player)
	for _, n := range neighbors(grid, s, players) {
		if c, ok := players[n]; ok && players[s].t != c.t {
			//fmt.Println("adding badguy ", n, c)
			out[n] = c
		}
	}
	return out
}

func empty(grid [][]string, s coord, players map[coord]Player) bool {
	if _, ok := players[s]; ok {
		return false
	}
	if grid[s[1]][s[0]] == "." {
		return true
	}
	return false
}

func neighbors(grid [][]string, s coord, players map[coord]Player) CoordSlice {
	avail := make([]coord, 0, 4)
	avail = append(avail, coord{s[0], s[1] - 1})
	avail = append(avail, coord{s[0] - 1, s[1]})
	avail = append(avail, coord{s[0] + 1, s[1]})
	avail = append(avail, coord{s[0], s[1] + 1})
	return avail
}

func avail(grid [][]string, s coord, players map[coord]Player) []coord {
	avail := make([]coord, 0, 4)
	for _, n := range neighbors(grid, s, players) {
		if empty(grid, n, players) {
			avail = append(avail, n)
		}
	}
	return avail
}

func multiDyk(grid [][]string, s coord, players map[coord]Player, goal map[coord]bool) coord {
	if _, ok := goal[s]; ok {
		return s
	}
	if len(goal) == 0 {
		return s
	}
	v := make(map[coord]bool)
	movelist := make([][]coord, 0, 10)
	movelist = append(movelist, avail(grid, s, players))

	//fmt.Printf("search: %v for %v\n", s, goal)
	//	fmt.Println("m: ", movelist)
	for len(movelist) > 0 && len(movelist[0]) > 0 {
	//	fmt.Println("m: ", movelist)
		cur := movelist[0]
		lastspot := cur[len(cur) - 1]
		if _, found := goal[lastspot]; found {
			//fmt.Println(cur)
			return cur[0]
		}
		v[lastspot] = true
		nextmoves := avail(grid, lastspot, players)
		for _, a := range nextmoves {
			if _, visited := v[a]; !visited {
				search := []coord{ cur[0], a }
				//fmt.Println("  s: ", search)
				movelist = append(movelist, search)
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
	p  int
	hp int
}
