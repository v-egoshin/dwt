package server

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/v-egoshin/dwt"
)

//https://stackoverflow.com/questions/62906766/how-to-group-routes-in-gin

func InitializeRoutes(g *gin.RouterGroup) *gin.RouterGroup {
	v1 := g.Group("/")
	{
		manage := v1.Group("/manage")
		{
			manage.GET("/job/:job_id", ListJobById)
			manage.GET("/job/:job_id/cancel", CancelJobById)
			manage.GET("/job", ListJobs)
			manage.POST("/job/create", CreateNewJob)

			manage.GET("/wordlist", ListWordlists)
			manage.POST("/wordlist/upload", UploadWordlist)
			manage.GET("/wordlist/reindex", ReindexWordlists)
		}

		jobs := v1.Group("/runner")
		{
			jobs.POST("/runner", RegisterRunner) // register client seat for jobs
			jobs.GET("/:job_id/get/:number", GetJobChunk)
			jobs.POST("/poll", GetJob)
		}

	}
	return v1
}

func UploadWordlist(context *gin.Context) {

}

func CancelJobById(context *gin.Context) {
	cjob, notfound := FindJobByUID(context)
	if notfound {
		return
	}

	cjob.Delete()

	context.JSON(200, gin.H{
		"status": "OK",
		"code":   200,
	})
}

func GetJob(context *gin.Context) {

}

func RegisterRunner(context *gin.Context) {
	// register
}

func ReindexWordlists(context *gin.Context) {
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

func ListWordlists(context *gin.Context) {
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

func GetJobChunk(context *gin.Context) {
	var number uint32
	if cnum, err := strconv.Atoi(context.Param("number")); err == nil {
		number = uint32(cnum)
	} else {
		context.JSON(404, gin.H{
			"status": "Error",
			"error":  fmt.Sprintf("Bad number: %d", number),
		})
		return
	}

	cjob, notfound := FindJobByUID(context)
	if notfound {
		return
	}

	context.JSON(200, gin.H{
		"status": "Ok",
		"data":   cjob.Get(number),
	})
}

func FindJobByUID(context *gin.Context) (*Job, bool) {

	id := context.Param("job_id")

	uid, err := uuid.FromString(id)
	if err != nil {
		context.JSON(404, gin.H{
			"status": "Error",
			"error":  fmt.Sprintf("Job %s not found", id),
		})
		return nil, true
	}

	cjob := new(Job)
	for _, j := range Jobs {
		if j.ID == uid {
			cjob = j
		}
	}
	return cjob, false
}

func ListJobById(context *gin.Context) {
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

func CreateNewJob(context *gin.Context) {
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

func ListJobs(context *gin.Context) {
	context.JSON(200, gin.H{
		"status": "OK",
		"code":   200,
		"jobs":   Jobs,
	})
}
