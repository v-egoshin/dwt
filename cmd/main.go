package main

import (
	"fmt"

	. "github.com/v-egoshin/dwt"
)

func main() {
	var wp WordlistPermutations

	fpaths := []string{
		"../dwt/test/wl1.txt",
		"../dwt/test/wl2.txt",
		"../dwt/test/wl3.txt",
		"../dwt/test/wl4.txt",
		"../dwt/test/1_000_000.txt",
	}

	wp.Initialize(fpaths)
	fmt.Println(wp.Count)

	p := make(chan []int, 0)
	go wp.Permute(p, 129000000, 129000020)
	for {
		_, ok := <-p
		if !ok {
			break
		}
		//	fmt.Println(wp.GetPermuteByState(pair))
	}

}
