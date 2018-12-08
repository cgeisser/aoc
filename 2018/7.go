package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var graph = make(map[string]map[string]bool)

func addEdge(from, to string) {
	edges := graph[from]
	if edges == nil {
		edges = make(map[string]bool)
		graph[from] = edges
	}
	edges[to] = true

	// make a stub for the destination node
	if graph[to] == nil {
		graph[to] = make(map[string]bool)
	}
}

const maxelves = 2
const timebonus = 0

var tokens = make(chan struct{}, maxelves)
var clock = make(chan int)

func main() {
	f, err := os.Open("7.data")
	if err != nil {
		fmt.Printf("can't open file: %v\n", err)
		os.Exit(1)
	}
	input := bufio.NewScanner(f)
	for input.Scan() {
		tokens := strings.Split(input.Text(), " ")
		start := tokens[1]
		end := tokens[7]
		fmt.Printf("%v -> %v\n", start, end)
		addEdge(end, start)
	}
	fmt.Printf("%v", graph)
	f.Close()

	var alphanodes = make([]string, 0, len(graph))
	for name := range graph {
		alphanodes = append(alphanodes, name)
	}
	sort.Strings(alphanodes)
	fmt.Printf("alphanodes: %v\n", alphanodes)

	for len(graph) > 0 {
		var tocleanup = make(chan string)
		var working = false
		for _, n := range alphanodes {
			if graph[n] != nil && len(graph[n]) == 0 {
				go func(nextnode string) {
					tokens <- struct{}{}
					defer func() { <-tokens }()

					timevalue := int(nextnode[0]) - 'A' + 1 + timebonus
					fmt.Printf("%v%v ", nextnode, timevalue)
					tocleanup <- nextnode
				}(n)
				working = true
				break

			}
		}
		// wait and cleanup
		if working {
		for val := range tocleanup {
			for cleanup := range graph {
				delete(graph[cleanup], val)

			}
			delete(graph, val)
		}
			working = false
		}
	}
	fmt.Printf("\n")

}
