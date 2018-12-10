package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	sum, consumed, val := sumMetaData(license_file)
	fmt.Printf("sum: %v consumed: %v val: %v len: %v\n",
		sum, consumed, val, len(license_file))
}

func sumMetaData(license []int) (int, int, int) {
	var sum int = 0
	var child_ptr int = 2
	if len(license) < 40 {
		fmt.Printf("processing: %v\n", license)
	}
	fmt.Printf("c:%v m:%v len:%v\n", license[0], license[1], len(license))

	// Parse out children recursively
	var childsums = make([]int, license[0])
	for i := 0; i < license[0]; i++ {
		s, c, v := sumMetaData(license[child_ptr : len(license)-license[1]])
		child_ptr += c
		sum += s
		childsums[i] = v
	}
	if child_ptr+license[1] > len(license) {
		panic("consumed too much")
	}

	var val = 0
	// add up our own metadata
	fmt.Println(childsums)
	for i := 0; i < license[1]; i++ {
		sum = sum + license[child_ptr]
		fmt.Println("metadata: ", license[child_ptr])
		if license[child_ptr] > 0 &&
			license[child_ptr]-1 < len(childsums) {
			val += childsums[license[child_ptr]-1]
		}
		child_ptr++
	}
	if len(childsums) == 0 {
		val = sum
	}

	fmt.Printf("sum: %v consumed: %v val: %v\n", sum, child_ptr, val)
	return sum, child_ptr, val
}
