package main

import (
	"io/ioutil"
	"log"
	"strings"
)

func readfile(s string) ([]byte, error) {
	return ioutil.ReadFile(s)
}

func noerr(e error, m string) {
	if e != nil {
		log.Fatalf(m, e)
	}
}

func main() {
	data, err := readfile("input")
	noerr(err, "failed top open input file")

	groups := strings.Split(
		strings.TrimRight(string(data), "\n"),
		"\n\n",
	)

	res := 0

	for _, group := range groups {
		answers := cntAnswers(group)
		res += len(answers)
	}

	log.Printf("the answer is: %d", res)
}

func cntAnswers(s string) map[rune]int {
	answers := make(map[rune]int)
	ppl := 1
	for _, ch := range s {
		if ch == '\n' {
			ppl++
			continue
		}
		answers[ch]++
	}
	res := make(map[rune]int)
	for k, v := range answers {
		if v == ppl {
			res[k] = v
		}
	}
	return res
}
