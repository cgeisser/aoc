package main

import (
	"bufio"
	"bytes"
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

type node struct {
	value       int
	left, right *node
}

func (n *node) String() string {
	var start *node = nil
	var buf bytes.Buffer
	c := n
	for c != start {
		fmt.Fprintf(&buf, "%d ", c.value)
		c = c.right
		start = n
	}
	return buf.String()

}

func marbleGame(elves, max_marble int) int {
	var elfscores = make([]int, elves)
	var cur_marble = new(node)
	cur_marble.value = 0
	cur_marble.left = cur_marble
	cur_marble.right = cur_marble
	for i := 1; i <= max_marble; i++ {
		if i%23 == 0 {
			for j := 0; j < 7; j++ {
				cur_marble = cur_marble.left
			}
			//fmt.Println("removing: ", cur_marble.value)
			elfscores[i%len(elfscores)] += i + cur_marble.value
			cur_marble.left.right = cur_marble.right
			cur_marble.right.left = cur_marble.left
			cur_marble = cur_marble.right
		} else {
			cur_marble = cur_marble.right

			new_spot := new(node)
			new_spot.value = i
			new_spot.right = cur_marble.right
			new_spot.left = cur_marble
			cur_marble.right.left = new_spot
			cur_marble.right = new_spot

			cur_marble = new_spot
		}

		//fmt.Println("step ", i)
		//fmt.Printf("%v cur: %v\n", cur_marble, cur_marble.value)

	}

	max := elfscores[0]
	for _, v := range elfscores[1:] {
		if v > max {
			max = v
		}
	}

	return max

}
