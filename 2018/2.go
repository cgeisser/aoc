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

	for _, val := range values {
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
	}
	fmt.Println("twos: %v threes: %v", counthist[2], counthist[3])
	fmt.Println(counthist[2] * counthist[3])

}
