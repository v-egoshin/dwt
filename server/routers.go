package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//https://stackoverflow.com/questions/62906766/how-to-group-routes-in-gin

func InitializeRoutes(g *gin.RouterGroup) *gin.RouterGroup {
	v1 := g.Group("/")
	{
		jobs := v1.Group("/job")
		{
			jobs.GET("/", showRunningJobs)
			jobs.GET("/:job_id", showJobById)
			jobs.POST("/create", createNewJob)
		}
	}
	return v1
	//Router.
	//Router.POST("/job/create", createNewJob)
	//Router.GET("/job/:job_id", showJobById)
	//Router.GET("/job/:job_id/chunk", getJobChunk)
	//
	//Router.GET("/wordlist", showWordlist)
	//Router.GET("/wordlist/reindex", reindexWordlist)
}

func reindexWordlist(context *gin.Context) {
	//TODO: redindex
}

func showWordlist(context *gin.Context) {
	//TODO: get from database after reindex
}

func getJobChunk(context *gin.Context) {

}

func showJobById(context *gin.Context) {
	name := context.Param("job_id")
	context.String(http.StatusOK, "Hello %s", name)
}

func createNewJob(context *gin.Context) {

}

func showRunningJobs(context *gin.Context) {
	context.String(http.StatusOK, "ohai")
}
