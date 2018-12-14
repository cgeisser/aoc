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
		cycles, _ := strconv.Atoi(input.Text())
		recipeDance(cycles)
	}
}

func recipeDance(cycles int) {
	el := []int{0, 1}
	curelf := 0
	r := make([]int, 2, cycles+10)
	r[0], r[1] = 3, 7
	for i:=0; i< cycles + 10; i++ {
		el, r = onemove(el, curelf, r)
		//fmt.Println(el)
		//fmt.Println(r)	
	}
	fmt.Printf("c: %v out: %v\n", cycles, r[cycles:cycles+10])
}

func onemove(el []int, curelf int, r []int) ([]int, []int) {
	sum := r[el[0]] + r[el[1]]
	if sum / 10 > 0 {
		r = append(r, sum / 10)
	}
	r = append(r, sum % 10)

	for i, v := range el {
		el[i] = (r[v] + v + 1) % len(r)
	}
	return el, r
}
