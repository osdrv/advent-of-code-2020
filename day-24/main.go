package main

import (
	"io/ioutil"
	"log"
	"strings"
)

type Step int

func (s Step) String() string {
	switch s {
	case E:
		return "E"
	case SE:
		return "SE"
	case SW:
		return "SW"
	case W:
		return "W"
	case NW:
		return "NW"
	case NE:
		return "NE"
	default:
		return "UNKNWN"
	}
}

const (
	E Step = iota
	SE
	SW
	W
	NW
	NE
)

var (
	STEPS = map[Step][3]int{
		E:  {1, -1, 0},
		W:  {-1, 1, 0},
		SE: {0, -1, 1},
		SW: {-1, 0, 1},
		NE: {1, 0, -1},
		NW: {0, 1, -1},
	}
)

func parsePath(p string) []Step {
	res := []Step{}
	i := 0
	for i < len(p) {
		var step Step
		switch p[i] {
		case 'e':
			step = E
		case 'w':
			step = W
		case 's':
			i++
			switch p[i] {
			case 'e':
				step = SE
			case 'w':
				step = SW
			}
		case 'n':
			i++
			switch p[i] {
			case 'e':
				step = NE
			case 'w':
				step = NW
			}
		}
		res = append(res, step)
		i++
	}
	return res
}

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

func resolvePath(p []Step) [3]int {
	var pos [3]int
	for _, s := range p {
		dp := STEPS[s]
		pos[0] += dp[0]
		pos[1] += dp[1]
		pos[2] += dp[2]
	}
	return pos
}

func evolve(floor map[[3]int]bool) map[[3]int]bool {
	var minx, maxx, miny, maxy, minz, maxz int
	for pos := range floor {
		x, y, z := pos[0], pos[1], pos[2]
		if x < minx {
			minx = x
		}
		if x > maxx {
			maxx = x
		}
		if y < miny {
			miny = y
		}
		if y > maxy {
			maxy = y
		}
		if z < minz {
			minz = z
		}
		if z > maxz {
			maxz = z
		}
	}

	minx--
	miny--
	minz--
	maxx++
	maxy++
	maxz++

	newfloor := make(map[[3]int]bool)

	for i := minz; i <= maxz; i++ {
		for j := miny; j <= maxy; j++ {
			for k := minx; k <= maxx; k++ {
				pos := [3]int{k, j, i}
				var isBlack bool // white
				if b, ok := floor[pos]; ok {
					isBlack = b
				}
				if isBlack {
					nbrs := countNbrs(pos, floor)
					if nbrs == 0 || nbrs > 2 {
						isBlack = false
					}
				} else {
					if countNbrs(pos, floor) == 2 {
						isBlack = true
					}
				}
				if isBlack {
					newfloor[pos] = isBlack
				}
			}
		}
	}

	return newfloor
}

func countNbrs(p [3]int, floor map[[3]int]bool) int {
	cnt := 0
	for _, step := range STEPS {
		pp := [3]int{p[0] + step[0], p[1] + step[1], p[2] + step[2]}
		b, ok := floor[pp]
		if !ok {
			continue
		}
		if b {
			cnt++
		}
	}
	return cnt
}

func countPop(floor map[[3]int]bool) int {
	blackCnt := 0
	for _, isBlack := range floor {
		if isBlack {
			blackCnt++
		}
	}
	return blackCnt
}

func main() {
	lines := strings.Split(readfile("input"), "\n")
	paths := make([][]Step, 0, len(lines))
	for _, line := range lines {
		paths = append(paths, parsePath(line))
	}

	// white = false, black = true
	tiles := make(map[[3]int]bool)

	for _, pp := range paths {
		pos := resolvePath(pp)
		tiles[pos] = !tiles[pos]
	}
	log.Printf("%+v", paths)

	blackCnt := countPop(tiles)

	log.Printf("the result is: %d", blackCnt)

	for i := 1; i <= 100; i++ {
		log.Printf("===== day %d =====", i)
		tiles = evolve(tiles)
		log.Printf("count: %d", countPop(tiles))
	}
}
