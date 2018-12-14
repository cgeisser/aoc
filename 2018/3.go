package main

import (
	"bufio"
	"fmt"
	"os"
	//"strconv"
)

func main() {
	f, err := os.Open("3.data")
	if err != nil {
		fmt.Printf("can't open file: %v\n", err)
		os.Exit(1)
	}
	input := bufio.NewScanner(f)

	var values = make([][]int, 0)
	for i:=0; i<1000; i++ {
		values=append(values, make([]int, 1000))
	}
	clean := make(map[int]bool)
	for input.Scan() {
		var e, x, y, xl, yl int
		fmt.Sscanf(input.Text(), "#%d @ %d,%d: %dx%d" , &e, &x, &y, &xl, &yl)
		fmt.Printf("#%d @ %d,%d: %dx%d\n" , e, x, y, xl, yl)
		clean[e] = true
		for i := x; i<x+xl; i++ {
			for j := y; j<y+yl; j++ {
				if values[j][i] == 0 {
					values[j][i] = e
				} else {
					clean[e] = false
					clean[values[j][i]] = false
					values[j][i] = -1
				}
			}
		}
		//for _, r:= range values {
		//	fmt.Println(r)
		//}
	}

	cnt := 0
	for _, r := range values {
		for _, v := range r {
			if v==-1 {
				cnt++
			}

		}
	}
	fmt.Println(cnt)
	for e, c := range clean {
		if c {
			fmt.Println("still clean: ", e)
		}
	}
}
