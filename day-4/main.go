package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	data, err := ioutil.ReadFile("input")
	if err != nil {
		log.Fatalf("failed to read from input file: %s", err)
	}
	chunks := strings.Split(strings.TrimRight(string(data), "\n"), "\n\n")

	valid := 0
	for _, chunk := range chunks {
		passport := parsePassport(chunk)
		if isValid(passport) {
			valid++
		}
	}
	log.Printf("number of valid passports: %d", valid)
}

func parsePassport(s string) map[string]string {
	p := make(map[string]string)
	for _, line := range strings.Split(s, "\n") {
		for _, kv := range strings.Split(line, " ") {
			kvarr := strings.Split(kv, ":")
			k, v := kvarr[0], kvarr[1]
			p[k] = v
		}
	}
	return p
}

var (
	VALIDATORS = map[string]func(string) bool{
		"byr": func(v string) bool {
			if ok, err := regexp.MatchString("^\\d{4}$", v); err != nil {
				log.Fatalf("failed to parse regex: %s", err)
			} else if !ok {
				return false
			}
			i, _ := strconv.Atoi(v)
			return i >= 1920 && i <= 2002
		},
		"iyr": func(v string) bool {
			if ok, err := regexp.MatchString("^\\d{4}$", v); err != nil {
				log.Fatalf("failed to parse regex: %s", err)
			} else if !ok {
				return false
			}
			i, _ := strconv.Atoi(v)
			return i >= 2010 && i <= 2020
		},
		"eyr": func(v string) bool {
			if ok, err := regexp.MatchString("^\\d{4}$", v); err != nil {
				log.Fatalf("failed to parse regex: %s", err)
			} else if !ok {
				return false
			}
			i, _ := strconv.Atoi(v)
			return i >= 2020 && i <= 2030
		},
		"hgt": func(v string) bool {
			re := regexp.MustCompile("^(\\d{2,3})(cm|in)$")
			m := re.FindStringSubmatch(v)
			if len(m) != 3 {
				return false
			}
			i, _ := strconv.Atoi(m[1])
			switch m[2] {
			case "cm":
				return i >= 150 && i <= 193
			case "in":
				return i >= 59 && i <= 76
			}
			return false
		},
		"hcl": func(v string) bool {
			ok, err := regexp.MatchString("^#[0-9a-f]{6}$", v)
			if err != nil {
				log.Fatalf("failed to compile regex: %s", err)
			}
			return ok
		},
		"ecl": func(v string) bool {
			ok, err := regexp.MatchString("^(amb|blu|brn|gry|grn|hzl|oth)$", v)
			if err != nil {
				log.Fatalf("failed to comile ecl regexp: %s", err)
			}
			return ok
		},
		"pid": func(v string) bool {
			ok, err := regexp.MatchString("^\\d{9}$", v)
			if err != nil {
				log.Fatalf("failed to compile pid regexp: %s", err)
			}
			return ok
		},
	}
)

func isValid(p map[string]string) bool {
	log.Printf("-----")
	log.Printf("passport: %+v", p)
	for k, vldtr := range VALIDATORS {
		v, ok := p[k]
		if !ok {
			log.Printf("missing field: %s", k)
			return false
		}
		if !vldtr(v) {
			log.Printf("invalid attribute: %s: %q", k, v)
			return false
		}
	}
	return true
}
