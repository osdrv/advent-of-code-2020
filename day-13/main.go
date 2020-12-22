package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func noerr(e error) {
	if e != nil {
		panic(e.Error())
	}
}

func readfile(s string) string {
	data, err := ioutil.ReadFile(s)
	noerr(err)
	return strings.TrimRight(string(data), "\r\n")
}

func parseIds(s string) []int {
	chs := strings.Split(s, ",")
	res := []int{}
	for _, ch := range chs {
		if ch == "x" {
			continue
		}
		n, err := strconv.Atoi(ch)
		noerr(err)
		res = append(res, n)
	}
	return res
}

func parseIds2(s string) [][2]int {
	chunks := strings.Split(s, ",")
	res := [][2]int{}
	for ix, ch := range chunks {
		if ch == "x" {
			continue
		}
		n, err := strconv.Atoi(ch)
		noerr(err)
		res = append(res, [2]int{ix, n})
	}
	return res
}

func main() {
	lines := strings.Split(readfile("input"), "\n")
	ids := parseIds2(lines[1])
	log.Printf("%+v", ids)
	jmp := ids[0][1]
	ts := 0
	for _, pair := range ids[1:] {
		delta := pair[0]
		id := pair[1]
		for (ts+delta)%id != 0 {
			ts += jmp
		}
		jmp *= id
	}

	log.Printf("%d", ts)
}
