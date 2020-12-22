package main

import (
	"fmt"
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

type Filter struct {
	name   string
	ranges [][2]int
}

func NewFilter(name string, ranges [][2]int) *Filter {
	return &Filter{
		name:   name,
		ranges: ranges,
	}
}

func (f *Filter) IsValid(v int) bool {
	for _, r := range f.ranges {
		if r[0] <= v && r[1] >= v {
			return true
		}
	}
	return false
}

func (f *Filter) String() string {
	return fmt.Sprintf("%s: %+v", f.name, f.ranges)
}

var (
	fre = regexp.MustCompile("^([\\w\\s]+): (\\d+)-(\\d+) or (\\d+)-(\\d+)$")
)

func parseInt(s string) int {
	n, err := strconv.Atoi(s)
	noerr(err)
	return n
}

func parseFilter(s string) *Filter {
	match := fre.FindStringSubmatch(s)
	name := match[1]
	n1, n2, n3, n4 := parseInt(match[2]), parseInt(match[3]), parseInt(match[4]), parseInt(match[5])
	ranges := [][2]int{{n1, n2}, {n3, n4}}
	return NewFilter(name, ranges)
}

type Ticket struct {
	fields []int
}

func NewTicket(f []int) *Ticket {
	return &Ticket{
		fields: f,
	}
}

func (t *Ticket) String() string {
	return fmt.Sprintf("%+v", t.fields)
}

func parseTicket(s string) *Ticket {
	chs := strings.Split(s, ",")
	fields := make([]int, 0, len(chs))
	for _, ch := range chs {
		fields = append(fields, parseInt(ch))
	}
	return NewTicket(fields)
}

func main() {
	fdata := readfile("filters")
	filters := []*Filter{}
	for _, line := range strings.Split(fdata, "\n") {
		filter := parseFilter(line)
		filters = append(filters, filter)
	}
	log.Printf("filters: %+v", filters)

	ticket := parseTicket(readfile("ticket"))
	log.Printf("ticket: %+v", ticket)

	tdata := readfile("tickets")
	tickets := []*Ticket{}
	for _, ch := range strings.Split(tdata, "\n") {
		tickets = append(tickets, parseTicket(ch))
	}
	log.Printf("tickets: %+v", tickets)

	valid := []*Ticket{}

	for _, t := range tickets {
		invcnt := 0
	Fld:
		for _, fld := range t.fields {
			for _, filter := range filters {
				//log.Printf("field: %d, filter: %s, valid: %t", fld, filter, filter.IsValid(fld))
				if filter.IsValid(fld) {
					continue Fld
				}
			}
			//log.Printf("field %d is invalid", fld)
			invcnt++
		}
		if invcnt == 0 {
			valid = append(valid, t)
		}
	}

	bind := map[string][]int{}

	for _, filter := range filters {
		log.Printf("filter: %s", filter)
		pos := map[int]struct{}{}
		for _, t := range valid {
			cnds := map[int]struct{}{}
			for ix, f := range t.fields {
				if filter.IsValid(f) {
					cnds[ix] = struct{}{}
				}
			}
			log.Printf("candidates: %+v", cnds)
			if len(pos) == 0 {
				pos = cnds
			}
			newpos := map[int]struct{}{}
			for ix := range pos {
				if _, ok := cnds[ix]; ok {
					newpos[ix] = struct{}{}
				}
			}
			log.Printf("newpos: %+v", newpos)
			pos = newpos
		}
		for ix := range pos {
			bind[filter.name] = append(bind[filter.name], ix)
		}
	}

	// map[class:[1 2] row:[2 0 1] seat:[2]]
	names := make([]string, 0, len(filters))
	for _, f := range filters {
		names = append(names, f.name)
	}
	var solve func([]string, map[string][]int, map[int]string) map[int]string
	solve = func(names []string, b map[string][]int, taken map[int]string) map[int]string {
		if len(names) == 0 {
			return taken
		}
		name := names[0]
	Ix:
		for _, ix := range b[name] {
			if _, ok := taken[ix]; ok {
				continue Ix
			}
			taken[ix] = name
			if res := solve(names[1:], b, taken); res != nil {
				return taken
			}
			delete(taken, ix)
		}
		return nil
	}
	res := solve(names, bind, map[int]string{})

	mult := 1
	cnt := 0
	for k, v := range res {
		if ok, _ := regexp.MatchString("^departure", v); ok {
			mult *= ticket.fields[k]
			cnt++
		}
	}

	log.Printf("cnt: %d, res: %d", cnt, mult)
}
