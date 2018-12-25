package main

import (
	"bufio"
	"fmt"
	"sort"

	"os"
	"regexp"
	"strconv"
	"strings"
)

type team int

type group struct {
	id         int
	side       string
	units      int
	attack     int
	weapon     string
	hitpoints  int
	initiative int
	immune     map[string]bool
	weak       map[string]bool
	curtarget  *group
	targetted  bool
}

func (a group) String() string {
	return fmt.Sprintf("%v group %v %v units", a.side, a.id, a.units)
}

func (a group) getDamage(t group) int {
	if a.side == t.side {
		//don't hurt our friends
		return 0
	}
	if t.immune[a.weapon] {
		return 0
	}
	pwr := a.pwr()
	if t.weak[a.weapon] {
		return 2 * pwr
	}
	//fmt.Printf("     %v %v would damage %v by %v\n", a.side, a.id, t.id, pwr)
	return pwr
}

func (a group) pwr() int {
	return a.units * a.attack
}

func (a group) hit() {
	if a.units == 0 {
		fmt.Println("can't hit, I'm dead: ", a)
		return
	}
	if a.curtarget != nil {
		damage := a.getDamage(*a.curtarget)

		if damage > 0 {
			unitslost := damage / a.curtarget.hitpoints
			fmt.Printf("   %v group %v hits %v damage %d killed %d\n", a.side, a.id, a.curtarget.id, damage, unitslost)
			if a.curtarget.units > unitslost {
				a.curtarget.units -= unitslost
			} else {
				a.curtarget.units = 0
			}
		}
	}
	a.curtarget = nil
}

type players []*group
type byInitiative []*group

func (x byInitiative) Len() int      { return len(x) }
func (x byInitiative) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x byInitiative) Less(i, j int) bool {
	if x[i].initiative > x[j].initiative {
		return true
	}

	return false
}

type byPower []*group

func (x byPower) Len() int      { return len(x) }
func (x byPower) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x byPower) Less(i, j int) bool {
	if x[i].pwr() == x[j].pwr() {
		return byInitiative(x).Less(i, j)
	}
	return x[i].pwr() > x[j].pwr()
}

type byHits struct {
	a *group
	t players
}

func (x byHits) Len() int      { return len(x.t) }
func (x byHits) Swap(i, j int) { x.t[i], x.t[j] = x.t[j], x.t[i] }

func (x byHits) Less(i, j int) bool {
	if x.a.getDamage(*x.t[i]) == x.a.getDamage(*x.t[j]) {
		return byPower(x.t).Less(i, j)
	}
	return x.a.getDamage(*x.t[i]) > x.a.getDamage(*x.t[j])
}

func main() {
	f, _ := os.Open(os.Args[1])
	p := make(players, 0)
	re := regexp.MustCompile("(\\d+).*?(\\d+).*?(\\d+) (\\w+).*?(\\d+)")
	weakness := regexp.MustCompile("weak to ([\\w, ]+)")
	immunity := regexp.MustCompile("immune to ([\\w, ]+)")
	input := bufio.NewScanner(f)
	input.Scan()
	side := input.Text()[:len(input.Text())-1]
	id := 1
	for input.Scan() {
		if len(input.Text()) == 0 {
			input.Scan()
			side = input.Text()[:len(input.Text())-1]
			input.Scan()
			id = 1
		}
		g := group{}
		g.id = id
		id++
		p = append(p, &g)
		g.immune = make(map[string]bool)
		g.weak = make(map[string]bool)
		g.side = side
		//var submatch string
		matches := re.FindStringSubmatch(input.Text())
		g.units, _ = strconv.Atoi(matches[1])
		g.hitpoints, _ = strconv.Atoi(matches[2])
		g.attack, _ = strconv.Atoi(matches[3])
		g.weapon = matches[4]
		g.initiative, _ = strconv.Atoi(matches[5])

		// weaknesses
		matches = weakness.FindStringSubmatch(input.Text())
		if len(matches) == 2 {
			for _, w := range strings.Split(matches[1], ", ") {
				g.weak[w] = true
			}
		}

		// imunity
		matches = immunity.FindStringSubmatch(input.Text())
		if len(matches) == 2 {
			for _, i := range strings.Split(matches[1], ", ") {
				g.immune[i] = true
			}
		}
	}

	p.tothedeath()
}

func (p players) count() (map[string]int, bool) {
	ret := make(map[string]int)
	for _, g := range p {
		ret[g.side] += g.units
	}
	var done bool
	for _, s := range ret {
		if s == 0 {
			done = true
			break
		}
	}
	return ret, done
}

func (p players) tothedeath() {
	result, done := p.count()
	for !done {

		// copy active players
		active := make(players, 0)
		for _, g := range p {
			g.targetted = false
			g.curtarget = nil
			if g.units > 0 {
				active = append(active, g)
				fmt.Println(g)
			}
		}

		sort.Sort(byPower(active))

		targetlist := make(players, len(active))
		copy(targetlist, active)
		// for each active player, pick a target
		for _, a := range active {
			sort.Sort(byHits{a, targetlist})
			fmt.Println(" choosing target for", a, targetlist)
			for _, t := range targetlist {
				if a.side != t.side && !t.targetted && a.getDamage(*t) > 0 {
					fmt.Println("     picked:", t)
					a.curtarget = t
					t.targetted = true
					break
				}
			}
		}
		// fight
		sort.Sort(byInitiative(active))
		for _, a := range active {
			fmt.Println(a, a.initiative)
			a.hit()
		}

		result, done = p.count()
		fmt.Println(" end of round:", result)
	}
}
