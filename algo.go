package main

const (
	UNKNOWN = iota
	CN
	CN_LITE
	CN_HEAVY
	CN_PICO
	RANDOM_X
	ARGON2
	ASTROBWT
	KAWPOW
)

type Algorithm struct {
}

func NewAlgorithm(algo string) *Algorithm {
	a := &Algorithm{}
	// TODO
	return a
}

func (a *Algorithm) family() int {
	return UNKNOWN
}

func (a *Algorithm) name() string {
	return ""
}
