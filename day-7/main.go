package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func readfile(f string) ([]byte, error) {
	return ioutil.ReadFile(f)
}

func noerr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

var (
	//re = regexp.MustCompile("^(\\w+)\\s(\\w+)\\s bags\\s contain\\s((no\\s other\\s bags | (\\d+)\\s(\\w+)\\s(\\w+)\\s(bag|bags)(,\\s)?)+)\\.$")
	reObj  = regexp.MustCompile("^(\\w+ \\w+) bags contain")
	reSubj = regexp.MustCompile("(\\d+) (\\w+) (\\w+)|(no other)")
)

type Contains struct {
	subj  string
	quant int
}

func main() {
	data, err := readfile("input")
	noerr(err)
	rules := strings.Split(strings.TrimRight(string(data), "\r\n"), "\n")

	graph := make(map[string][]Contains)

	for _, rule := range rules {
		log.Printf("rule: %q", rule)
		res := reObj.FindStringSubmatch(rule)
		log.Printf("matches: %+v", res)
		obj := strings.TrimRight(res[1], " ")
		rule = rule[len(res[0]):]

		log.Printf("obj: %q, rule: %q", obj, rule)

		subjMatch := reSubj.FindAllStringSubmatch(rule, -1)
		graph[obj] = make([]Contains, 0)

		for _, sm := range subjMatch {
			if sm[0] == "no other" {
				continue
			}
			subj := strings.TrimRight(sm[2]+" "+sm[3], " ")
			quant, err := strconv.Atoi(sm[1])
			noerr(err)
			log.Printf("Subj: %q", subj)
			graph[obj] = append(graph[obj], Contains{subj: subj, quant: quant})
		}
	}

	log.Printf("graph: %+v", graph)

	var visit func(string) int
	visit = func(s string) int {
		if len(graph[s]) == 0 {
			return 1
		}
		res := 1
		for _, ch := range graph[s] {
			i := ch.quant * visit(ch.subj)
			log.Printf("%s bag contains %d %s bags (%d)", s, ch.quant, ch.subj, i)
			res += i
		}
		return res
	}

	answer := visit("shiny gold") - 1

	log.Printf("the answer is: %d", answer)

	//ptr := "dark orange"
	//i := 0
	//for ptr != graph[ptr] {
	//	i++
	//	ptr = graph[ptr]
	//}
	//log.Printf("the answer is: %d", i)
}
