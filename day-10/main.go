package main

import (
	"io/ioutil"
	"log"
	"sort"
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

func main() {
	lines := strings.Split(readfile("input"), "\n")
	nums := make([]int, 0, len(lines))
	for _, line := range lines {
		n, err := strconv.Atoi(line)
		noerr(err)
		nums = append(nums, n)
	}

	var deltas [3]int
	deltas[2] = 1
	sort.Ints(nums)

	prev := 0
	for _, n := range nums {
		d := n - prev
		deltas[d-1]++
		prev = n
	}

	log.Printf("deltas: %+v", deltas)

	//log.Printf("the answer is: %d", deltas[0]*deltas[2])

	var solve func(int) int
	nums = append([]int{0}, nums...)
	nums = append(nums, nums[len(nums)-1]+3)
	max := len(nums) - 1
	memo := make(map[int]int)
	solve = func(ix int) int {
		log.Printf("solving for ix: %d(%d)", ix, max)
		if ix == max {
			return 1
		}
		if v, ok := memo[ix]; ok {
			return v
		}
		res := 0
		if ix+1 <= max && nums[ix+1]-nums[ix] <= 3 {
			res += solve(ix + 1)
		}
		if ix+2 <= max && nums[ix+2]-nums[ix] <= 3 {
			res += solve(ix + 2)
		}
		if ix+3 <= max && nums[ix+3]-nums[ix] <= 3 {
			res += solve(ix + 3)
		}
		memo[ix] = res
		return res
	}

	res := solve(0)

	log.Printf("The result is: %d", res)
}
