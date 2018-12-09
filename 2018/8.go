package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("can't open file: %v\n", err)
		os.Exit(1)
	}
	input := bufio.NewScanner(f)
	input.Scan()
	var contents = input.Text()
	f.Close()
	
	tokens := strings.Split(contents, " ")
	var license_file = make([]int, 0, len(tokens))
	for _, v := range tokens {
		i, _ := strconv.Atoi(v)
		license_file = append(license_file, i)
	}
	fmt.Printf("input data: %v\n", license_file)	
        sum, consumed := sumMetaData(license_file)
	fmt.Printf("sum: %v consumed: %v len: %v\n",
		sum, consumed, len(license_file))
}

func sumMetaData(license []int) (int, int) {
	var sum int = 0
	var child_ptr int = 2
	if (len(license) < 4000) {
	  fmt.Printf("processing: %v\n", license)
	}
	if (len(license) < 3*license[0] + license[1] + 2) {
		panic("broken array length")
	}
	if (license[1] < 1) {
		panic("no metadata")
	}
	fmt.Printf("c:%v m:%v len:%v\n", license[0], license[1], len(license))

	// Parse out children recursively
	for i:=0 ; i < license[0]; i++ {
		s, c := sumMetaData(license[child_ptr:len(license) - license[1]])		
		child_ptr = child_ptr + c
		sum = sum + s
	}
	if child_ptr + license[1] > len(license) {
		panic("consumed too much")
	}

	// add up our own metadata
	for i:= 0; i < license[1]; i++ {
		sum = sum + license[child_ptr]
		child_ptr++
	}
	fmt.Printf("sum: %v consumed: %v\n", sum, child_ptr)
	return sum, child_ptr			
}
