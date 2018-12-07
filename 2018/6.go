package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct { X, Y int }

func Boundaries(points []Point) (Point, Point) {
        min, max := points[0], points[0]
	for _, p := range points[1:] {
		if min.X > p.X {
			min.X = p.X
		}
		if min.Y > p.Y {
			min.Y = p.Y}
		if max.X < p.X {
			max.X = p.X
		}
		if max.Y < max.Y {
			max.Y = p.Y
		}
	}
	return min, max
}
		

func main() {
	var points []Point 
	f, err:= os.Open("6.data")
        if err != nil {
		fmt.Printf("can't open file: %v\n", err)
		os.Exit(1)
	}
	input := bufio.NewScanner(f)
	for input.Scan() {
		var p Point
		fmt.Sscanf(input.Text(), "%d, %d", &p.X, &p.Y)
		points = append(points, p)
	}
	f.Close()

	fmt.Printf("%v\n", points)
	min, max := Boundaries(points)
	fmt.Printf("%v %v\n", min, max)
}
