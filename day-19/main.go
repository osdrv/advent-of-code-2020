package main

import (
	"fmt"
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

type Rule interface {
	Match(s string) (string, bool)
	String() string
}

type TermRule struct {
}

func (r *TermRule) Match(s string) (string, bool) {
	if len(s) == 0 {
		return "", true
	}
	return "", false
}

func (r *TermRule) String() string {
	return "[TERM]"
}

type RefRule struct {
	ix int
}

func (r *RefRule) Match(s string) (string, bool) {
	rr, ok := rules[r.ix]
	if !ok {
		log.Fatalf("unknown rule: %d", r.ix)
	}
	return rr.Match(s)
}

func (r *RefRule) String() string {
	return fmt.Sprintf("(->%d)", r.ix)
}

type AndRule struct {
	rules []Rule
}

func (r *AndRule) Match(s string) (string, bool) {
	pat := s
	var ok bool
	for _, sr := range r.rules {
		pat, ok = sr.Match(pat)
		if !ok {
			return "", false
		}
	}
	return pat, true
}

func (r *AndRule) String() string {
	chs := make([]string, 0, len(r.rules))
	for _, rr := range r.rules {
		chs = append(chs, rr.String())
	}
	return "(" + strings.Join(chs, " AND ") + ")"
}

type OrRule struct {
	rules []Rule
}

func (r *OrRule) Match(s string) (string, bool) {
	pat := s
	var ok bool
	for _, sr := range r.rules {
		pat, ok = sr.Match(s)
		if ok {
			return pat, ok
		}
	}
	return "", false
}

func (r *OrRule) String() string {
	chs := make([]string, 0, len(r.rules))
	for _, rr := range r.rules {
		chs = append(chs, rr.String())
	}
	return "(" + strings.Join(chs, " OR ") + ")"
}

type CharRule struct {
	ch rune
}

func (r *CharRule) Match(s string) (string, bool) {
	if len(s) > 0 && rune(s[0]) == r.ch {
		return s[1:], true
	}
	return "", false
}

func (r *CharRule) String() string {
	return fmt.Sprintf("\"%c\"", r.ch)
}

var (
	rules = make(map[int]Rule)
)

type Regex struct {
}

func (r *Regex) Register(ix int, rule Rule) {
	rules[ix] = rule
}

var (
	TR = &TermRule{}
)

func (r *Regex) Match(s string) bool {
	tr := &AndRule{rules: []Rule{rules[0], TR}}
	_, ok := tr.Match(s)
	return ok
}

func (r *Regex) MatchRec(d int, chain, s string, cur Rule) ([]string, bool) {
	if d > 64 {
		return nil, false
	}
	log.Printf(chain)
	switch rule := cur.(type) {
	case *TermRule:
		_, ok := rule.Match(s)
		if ok {
			return []string{""}, true
		}
		return nil, false
	case *CharRule:
		ss, ok := rule.Match(s)
		log.Printf("matching char %c against %s", rule.ch, s)
		return []string{ss}, ok
	case *RefRule:
		//if rule.ix == 8 || rule.ix == 11 {
		//	runtime.Breakpoint()
		//}
		rr := rules[rule.ix]
		return r.MatchRec(d, fmt.Sprintf("%s->%d", chain, rule.ix), s, rr)
	case *AndRule:
		opts := map[string]struct{}{s: struct{}{}}
		for _, rr := range rule.rules {
			newopts := make(map[string]struct{})
			for s1 := range opts {
				ss, ok := r.MatchRec(d, chain, s1, rr)
				if ok {
					for _, s2 := range ss {
						newopts[s2] = struct{}{}
					}
				}
			}
			if len(newopts) == 0 {
				return nil, false
			}
			opts = newopts
		}
		res := make([]string, 0, len(opts))
		for s1 := range opts {
			res = append(res, s1)
		}
		return res, true
	case *OrRule:
		opts := map[string]struct{}{}
		for _, rr := range rule.rules {
			ss, ok := r.MatchRec(d+1, chain, s, rr)
			if ok {
				for _, s1 := range ss {
					opts[s1] = struct{}{}
				}
			}
		}
		if len(opts) > 0 {
			res := make([]string, 0, len(opts))
			for s1 := range opts {
				res = append(res, s1)
			}
			return res, true
		}
		return nil, false
	default:
		log.Fatalf("unexpected rule type: %T", rule)
	}
	return nil, false
}

func parseIxs(s string) []int {
	chs := strings.Split(s, " ")
	res := make([]int, 0, len(chs))
	for _, ch := range chs {
		n, err := strconv.Atoi(ch)
		noerr(err)
		res = append(res, n)
	}
	return res
}

func parseRule(p string) (int, Rule) {
	chs := strings.SplitN(p, ": ", 2)
	ix, err := strconv.Atoi(chs[0])
	noerr(err)
	rem := chs[1]
	if rem[0] == '"' {
		ch := rune(rem[1])
		return ix, &CharRule{ch: ch}
	}
	chs = strings.Split(rem, " | ")
	ror := []Rule{}
	for _, ch := range chs {
		ixs := parseIxs(ch)
		rand := make([]Rule, 0, len(ixs))
		for _, ix := range ixs {
			rand = append(rand, &RefRule{ix: ix})
		}
		r := &AndRule{rules: rand}
		ror = append(ror, r)
	}
	if len(ror) == 1 {
		return ix, ror[0]
	}
	return ix, &OrRule{rules: ror}
}

func main() {
	reg := &Regex{}
	rulestr := strings.Split(readfile("rules"), "\n")
	for _, rs := range rulestr {
		ix, rule := parseRule(rs)
		//log.Printf("source: %s, rule: %+v", rs, rule)
		reg.Register(ix, rule)
	}

	cnt := 0
	messages := strings.Split(readfile("messages"), "\n")
	for _, m := range messages {
		if _, ok := reg.MatchRec(0, "0", m, &AndRule{rules: []Rule{rules[0], TR}}); ok {
			log.Printf("message %s matches", m)
			cnt++
		}
	}

	log.Printf("the result is: %d", cnt)
}
