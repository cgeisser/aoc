package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, _ := os.Open("9.data")
	input := bufio.NewScanner(f)
	for input.Scan() {
		tokens := strings.Split(input.Text(), " ")
		elves, _ := strconv.Atoi(tokens[0])
		max_marble, _ := strconv.Atoi(tokens[6])
		fmt.Printf("%v %v \n", elves, max_marble)
		fmt.Println("final score: ", marbleGame(elves, max_marble))
		//break
	}
}

func marbleGame(elves, max_marble int) int {
	var elfscores = make([]int, elves)

	var marbles = make([]int, 1, max_marble)
	var cur_marble = 0
	for i := 1; i <= max_marble; i++ {
		if i%23 == 0 {
			rem := (len(marbles) + cur_marble - 7) % len(marbles)
			elfscores[i%len(elfscores)] += i + marbles[rem]
			copy(marbles[rem:], marbles[rem+1:])
			marbles = marbles[:len(marbles)-1]
			cur_marble = rem % len(marbles)
		} else {
			new_spot := (cur_marble + 2) % len(marbles)
			marbles = append(marbles, -1)
			copy(marbles[new_spot+1:], marbles[new_spot:])

			marbles[new_spot] = i
			cur_marble = new_spot

		}
		if i%100000 == 0 {
			fmt.Println("step ", i)
			//fmt.Printf("step: %v %v cur: %v\n", i, marbles, marbles[cur_marble])
		}
	}

	max := elfscores[0]
	for _, v := range elfscores[1:] {
		if v > max {
			max = v
		}
	}

	return max

}
