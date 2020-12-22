package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Instr struct {
	n string
	v int
}

func Parse(s string) Instr {
	n := s[:3]
	v, err := strconv.Atoi(s[4:])
	noerr(err)
	return Instr{n: n, v: v}
}

type State struct {
	instr []Instr
	acc   int
	ptr   int
}

func NewState(instr []Instr) *State {
	return &State{
		instr: instr,
		acc:   0,
		ptr:   0,
	}
}

func (s *State) Next() (int, int) {
	if s.ptr < 0 || s.ptr >= len(s.instr) {
		return -1, s.acc
	}
	i := s.instr[s.ptr]
	switch i.n {
	case "nop":
		s.ptr++
	case "jmp":
		s.ptr += i.v
	case "acc":
		s.acc += i.v
		s.ptr++
	}
	return s.ptr, s.acc
}

func main() {
	s := readfile("input")
	lines := strings.Split(s, "\n")
	instr := make([]Instr, 0, len(lines))
	for _, l := range lines {
		instr = append(instr, Parse(l))
	}
Replace:
	for i := 0; i < len(instr); i++ {
		if instr[i].n != "jmp" && instr[i].n != "nop" {
			continue
		}
		memo := make(map[int]struct{})
		instrcp := make([]Instr, len(instr))
		copy(instrcp, instr)
		switch instrcp[i].n {
		case "jmp":
			instrcp[i].n = "nop"
		case "nop":
			instrcp[i].n = "jmp"
		}
		st := NewState(instrcp)
		for {
			ptr, acc := st.Next()
			if _, ok := memo[ptr]; ok {
				log.Printf("loop at instruction: %d", ptr)
				break
			}
			memo[ptr] = struct{}{}
			if ptr == -1 {
				log.Printf("the program terminates successfully with acc: %d", acc)
				break Replace
			}
		}
	}
}

func readfile(f string) string {
	data, err := ioutil.ReadFile(f)
	noerr(err)
	return strings.TrimRight(string(data), "\r\n")
}

func noerr(e error) {
	if e != nil {
		panic(e.Error())
	}
}
