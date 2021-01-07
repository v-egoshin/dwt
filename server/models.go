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
	Position    int
	Done        bool
	Start       time.Time
	End         time.Time
	Permutate   *dwt.WordlistPermutations
	*Authenticate
}

func NewJob(wlp *dwt.WordlistPermutations) *Job {
	return &Job{
		ID:           uuid.NewV4(),
		Description:  "",
		Position:     0,
		Done:         false,
		Start:        time.Now(),
		End:          time.Time{},
		Permutate:    wlp,
		Authenticate: NewAuthenticate(),
	}
}
