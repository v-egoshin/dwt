package client

import "time"

type Job struct {
	//TODO: client job struct
}

type Runner struct {
	Arch      string
	OS        string
	MaxMemory uint32
	MaxSpace  uint32
	LastSeen  time.Time
}

func Register() {
	// register on server
}
