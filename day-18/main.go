package main

import (
	"bytes"
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

type Expr interface {
	Eval() int
	String() string
}

type Math struct {
	Op          rune
	Left, Right Expr
}

func (m *Math) Eval() int {
	switch m.Op {
	case '+':
		return m.Left.Eval() + m.Right.Eval()
	case '-':
		return m.Left.Eval() - m.Right.Eval()
	case '*':
		return m.Left.Eval() * m.Right.Eval()
	case '/':
		return m.Left.Eval() / m.Right.Eval()
	default:
		panic("wtf")
	}
}

func (m *Math) String() string {
	return fmt.Sprintf("(%s %c %s)", m.Left.String(), m.Op, m.Right.String())
}

type Number struct {
	raw string
}

func (n *Number) Eval() int {
	num, err := strconv.Atoi(n.raw)
	noerr(err)
	return num
}

func (n *Number) String() string {
	return fmt.Sprintf("%d", n.Eval())
}

type Parser struct {
	s   string
	cur int
}

func NewParser(s string) *Parser {
	return &Parser{
		s:   s,
		cur: 0,
	}
}

func (p *Parser) HasNext() bool {
	return p.cur < len(p.s)
}

func (p *Parser) Peek() rune {
	if p.HasNext() {
		return rune(p.s[p.cur])
	}
	return 0
}

func (p *Parser) CheckNumber() bool {
	r := p.Peek()
	return r >= '0' && r <= '9'
}

func (p *Parser) CheckOpenBrac() bool {
	return p.Peek() == '('
}

func (p *Parser) CheckClosBrac() bool {
	return p.Peek() == ')'
}

func (p *Parser) CheckMathOp() bool {
	switch p.Peek() {
	case '+', '-', '*', '/':
		return true
	default:
		return false
	}
}

func (p *Parser) EatWhitespace() {
	log.Printf("eating whitespace")
	for p.Peek() == ' ' {
		p.Advance()
	}
}

func (p *Parser) Advance() rune {
	if p.HasNext() {
		res := rune(p.s[p.cur])
		p.cur++
		return res
	}
	return 0
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func parseExpr(p *Parser) Expr {
	p.EatWhitespace()
	return parseMult(p)
}

func parseMult(p *Parser) Expr {
	left := parseSum(p)
	p.EatWhitespace()
	for p.Peek() == '*' {
		p.Advance()
		right := parseSum(p)
		left = &Math{
			Left:  left,
			Right: right,
			Op:    '*',
		}
		p.EatWhitespace()
	}
	return left
}

func parseSum(p *Parser) Expr {
	left := parseNumber(p)
	p.EatWhitespace()
	for p.Peek() == '+' {
		p.Advance()
		p.EatWhitespace()
		right := parseNumber(p)
		left = &Math{
			Left:  left,
			Right: right,
			Op:    '+',
		}
		p.EatWhitespace()
	}
	return left
}

func parseNumber(p *Parser) Expr {
	p.EatWhitespace()
	if p.CheckOpenBrac() {
		p.Advance()
		res := parseExpr(p)
		p.EatWhitespace()
		if p.CheckClosBrac() {
			p.Advance()
		} else {
			log.Fatalf("expected closing paren: %d", p.cur)
		}
		return res
	}
	if !p.CheckNumber() {
		log.Fatalf("expected a number at %d", p.cur)
	}
	var buf bytes.Buffer
	for p.CheckNumber() {
		buf.WriteRune(p.Peek())
		p.Advance()
	}
	return &Number{raw: buf.String()}
}

//func parseExprPart1(p *Parser) Expr {
//	var left, right Expr
//	var op rune
//	for p.HasPrev() {
//		p.EatWhitespace()
//		if p.CheckNumber() {
//			log.Printf("parsing number")
//			var b bytes.Buffer
//			for p.CheckNumber() {
//				b.WriteRune(p.Peek())
//				p.Advance()
//			}
//			right = &Number{raw: reverse(b.String())}
//			log.Printf("right is a number: %s", right)
//		} else if p.CheckClosBrac() {
//			log.Printf("parsing expr in parenteses")
//			p.Advance()
//			right = parseExprPart1(p)
//			if !p.CheckOpenBrac() {
//				log.Fatalf("no opening bracket around %d", p.cur)
//			}
//			p.Advance()
//		}
//		p.EatWhitespace()
//		if !p.HasPrev() {
//			return right
//		}
//		if p.CheckOpenBrac() {
//			return right
//		}
//		if p.CheckMathOp() {
//			op = p.Peek()
//			p.Advance()
//		} else {
//			log.Fatalf("expected a math expr around %d", p.cur)
//		}
//		p.EatWhitespace()
//		left = parseExprPart1(p)
//		right = &Math{
//			Left:  left,
//			Right: right,
//			Op:    op,
//		}
//		p.EatWhitespace()
//	}
//	return right
//}

func main() {
	lines := strings.Split(readfile("input"), "\n")
	sum := 0
	for _, line := range lines {
		//expr := parseExprPart1(NewParser(line))
		//runtime.Breakpoint()
		expr := parseExpr(NewParser(line))
		log.Printf("expr: %s", expr)
		res := expr.Eval()
		log.Printf("expr: %s = %d", line, res)
		sum += res
	}
	log.Printf("sum is: %d", sum)
}
