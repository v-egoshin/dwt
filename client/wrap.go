package client

type Job struct {
	//TODO: client job struct
}

type Jobs []Job

type Runner struct {
	Arch      string
	OS        string
	MaxMemory uint32
	MaxSpace  uint32
}

func Register() {
	// register on server
}
