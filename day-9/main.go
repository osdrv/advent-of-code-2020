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

func solve1(ns []uint64, w int) uint64 {
	ix := w
Ix:
	for ix < len(ns) {
		for i := ix - w; i < ix; i++ {
			for j := i + 1; j < ix; j++ {
				if ns[i]+ns[j] == ns[ix] {
					ix++
					continue Ix
				}
			}
		}
		return ns[ix]
	}
	return 0
}

func min(ns []uint64) uint64 {
	res := ns[0]
	for _, n := range ns {
		if n < res {
			res = n
		}
	}
	return res
}

func max(ns []uint64) uint64 {
	res := ns[0]
	for _, n := range ns {
		if n > res {
			res = n
		}
	}
	return res
}

func solve2(ns []uint64, sum uint64) uint64 {
	ws := ns[0] + ns[1]
	l, h := 0, 1
	log.Printf("sum: %d", sum)
	for h < len(ns) {
		if ws > sum {
			ws -= ns[l]
			l++
		} else if ws < sum {
			h++
			ws += ns[h]
		} else {
			a1, a2 := min(ns[l:h+1]), max(ns[l:h+1])
			log.Printf("nums: %+v, min: %d, max: %d", ns[l:h+1], a1, a2)
			return a1 + a2
		}
	}
	return 0
}

func main() {
	lines := strings.Split(readfile("input"), "\n")
	nums := make([]uint64, 0, len(lines))
	for _, line := range lines {
		n, err := strconv.ParseUint(line, 10, 64)
		noerr(err)
		nums = append(nums, n)
	}

	res1 := solve1(nums, 25)
	log.Printf("the result1 is: %d", res1)

	res2 := solve2(nums, res1)
	log.Printf("the result2 is: %d", res2)
}
