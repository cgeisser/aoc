package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	f, _ := os.Open("14.data")
	input := bufio.NewScanner(f)
	for input.Scan() {
		digits := make([]int, 0)
		cycles, _ := strconv.Atoi(input.Text())
		for _, b := range input.Text() {
			d, _ := strconv.Atoi(string(b))
			digits = append(digits, d)
		}
		recipeDance(cycles, digits)
		//break
	}
}

func equal(x, y []int) bool {
	//	fmt.Println("cmp ", x, y)
	if len(x) != len(y) {
		return false
	}
	for i, v := range x {
		if v != y[i] {
			return false
		}
	}
	return true
}

func recipeDance(cycles int, digits []int) {
	el := []int{0, 1}
	r := make([]int, 2, cycles+10)
	r[0], r[1] = 3, 7
	found := false
	tail := 0
	for i := 0; !found || i < cycles+1; i++ {
		el, r = onemove(el, r)
		//fmt.Println(el)
		//fmt.Println(r)
		// search for my sequence
		for !found && tail+len(digits) < len(r) {
			//	fmt.Println(r)
			//	fmt.Println("checking tail", tail, r[tail:tail+len(digits)])
			found = equal(r[tail:tail+len(digits)], digits)
			if found {
				fmt.Println("found at: ", tail)
			}
			tail++
		}
	}
	fmt.Printf("c: %v out: %v\n", cycles, r[cycles:cycles+10])
}

func onemove(el []int, r []int) ([]int, []int) {
	sum := r[el[0]] + r[el[1]]
	if sum/10 > 0 {
		r = append(r, sum/10)
	}
	r = append(r, sum%10)

	for i, v := range el {
		el[i] = (r[v] + v + 1) % len(r)
	}

	return el, r
}
