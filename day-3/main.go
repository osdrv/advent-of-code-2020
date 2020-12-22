package main

import (
	"io/ioutil"
	"log"
	"strings"
)

func findSlopes(lines []string, right, down int) int {
	x, y := 0, 0
	treeCnt := 0
	width := len(lines[0])
	for y < len(lines) {
		if lines[y][x] == '#' {
			treeCnt++
		}
		x += right
		y += down
		if x >= width {
			x %= width
		}
	}
	return treeCnt
}

func main() {
	data, err := ioutil.ReadFile("input")
	if err != nil {
		log.Fatalf("failed to read input file: %s", err)
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")

	steps := [][2]int{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}
	mult := 1
	for _, step := range steps {
		mult *= findSlopes(lines, step[0], step[1])
	}

	log.Printf("tree count mult: %d", mult)
}
