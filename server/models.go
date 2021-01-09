package server

// Configuration
// Auth
// Jobs

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/v-egoshin/dwt"
	"gorm.io/gorm"
)

type Configuration struct {
	gorm.Model
	Auth string
}

type Role int

const (
	Read Role = iota
	Write
)

type Authenticate struct {
	gorm.Model
	Access Role      // access level: read, write
	Token  uuid.UUID // for manipulate wordlists and access to jobs
}

func NewAuthenticate() *Authenticate {
	return &Authenticate{
		Access: Read,
		Token:  uuid.NewV4(),
	}
}

type Job struct {
	gorm.Model
	ID          uuid.UUID
	Description string
	Position    uint32
	Done        bool
	Start       time.Time
	End         time.Time
	Permutes    *dwt.WordlistPermutations
	receiver    chan []uint32
	*Authenticate
}

func NewJob(wlp *dwt.WordlistPermutations) *Job {
	return &Job{
		ID:           uuid.NewV4(),
		Description:  "",
		Position:     0,
		Done:         false,
		receiver:     make(chan []uint32, 1),
		Start:        time.Now(),
		End:          time.Time{},
		Permutes:     wlp,
		Authenticate: NewAuthenticate(),
	}
}

func (j *Job) Get(i uint32) [][]string {

	if j.Position == 0 {
		go j.Permutes.PermuteAll(j.receiver)
	}

	if (j.Position + i) > j.Permutes.Count {
		i = j.Permutes.Count - (j.Position + i)
	}
	var collect [][]string

	var c uint32

	for {
		pair, ok := <-j.receiver
		if !ok {
			break
		}
		if c > i {
			break
		}
		r, _ := j.Permutes.GetPermuteByState(pair)
		collect = append(collect, r)
		c += 1
		j.Position += 1
	}

	return collect
}
