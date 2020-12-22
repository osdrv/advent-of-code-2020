package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
)

var (
	re = regexp.MustCompile("^(\\d+)\\-(\\d+)\\s(\\w{1}):\\s(.+)$")
)

func main() {
	file, err := os.Open("input")
	if err != nil {
		log.Fatalf("failed to open input file: %s", err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)

	validCnt := 0

	for sc.Scan() {
		line := sc.Text()
		match := re.FindStringSubmatch(line)
		lows, highs, chars, pwd := match[1], match[2], match[3], match[4]
		low, err := strconv.Atoi(lows)
		if err != nil {
			log.Fatalf("failed to parse number %q: %s", lows, err)
		}
		high, err := strconv.Atoi(highs)
		if err != nil {
			log.Fatalf("failed to parse number %q: %s", lows, err)
		}
		//if validateCharNum(pwd, rune(chars[0]), low, high) {
		//	validCnt++
		//}
		if validateNewPolicy(pwd, rune(chars[0]), low, high) {
			validCnt++
		}
		log.Printf("%+v", match)
	}

	if err := sc.Err(); err != nil {
		log.Fatalf("failed to read line: %s", err)
	}

	log.Printf("the number of valid passwords: %d", validCnt)
}

func validateCharNum(s string, r rune, low, high int) bool {
	cnt := 0
	for _, ch := range s {
		if ch == r {
			cnt++
		}
	}
	return cnt >= low && cnt <= high
}

func validateNewPolicy(s string, r rune, mustIx, MustNotIx int) bool {
	return (rune(s[mustIx-1]) == r && rune(s[MustNotIx-1]) != r) ||
		rune(s[mustIx-1]) != r && rune(s[MustNotIx-1]) == r
}
