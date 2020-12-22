package main

import (
	"fmt"
	"io/ioutil"
	"log"
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

func main() {
	s := readfile("input")
	lines := strings.Split(s, "\n")
	field := make([][]rune, len(lines))
	for ix, line := range lines {
		field[ix] = []rune(line)
	}

	width, height := len(field[0]), len(field)

	evolve := func() int {
		diff := 0
		cp := copyfield(field)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				switch cp[y][x] {
				case 'L':
					if neighbours2(field, x, y) == 0 {
						cp[y][x] = '#'
						diff++
					}
				case '#':
					if neighbours2(field, x, y) >= 5 {
						cp[y][x] = 'L'
						diff++
					}
				}
			}
		}
		field = cp
		printfield(field)
		return diff
	}

	for evolve() > 0 {
	}
	log.Printf("ocupied seats: %d", countseats(field))
}

var (
	deltas = [][2]int{
		{1, 0},
		{1, 1},
		{0, 1},
		{-1, 1},
		{-1, 0},
		{-1, -1},
		{0, -1},
		{1, -1},
	}
)

//neighbours := func(x, y int) int {
//	cnt := 0
//	for _, d := range deltas {
//		x1, y1 := x+d[0], y+d[1]
//		if x1 >= 0 && y1 >= 0 && x1 < width && y1 < height {
//			if field[y1][x1] == '#' {
//				cnt++
//			}
//		}
//	}
//	return cnt
//}

func neighbours2(field [][]rune, x, y int) int {
	cnt := 0
	width, height := len(field[0]), len(field)
Delta:
	for _, d := range deltas {
		x1, y1 := x+d[0], y+d[1]
		for x1 >= 0 && y1 >= 0 && x1 < width && y1 < height {
			switch field[y1][x1] {
			case 'L':
				continue Delta
			case '#':
				cnt++
				continue Delta
			default:
				x1 += d[0]
				y1 += d[1]
			}
		}
	}
	return cnt
}

func countseats(field [][]rune) int {
	cnt := 0
	for y := 0; y < len(field); y++ {
		for x := 0; x < len(field[0]); x++ {
			if field[y][x] == '#' {
				cnt++
			}
		}
	}
	return cnt
}

func printfield(field [][]rune) {
	var b strings.Builder
	for _, line := range field {
		for _, ch := range line {
			b.WriteRune(ch)
			b.WriteRune(' ')
		}
		b.WriteRune('\n')
	}
	b.WriteRune('\n')
	fmt.Print(b.String())
}

func copyfield(field [][]rune) [][]rune {
	cp := make([][]rune, len(field))
	for ix, line := range field {
		lcp := make([]rune, len(line))
		copy(lcp, line)
		cp[ix] = lcp
	}
	return cp
}
