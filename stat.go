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
	s.job = 0
	s.submitJob = 0
	s.submitJobOK = 0
	s.submitJobFail = 0
}
