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

type Ceil struct {
	x, y, z, w int
}

type Area struct {
	rangex, rangey, rangez, rangew [2]int
	cubes                          map[Ceil]struct{}
}

func NewArea(rangex, rangey, rangez, rangew [2]int) *Area {
	return &Area{
		rangex: rangex, rangey: rangey, rangez: rangez, rangew: rangew,
		cubes: make(map[Ceil]struct{}),
	}
}

func (a *Area) IsActive(c Ceil) bool {
	_, ok := a.cubes[c]
	return ok
}

func (a *Area) Activate(c Ceil) {
	a.cubes[c] = struct{}{}
}

func (a *Area) Deactivate(c Ceil) {
	delete(a.cubes, c)
}

func (a *Area) String() string {
	return fmt.Sprintf("area:\nrangeX: %+v\nrangeY:%+v\nrangeZ:%+v\ncubes: %+v", a.rangex, a.rangey, a.rangez, a.cubes)
}

func neighbors(c Ceil) []Ceil {
	res := make([]Ceil, 0, 80)
	for _, dx := range []int{-1, 0, 1} {
		for _, dy := range []int{-1, 0, 1} {
			for _, dz := range []int{-1, 0, 1} {
				for _, dw := range []int{-1, 0, 1} {
					nc := Ceil{x: c.x + dx, y: c.y + dy, z: c.z + dz, w: c.w + dw}
					if nc != c {
						res = append(res, nc)
					}
				}
			}
		}
	}
	return res
}

func evolve(area *Area) *Area {
	newarea := NewArea(
		[2]int{area.rangex[0] - 1, area.rangex[1] + 1},
		[2]int{area.rangey[0] - 1, area.rangey[1] + 1},
		[2]int{area.rangez[0] - 1, area.rangez[1] + 1},
		[2]int{area.rangew[0] - 1, area.rangew[1] + 1},
	)
	for x := area.rangex[0] - 1; x <= area.rangex[1]+1; x++ {
		for y := area.rangey[0] - 1; y <= area.rangey[1]+1; y++ {
			for z := area.rangez[0] - 1; z <= area.rangez[1]+1; z++ {
				for w := area.rangew[0] - 1; w <= area.rangew[1]+1; w++ {
					c := Ceil{x: x, y: y, z: z, w: w}
					nbrs := neighbors(c)
					actcnt := 0
					for _, nb := range nbrs {
						if area.IsActive(nb) {
							actcnt++
						}
					}
					if area.IsActive(c) {
						if actcnt >= 2 && actcnt <= 3 {
							newarea.Activate(c)
						}
					} else {
						if actcnt == 3 {
							newarea.Activate(c)
						}
					}
				}
			}
		}
	}
	return newarea
}

func main() {
	lines := strings.Split(readfile("input"), "\n")
	ceils := make(map[Ceil]struct{})
	for y, line := range lines {
		for x, ch := range line {
			if ch == '#' {
				ceils[Ceil{x: x, y: y, z: 0, w: 0}] = struct{}{}
			}
		}
	}
	area := NewArea(
		[2]int{0, len(lines[0]) - 1},
		[2]int{0, len(lines) - 1},
		[2]int{0, 0},
		[2]int{0, 0},
	)
	for c := range ceils {
		area.Activate(c)
	}

	log.Printf("area: %s", area)

	for i := 0; i < 6; i++ {
		area = evolve(area)
		log.Printf("new generation: %s", area)
	}

	log.Printf("active cubes: %d", len(area.cubes))
}
