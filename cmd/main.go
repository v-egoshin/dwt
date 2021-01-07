package main

import (
	"os"

	"github.com/v-egoshin/dwt/server"
)

func main() {

	server.Run()
	os.Exit(0)
	//var wlp WordlistPermutations
	//
	//fpaths := []string{
	//	"../dwt/test/wl1.txt",
	//	"../dwt/test/wl2.txt",
	//	"../dwt/test/wl3.txt",
	//	"../dwt/test/wl4.txt",
	//	"../dwt/test/1_000_000.txt",
	//}
	//
	//wlp.Initialize(fpaths)
	//fmt.Println(wlp.Count)
	//
	//receiver := make(chan []int, 0)
	//go wlp.Permute(receiver, 129000000, 129000020)
	//for {
	//	_, ok := <-receiver
	//	if !ok {
	//		break
	//	}
	//	//	fmt.Println(wlp.GetPermuteByState(pair))
	//}

}
