package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("1.data")
	if err != nil {
		fmt.Printf("can't open file: %v\n", err)
		os.Exit(1)
	}
	input := bufio.NewScanner(f)

	var values = make([]int, 0)
	for input.Scan() {

		val, _ := strconv.Atoi(input.Text())
		values = append(values, val)
	}
	var s = 0
	for _, val := range values {
		s = s + val
	}
	fmt.Println("sum ", s)

	var sum int = 0
	var sumdupes = make(map[int]bool)
	var searching = true
	for searching {
		for _, val := range values {
			sum = sum + val
			if sumdupes[sum] {
				fmt.Println("dup sum: ", sum)
				searching = false
				break
			}
			sumdupes[sum] = true
		}
	}
	f.Close()
}
