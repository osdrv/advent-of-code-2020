package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func noErr(err error) {
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}
}

func main() {
	data, err := ioutil.ReadFile("input")
	if err != nil {
		log.Fatalf("failed to open input: %s", err)
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")

	taken := make([]int, 1024)
	for _, line := range lines {
		seatId := getSeatId(line)
		taken[seatId] = seatId
	}

	log.Printf("taken: %+v", taken)

	for i := 1; i < len(taken)-1; i++ {
		if taken[i] == 0 && taken[i-1] > 1 && taken[i+1] > 0 {
			log.Printf("seat candidate: %d", i)
		}
	}
}

func getSeatId(line string) int {
	vert, hor := line[:7], line[7:]
	log.Printf("vert: %q, hor: %q", vert, hor)
	v, err := convToInt(vert, 'F', 'B')
	noErr(err)
	h, err := convToInt(hor, 'L', 'R')
	noErr(err)
	log.Printf("v: %d, h: %d", v, h)
	return v*8 + h
}

func convToInt(s string, falseCh, trueCh byte) (int, error) {
	res := 0
	pow := 1
	for i := len(s) - 1; i >= 0; i-- {
		switch s[i] {
		case falseCh:
		case trueCh:
			res += pow
		default:
			return 0, fmt.Errorf("unexpected char: %c", s[i])
		}
		pow *= 2
	}
	return res, nil
}
