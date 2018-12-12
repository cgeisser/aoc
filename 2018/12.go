package main

import (
	"bufio"
	//	"bytes"
	"fmt"
	"os"
	//	"strconv"
	"strings"
)

func main() {
	f, _ := os.Open("12.data")
	input := bufio.NewScanner(f)
	var start string
	trans := make(map[string]byte)
	for input.Scan() {
		line := input.Text()
		tokens := strings.Split(line, " ")
		if len(tokens) > 1 {
			if tokens[0] == "initial" {
				start = tokens[2]
			} else {
				trans[tokens[0]] = tokens[2][0]
			}
		}
	}
	zero_loc := 0

	fmt.Println("start state: ", start)
	fmt.Println("transforms: ", trans)

	orig_start := start
	for i := 0; i < 2000; i++ {
		start, zero_loc = generate(start, zero_loc, trans)
		//fmt.Println(start,zero_loc)

		sum := 0
		for i, v := range start {
			if v == '#' {
				sum += i - zero_loc
			}
		}
		fmt.Printf("sum of remaining plants @ %v: %v\n", i, sum)
	}

	// Part 2 ??
	zero_loc = 0
	start = orig_start
	for i := 0; ; i++ {
		start, zero_loc = generate(start, zero_loc, trans)
		if i%1000 == 0 {
			fmt.Printf("%v ", i)
		}
		if strings.Contains(start, orig_start) {
			fmt.Println("cycle at %v %v %v\n", i, start, orig_start)
		}
	}
}

func generate(start string, zero int, trans map[string]byte) (string, int) {

	if start[0:4] != "...." {
		start = "...." + start
		zero += 4
	}
	if start[len(start)-5:len(start)-1] != "...." {
		start = start + "...."
	}

	var next []byte = make([]byte, 0, len(start))
	next = append([]byte(start[0:2]))

	for i := 2; i < len(start)-2; i++ {
		if trans[start[i-2:i+3]] != 0 {
			next = append(next, trans[start[i-2:i+3]])
		} else {
			next = append(next, '.')
		}

	}
	return string(next), zero
}
