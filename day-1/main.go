package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

const (
	SUM = 2020
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)

	mem := make(map[int]struct{})

	for sc.Scan() {
		line := sc.Text()
		log.Printf("line: %q", line)
		num, err := strconv.Atoi(line)
		if err != nil {
			log.Fatalf("failed to parse line: %q", line)
		}
		mem[num] = struct{}{}
	}

	if err := sc.Err(); err != nil {
		log.Fatalf("scanner error: %s", err)
	}

	for num := range mem {
		a, b, ok := findPair(SUM-num, mem, map[int]struct{}{num: struct{}{}})
		if ok {
			log.Printf("the result is: %d", a*b*num)
			return
		}
	}
}

func findPair(sum int, mem map[int]struct{}, taken map[int]struct{}) (int, int, bool) {
	for num := range mem {
		if _, ok := taken[num]; ok {
			continue
		}
		if _, ok := mem[sum-num]; ok {
			return num, sum - num, true
		}
	}
	return 0, 0, false
}
