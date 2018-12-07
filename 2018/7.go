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

const elves = 2

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
	var nextnode string
	for len(graph) > 0 {
		for _, n := range alphanodes {
			if graph[n] != nil &&
				len(graph[n]) == 0 {
				nextnode = n
				fmt.Printf("%v", nextnode)
				break
			}
		}
		for cleanup := range graph {
			delete(graph[cleanup], nextnode)

		}
		delete(graph, nextnode)
	}
	fmt.Printf("\n")

}
