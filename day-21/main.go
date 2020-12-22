package main

import (
	"io/ioutil"
	"log"
	"sort"
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
	foods := make(map[string]int)
	alrg := make(map[string]map[string]struct{})

	for _, line := range lines {
		line = line[:len(line)-1]
		chs := strings.SplitN(line, " (contains ", 2)
		fs, as := chs[0], chs[1]
		log.Printf("fs: %q, as: %q", fs, as)
		ff := []string{}
		for _, f := range strings.Split(fs, " ") {
			ff = append(ff, f)
			foods[f]++
		}
		aa := []string{}
		for _, a := range strings.Split(as, ", ") {
			aa = append(aa, a)
			if _, ok := alrg[a]; !ok {
				alrg[a] = make(map[string]struct{})
				for _, f := range ff {
					alrg[a][f] = struct{}{}
				}
			} else {
				nalrg := make(map[string]struct{})
				for _, f := range ff {
					if _, ok := alrg[a][f]; ok {
						nalrg[f] = struct{}{}
					}
				}
				alrg[a] = nalrg
			}
		}
	}

	log.Printf("%+v", foods)
	log.Printf("%+v", alrg)

	inert := make(map[string]struct{})
	cnt := 0
Food:
	for f, n := range foods {
		for _, a := range alrg {
			if _, ok := a[f]; ok {
				continue Food
			}
		}
		log.Printf("%s is not found anywhere", f)
		cnt += n
		inert[f] = struct{}{}
	}
	log.Printf("cnt is: %d", cnt)

	aa := []string{}
	for a, ff := range alrg {
		aa = append(aa, a)
		remove := []string{}
		for f := range ff {
			if _, ok := inert[f]; ok {
				remove = append(remove, f)
			}
		}
		for _, r := range remove {
			delete(alrg[a], r)
		}
	}

	var recurse func([]string, map[string]string) (map[string]string, bool)
	recurse = func(aa []string, dang map[string]string) (map[string]string, bool) {
		if len(aa) == 0 {
			return dang, true
		}
		var a string
		a, aa = aa[0], aa[1:]
		for f := range alrg[a] {
			if _, ok := dang[f]; ok {
				continue
			}
			dang[f] = a
			if dd, ok := recurse(aa, dang); ok {
				return dd, true
			}
			delete(dang, f) // did not work
		}
		return nil, false
	}

	dang, _ := recurse(aa, make(map[string]string))
	log.Printf("refined alergens: %+v", dang)

	res := make([]string, 0, len(dang))
	for f := range dang {
		res = append(res, f)
	}
	sort.Slice(res, func(i, j int) bool {
		return dang[res[i]] < dang[res[j]]
	})

	log.Print(strings.Join(res, ","))

}
