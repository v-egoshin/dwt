package main

import (
	"fmt"

	. "github.com/v-egoshin/distributed-word-tool"
)

func main() {
	var wp WordlistPermutations

	fpaths := []string{
		"../distributed-word-tool/test/wl1.txt",
		"../distributed-word-tool/test/wl2.txt",
		"../distributed-word-tool/test/wl3.txt",
		"../distributed-word-tool/test/wl4.txt",
		"../distributed-word-tool/test/1_000_000.txt",
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
