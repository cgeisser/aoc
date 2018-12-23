package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type coord [3]int

type bot struct {
	xyz coord
	pwr int
	ir  int
}

func (b bot) dist(t bot) int {
	x := math.Abs(float64(b.xyz[0] - t.xyz[0]))
	y := math.Abs(float64(b.xyz[1] - t.xyz[1]))
	z := math.Abs(float64(b.xyz[2] - t.xyz[2]))
	return int(x + y + z)
}

func (b bot) inrange(t bot) bool {
	if b.dist(t) <= b.pwr {
		return true
	}
	return false

}

func main() {
	f, _ := os.Open(os.Args[1])
	input := bufio.NewScanner(f)

	allbots := make([]*bot, 0)
	for input.Scan() {
		b := &bot{}
		_, ok := fmt.Sscanf(input.Text(), "pos=<%d,%d,%d>, r=%d", &b.xyz[0], &b.xyz[1], &b.xyz[2], &b.pwr)
		if ok != nil {
			panic(fmt.Sprintf("bad line %v, %s", ok, input.Text()))
		}
		allbots = append(allbots, b)
	}

	maxbot := allbots[0]

	for _, s := range allbots {
		if s.pwr > maxbot.pwr {
			maxbot = s
		}
		for _, t := range allbots {
			if s.inrange(*t) {
				s.ir++
			}
		}
	}

	fmt.Println("max", maxbot)

}
