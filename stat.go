package main

type Stat struct {
	hash          uint32
	job           uint32
	submitJob     uint32
	submitJobOK   uint32
	submitJobFail uint32
}

func (s *Stat) clear() {
	s.hash = 0
}
