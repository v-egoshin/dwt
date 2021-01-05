package server

// Configuration
// Auth
// Jobs

import (
	"time"

	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Configuration struct {
	gorm.Model
	Auth string
}

type Authenticate struct {
	gorm.Model
	Access int    // access level: read, write
	Token  string // for manipulate wordlists and access to jobs
}

type Job struct {
	gorm.Model
	ID          uuid.UUID
	Description string
	Position    int
	Done        bool
	Start       time.Time
	End         time.Time
	Authenticate
}
