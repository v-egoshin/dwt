package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/v-egoshin/dwt"
)

//https://stackoverflow.com/questions/62906766/how-to-group-routes-in-gin

func InitializeRoutes(g *gin.RouterGroup) *gin.RouterGroup {
	v1 := g.Group("/")
	{
		jobs := v1.Group("/job")
		{
			jobs.GET("/", showRunningJobs)
			jobs.GET("/:job_id", showJobById)
			jobs.GET("/:job_id/get/:number", getJobRange)
			jobs.POST("/create", createNewJob)
		}
		wordlists := v1.Group("/wordlist")
		{
			wordlists.GET("/", showWordlist)
			wordlists.GET("/reindex", reindexWordlist)
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

func getJobRange(context *gin.Context) {
	//TODO: permutes
}

func reindexWordlist(context *gin.Context) {
	Wordlists = dwt.ListWordlists(flagWordlistPath)
	context.JSON(200, gin.H{
		"status": "OK",
		"code":   200,
	})
}

type ResponseWordlist struct {
	Index int    `json:"index"`
	Path  string `json:"path"`
	Count int    `json:"count"`
}

func showWordlist(context *gin.Context) {
	var wl []ResponseWordlist
	for i, w := range Wordlists {
		wl = append(wl, ResponseWordlist{
			Index: i,
			Path:  w.Path,
			Count: int(w.Lines),
		})
	}
	context.JSON(200, gin.H{
		"status": "OK",
		"code":   200,
		"data":   wl,
	})
}

func getJobChunk(context *gin.Context) {

}

func showJobById(context *gin.Context) {
	id := context.Param("job_id")
	uid, err := uuid.FromString(id)
	if err != nil {
		context.JSON(404, gin.H{
			"status": "Error",
			"error":  fmt.Sprintf("Job %s not found", id),
		})
	}

	for _, j := range Jobs {
		if j.ID == uid {
			context.JSON(200, gin.H{
				"status": "Ok",
				"data":   j,
			})
		}
	}

}

type NewJsonJob struct {
	WordlistIds []int `json:"ids"`
}

func createNewJob(context *gin.Context) {
	job := new(NewJsonJob)
	if err := context.BindJSON(&job); err != nil {
		fmt.Println(err)
		context.AbortWithStatus(400)
		return
	}
	//fmt.Println(context.GetRawData())
	lenWordlists := len(Wordlists)

	var ws []*dwt.File

	if len(job.WordlistIds) == 0 {
		context.JSON(400, gin.H{
			"status": "Error",
			"error":  fmt.Sprintf("Wordlists have no ids"),
		})
	}

	for _, w := range job.WordlistIds {
		if w > lenWordlists-1 {

			context.JSON(400, gin.H{
				"status": "Error",
				"error":  fmt.Sprintf("Wordlist with bad id: %d", w),
			})
			return
		}
		ws = append(ws, Wordlists[w])
	}

	wlp := new(dwt.WordlistPermutations)
	wlp.Initialize(ws)
	newJob := NewJob(wlp)
	Jobs = append(Jobs, newJob)

	context.JSON(200, gin.H{
		"status": "OK",
		"code":   200,
		"id":     newJob.ID,
	})
}

func showRunningJobs(context *gin.Context) {
	context.JSON(200, gin.H{
		"status": "OK",
		"code":   200,
		"jobs":   Jobs,
	})
}
