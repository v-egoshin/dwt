package server

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/v-egoshin/dwt"
)

var flagWordlistPath string

func Run() {

	flag.StringVar(&flagWordlistPath, "path", "./test", "Path to wordlists.")
	flag.Parse()
	//flagWordlistPath := flag.String("path", "./wordlists", "Path to wordlists.")
	Wordlists = dwt.ListWordlists(flagWordlistPath)

	//
	//var wlp dwt.WordlistPermutations
	//wlp.Initialize(server.Wordlists[0:3])
	//fmt.Println(wlp.EndState())

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	rg := r.RouterGroup
	r.RouterGroup = *InitializeRoutes(&rg)
	r.Run(":8080")
}
