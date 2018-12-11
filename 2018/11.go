package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type point struct {
	x, y  int
	power int
}

func getMax(in <-chan point, out chan<- point) {
	max := <-in
	for p := range in {
		if p.power > max.power {
			max = p
		}
	}
	out <- max
	close(out)
}

func sumGrid(magic int, i point, out chan<- point) {
	//fmt.Println("thread processing: ", i)
	for gx := i.x; gx < i.x + 3; gx++ {
		for gy := i.y; gy < i.y + 3; gy++ {			
			i.power += gridPower(magic, gx, gy)
		}
	}
	//fmt.Println("thread done processing: ", i)
	out <- i
}

func gridPower(magic int, gx, gy int) int {
	rack := gx + 10
	pwr := rack * gy
	pwr += magic
	pwr *= rack
	h := (pwr % 1000) / 100
	return h - 5
}

const gridmax = 300

func main() {
	fmt.Println(gridPower(8, 3, 5))
	fmt.Println(gridPower(57, 122, 79))
	fmt.Println(gridPower(39, 217, 196))
	fmt.Println(gridPower(71, 101, 153))
	
	f, _ := os.Open("11.data")
	input := bufio.NewScanner(f)

	for input.Scan() {
		magic, _ := strconv.Atoi(input.Text())

		gather := make(chan point)
		findMax := make(chan point)
		var wg sync.WaitGroup
		go getMax(gather, findMax)
		for x := 1; x <= gridmax-2; x++ {
			for y := 1; y <= gridmax-2; y++ {
				wg.Add(1)
				p := point{x, y, 0}
				go func() {
					defer wg.Done()
					sumGrid(magic, p, gather)
				}()
			}
		}
		// closer
		go func() {
			wg.Wait()
			close(gather)
		}()

		answer := <-findMax
		fmt.Printf("magic: %v answer: %v\n", magic, answer)
		//break
	}
}
