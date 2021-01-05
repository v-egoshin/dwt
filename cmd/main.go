package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	. "github.com/v-egoshin/dwt"
	"github.com/v-egoshin/dwt/server"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	rg := r.RouterGroup
	r.RouterGroup = *server.InitializeRoutes(&rg)
	r.Run(":8080")
	os.Exit(0)
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

	receiver := make(chan []int, 0)
	go wp.Permute(receiver, 129000000, 129000020)
	for {
		_, ok := <-receiver
		if !ok {
			break
		}
		//	fmt.Println(wp.GetPermuteByState(pair))
	}

}
