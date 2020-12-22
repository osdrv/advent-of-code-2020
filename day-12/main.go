package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Vector struct {
	X, Y int
}

type Instr struct {
	Cmd rune
	Val int
}

func (i Instr) String() string { return fmt.Sprintf("%c%d", i.Cmd, i.Val) }

type Ship struct {
	WP  Vector
	Pos Vector
}

func (s *Ship) Exec(i Instr) {
	log.Printf("executing instr: %s", i)
	switch i.Cmd {
	case 'N':
		s.WP.Y += i.Val
	case 'S':
		s.WP.Y -= i.Val
	case 'E':
		s.WP.X += i.Val
	case 'W':
		s.WP.X -= i.Val
	case 'L', 'R':
		d := i.Val
		if i.Cmd == 'L' {
			d = 360 - d
		}
		switch d {
		case 90:
			var wp Vector
			wp.X = s.WP.Y
			wp.Y = -s.WP.X
			s.WP = wp
		case 180:
			var wp Vector
			wp.X = -s.WP.X
			wp.Y = -s.WP.Y
			s.WP = wp
		case 270:
			var wp Vector
			wp.X = -s.WP.Y
			wp.Y = s.WP.X
			s.WP = wp
		}
	case 'F':
		for n := 0; n < i.Val; n++ {
			s.Pos.X += s.WP.X
			s.Pos.Y += s.WP.Y
		}
	default:
		panic("unexpected")
	}
	log.Printf("%d %d %d %d", s.Pos.X, s.Pos.Y, s.WP.X, s.WP.Y)
	//log.Printf("the ship pos: %+v, waypoint: %+v", s.Pos, s.WP)
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func main() {
	lines := strings.Split(readfile("input"), "\n")
	instrs := make([]Instr, 0, len(lines))
	for _, line := range lines {
		instrs = append(instrs, parseinstr(line))
	}
	s := &Ship{
		WP: Vector{
			X: 10,
			Y: 1,
		},
	}
	for _, i := range instrs {
		s.Exec(i)
	}
	log.Printf("%+v", s)
	log.Printf("Manhattan distance is: %d (%d + %d)", abs(s.Pos.X)+abs(s.Pos.Y), s.Pos.X, s.Pos.Y)
}

func parseinstr(s string) Instr {
	var i Instr
	i.Cmd = rune(s[0])
	v, err := strconv.Atoi(s[1:])
	noerr(err)
	i.Val = v
	return i
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
