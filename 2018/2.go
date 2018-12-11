package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("2.data")
	if err != nil {
		fmt.Printf("can't open file: %v\n", err)
		os.Exit(1)
	}
	input := bufio.NewScanner(f)

	var values = make([]string, 0)
	for input.Scan() {

		val := input.Text()
		values = append(values, val)
	}
	f.Close()

	var counthist = make(map[int]int)
	var match, match2 string
	for loc, val := range values {
		fmt.Println("input: ", val)
		var stringhist = make(map[rune]int)
		for _, c := range val {
			stringhist[c]++
		}
		fmt.Println("stringhist: ", stringhist)
		var invert = make(map[int]bool)

		for c := range stringhist {
			invert[stringhist[c]] = true
		}
		for c := range invert {
			counthist[c]++
		}
		// do search for part 2
		if match == "" && match2 == "" {
			for _, m := range values[loc:] {
				var diffs = 0
				for i := range m {
					if val[i] != m[i] {
						diffs++
					}
					if diffs > 1 {
						break
					}
				}
				if diffs == 1 {
					match = m
					match2 = val
				}
			}
		}
	}
	fmt.Println("twos: %v threes: %v", counthist[2], counthist[3])
	fmt.Println(counthist[2] * counthist[3])
	fmt.Println(match)
	fmt.Println(match2)
}
