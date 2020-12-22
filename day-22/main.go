package main

import (
	"io/ioutil"
	"log"
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

func readints(s string) []int {
	lines := strings.Split(readfile(s), "\n")
	res := make([]int, 0, len(lines))
	for _, l := range lines {
		n, err := strconv.Atoi(l)
		noerr(err)
		res = append(res, n)
	}
	return res
}

func solve1() {
	p1, p2 := readints("player1_sample"), readints("player2_sample")
	var winner []int
	var c1, c2 int
	for {
		if len(p1) == 0 {
			winner = p2
			break
		} else if len(p2) == 0 {
			winner = p1
			break
		}
		c1, p1 = p1[0], p1[1:]
		c2, p2 = p2[0], p2[1:]
		if c1 > c2 {
			p1 = append(p1, c1, c2)
		} else {
			p2 = append(p2, c2, c1)
		}
	}

	res := 0
	for i := 1; i <= len(winner); i++ {
		res += winner[len(winner)-i] * i
	}

	log.Printf("the result is: %d", res)
}

func mkcp(src []int) []int {
	cp := make([]int, len(src))
	copy(cp, src)
	return cp
}

func makesum(in []int) int {
	res := 0
	for i := 1; i <= len(in); i++ {
		res += in[len(in)-i] * i
	}
	return res
}

func playGame(g int, p1, p2 []int) ([]int, []int) {
	log.Printf("game %d starts", g)
	memo := make(map[[2]int]struct{})
	var c1, c2 int
	r := 0
	for {
		r++
		log.Printf("round %d of game %d starts: p1: %+v, p2: %+v", r, g, p1, p2)
		if len(p1) == 0 {
			return nil, p2
		} else if len(p2) == 0 {
			return p1, nil
		}
		sum1, sum2 := makesum(p1), makesum(p2)
		if _, ok := memo[[2]int{sum1, sum2}]; ok {
			return p1, nil
		}
		memo[[2]int{sum1, sum2}] = struct{}{}
		c1, p1 = p1[0], p1[1:]
		c2, p2 = p2[0], p2[1:]
		if len(p1) >= c1 && len(p2) >= c2 {
			g++
			pp1, pp2 := playGame(g, mkcp(p1[:c1]), mkcp(p2[:c2]))
			if len(pp1) > 0 {
				p1 = append(p1, c1, c2)
			} else if len(pp2) > 0 {
				p2 = append(p2, c2, c1)
			} else {
				panic("no way")
			}
			continue
		}
		if c1 > c2 {
			p1 = append(p1, c1, c2)
		} else {
			p2 = append(p2, c2, c1)
		}
	}
}

func main() {
	p1, p2 := readints("player1"), readints("player2")
	pp1, pp2 := playGame(1, p1, p2)
	var res int
	if len(pp1) == 0 {
		res = makesum(pp2)
	} else if len(pp2) == 0 {
		res = makesum(pp1)
	}
	log.Printf("the answer is: %d", res)
}
