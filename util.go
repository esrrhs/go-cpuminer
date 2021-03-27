package main

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// DiffToTarget converts a whole number difficulty into a target.
func DiffToTarget(diff float64, powLimit *big.Int) (*big.Int, error) {
	if diff <= 0 {
		return nil, fmt.Errorf("invalid pool difficulty %v (0 or less than "+
			"zero passed)", diff)
	}

	// Round down in the case of a non-integer diff since we only support
	// ints (unless diff < 1 since we don't allow 0)..
	if diff < 1 {
		diff = 1
	} else {
		diff = math.Floor(diff)
	}
	divisor := new(big.Int).SetInt64(int64(diff))
	max := powLimit
	target := new(big.Int)
	target.Div(max, divisor)

	return target, nil
}

// reverseS reverses a hex string.
func reverseS(s string) (string, error) {
	a := strings.Split(s, "")
	sRev := ""
	if len(a)%2 != 0 {
		return "", fmt.Errorf("Incorrect input length")
	}
	for i := 0; i < len(a); i += 2 {
		tmp := []string{a[i], a[i+1], sRev}
		sRev = strings.Join(tmp, "")
	}
	return sRev, nil
}

// ReverseToInt reverse a string and converts to int32.
func ReverseToInt(s string) (int32, error) {
	sRev, err := reverseS(s)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(sRev, 10, 32)
	return int32(i), err
}

// RevHash reverses a hash in string format.
func RevHash(hash string) string {
	revHash := ""
	for i := 0; i < 7; i++ {
		j := i * 8
		part := fmt.Sprintf("%c%c%c%c%c%c%c%c",
			hash[6+j], hash[7+j], hash[4+j], hash[5+j],
			hash[2+j], hash[3+j], hash[0+j], hash[1+j])
		revHash += part
	}
	return revHash
}
