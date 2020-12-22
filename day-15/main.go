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

type Iter struct {
	nums []int
	last int
	memo map[int][2]int
}

func NewIter(nums []int) *Iter {
	return &Iter{
		nums: nums,
		memo: make(map[int][2]int),
	}
}

func (it *Iter) Next(turn int) int {
	//defer log.Printf("memo: %+v", it.memo)
	if len(it.memo) < len(it.nums) {
		n := it.nums[len(it.memo)]
		it.memo[n] = [2]int{turn, 0}
		it.last = n
		return n
	}
	var num int
	seq := it.memo[it.last]
	if seq[1] == 0 {
		num = 0
	} else {
		num = seq[1] - seq[0]
	}
	if s, ok := it.memo[num]; !ok {
		it.memo[num] = [2]int{turn, 0}
	} else {
		if s[1] == 0 {
			s[1] = turn
		} else {
			s[0] = s[1]
			s[1] = turn
		}
		it.memo[num] = s
	}
	it.last = num
	return num
}

func main() {
	data := readfile("input")
	chs := strings.Split(data, ",")
	nums := make([]int, 0, len(chs))
	for _, ch := range chs {
		n, err := strconv.Atoi(ch)
		noerr(err)
		nums = append(nums, n)
	}
	it := NewIter(nums)
	for i := 1; i <= 30000000; i++ {
		n := it.Next(i)
		log.Printf("turn %d: %d", i, n)
	}
}
