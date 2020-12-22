package main

import (
	"io/ioutil"
	"log"
	"regexp"
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

func parseMask(s string) Mask {
	var mor, mand uint64
	s = s[7:]
	for i := 35; i >= 0; i-- {
		switch s[i] {
		case '1':
			mor |= 1 << (35 - i)
		case '0':
		default:
			mand |= 1 << (35 - i)
		}
	}
	return Mask{
		Mor:  mor,
		Mand: mand,
	}
}

func parseMaskV2(s string) MaskV2 {
	s = s[7:]
	return MaskV2{Orig: s}
}

var (
	re = regexp.MustCompile("^mem\\[(\\d+)\\] = (\\d+)$")
)

func parseMem(s string) Mem {
	mtch := re.FindStringSubmatch(s)
	log.Printf("%+v", mtch)
	astr, vstr := mtch[1], mtch[2]
	addr, err := strconv.Atoi(astr)
	noerr(err)
	val, err := strconv.ParseUint(vstr, 10, 64)
	noerr(err)
	return Mem{
		Addr: addr,
		Val:  val,
	}
}

type Mask struct {
	Mor, Mand uint64
}

type MaskV2 struct {
	Orig string
}

type Mem struct {
	Addr int
	Val  uint64
}

type Memory struct {
	Regs   map[int]uint64
	Mask   Mask
	MaskV2 MaskV2
}

func NewMemory() *Memory {
	return &Memory{
		Regs: make(map[int]uint64),
	}
}

func (m *Memory) Exec(i interface{}) {
	switch instr := i.(type) {
	case Mask:
		m.Mask = instr
	case Mem:
		v := instr.Val
		v &= m.Mask.Mand
		v |= m.Mask.Mor
		m.Regs[instr.Addr] = v
		log.Printf("mem set: %d -> %036b", instr.Addr, v)
	}
}

func (m *Memory) ExecV2(i interface{}) {
	switch instr := i.(type) {
	case MaskV2:
		m.MaskV2 = instr
	case Mem:
		addrs := []int{instr.Addr}
		mask := m.MaskV2.Orig
		for off := 35; off >= 0; off-- {
			switch mask[off] {
			case '0':
			case '1':
				for i := range addrs {
					addrs[i] |= 1 << (35 - off)
				}
			case 'X':
				tmp := make([]int, 0, len(addrs)*2)
				for _, addr := range addrs {
					tmp = append(tmp, addr&(^(1 << (35 - off))))
					tmp = append(tmp, addr|(1<<(35-off)))
				}
				addrs = tmp
			}
		}
		for _, addr := range addrs {
			log.Printf("Addr: %035b", addr)
			m.Regs[addr] = instr.Val
		}
		log.Printf("========")
	}
}

func main() {
	lines := strings.Split(readfile("input"), "\n")
	instrs := make([]interface{}, 0, len(lines))
	for _, line := range lines {
		switch line[:3] {
		case "mem":
			instrs = append(instrs, parseMem(line))
		case "mas":
			instrs = append(instrs, parseMaskV2(line))
		}
	}

	mem := NewMemory()
	for _, instr := range instrs {
		mem.ExecV2(instr)
	}

	var sum uint64
	for _, v := range mem.Regs {
		sum += v
	}

	log.Printf("the answer is: %d", sum)
}
