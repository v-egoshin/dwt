package server

import "github.com/gin-gonic/gin"

var router *gin.Engine

func initializeRoutes() {
	router.GET("/job", showRunningJobs)
	router.POST("/job/create", createNewJob)
	router.GET("/job/:job_id", showJobById)
	router.GET("/job/:job_id/chunk", getJobChunk)

	router.GET("/wordlist", showWordlist)
	router.GET("/wordlist/reindex", reindexWordlist)
}

func showRunningJobs(context *gin.Context) {

}
