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
	dim   int
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

const gridmax = 300

func sumGrid(grid [gridmax][gridmax]int, i point, out chan<- point) {
	for ; i.dim <= gridmax && i.x+i.dim-1 <= gridmax && i.y+i.dim-1 <= gridmax; i.dim++ {
		//fmt.Println("thread processing: ", i)

		for gx := i.x; gx < i.x+i.dim; gx++ {
			//	fmt.Println(gx, i.y+i.dim - 1)
			i.power += grid[gx-1][i.y+i.dim-2]
		}
		for gy := i.y; gy < i.y+i.dim-1; gy++ {
			//	fmt.Println(i.x+i.dim - 1, gy)
			i.power += grid[i.x+i.dim-2][gy-1]
		}

		//fmt.Println("thread done processing: ", i)
		out <- i
	}

}

func gridPower(magic int, gx, gy int) int {
	rack := gx + 10
	pwr := rack * gy
	pwr += magic
	pwr *= rack
	h := (pwr % 1000) / 100
	return h - 5
}

func main() {
	fmt.Println(gridPower(8, 3, 5))
	fmt.Println(gridPower(57, 122, 79))
	fmt.Println(gridPower(39, 217, 196))
	fmt.Println(gridPower(71, 101, 153))

	f, _ := os.Open("11.data")
	input := bufio.NewScanner(f)

	for input.Scan() {
		magic, _ := strconv.Atoi(input.Text())

		gather := make(chan point, gridmax*gridmax)
		findMax := make(chan point)
		var wg sync.WaitGroup
		go getMax(gather, findMax)
		var grid [gridmax][gridmax]int
		// build grid
		for x := range grid {
			for y := range grid[x] {
				grid[x][y] = gridPower(magic, x+1, y+1)
			}
		}

		// search
		maxthreads := make(chan struct{}, gridmax)
		for x := 1; x <= gridmax; x++ {
			for y := 1; y <= gridmax; y++ {
				wg.Add(1)
				p := point{x, y, 1, 0}
				maxthreads <- struct{}{}
				go func() {
					defer wg.Done()
					sumGrid(grid, p, gather)
					<-maxthreads
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
	}
}
