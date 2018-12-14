package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {
	f, err := os.Open("5.data")
	if err != nil {
		fmt.Printf("can't open file: %v\n", err)
		os.Exit(1)
	}
	input := bufio.NewScanner(f)

	for input.Scan() {
		elements := list.New()
		for _, c := range input.Text() {
			elements.PushBack(string(c))
		}
		//fmt.Println("final length: ", reactAll(elements))
		var wg sync.WaitGroup
		l := make(chan int, 26)
		for b := byte('a'); b < byte('z')+1; b++ {
			wg.Add(1)
			b := b
			go func() {
				filtered := list.New()
				for e := elements.Front(); e != nil; e = e.Next() {
					c := e.Value.(string)[0]
					//fmt.Println(string(c), string(b), string(b-'a'+'A'))
					if c != b && c != b-'a'+'A' {
						filtered.PushBack(string(c))
					}
				}
				//for e := filtered.Front(); e != nil; e = e.Next() {
				//	fmt.Printf("%v", e.Value)
				//}
				//fmt.Println()
				if filtered.Len() != elements.Len() {
					l <- reactAll(filtered)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(l)
		min_length := elements.Len()
		for v := range l {
			if v < min_length {
				min_length = v
			}
		}
		fmt.Println("min length", min_length)
		//break
	}

}

func reactAll(elements *list.List) int {
	prevlen := 0
	for prevlen != elements.Len() {
		prevlen = elements.Len()
		for e := elements.Front(); e != nil && e.Next() != nil; e = e.Next() {
			cur := e.Value.(string)
			next := e.Next().Value.(string)
			if cur != next && strings.ToLower(cur) == strings.ToLower(next) {
				var save *list.Element
				if e.Prev() == nil {
					save = elements.Front()
				} else {
					save = e.Prev()
				}

				if e.Next() != nil {
					elements.Remove(e.Next())
				}
				elements.Remove(e)
				e = save
			}
		}

		//for e := elements.Front(); e != nil; e = e.Next() {
		//      fmt.Printf("%v", e.Value)
		//}
		//fmt.Println()
	}
	return elements.Len()
}
