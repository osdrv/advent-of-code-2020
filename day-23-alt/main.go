package main

import (
	"log"
	"strconv"
	"strings"
)

type Node struct {
	next  *Node
	value int
}

type List struct {
	head, tail *Node
}

func NewList(arr []int) *List {
	var head, tail *Node
	for _, value := range arr {
		node := &Node{value: value}
		if tail != nil {
			tail.next = node
		} else {
			head = node
		}
		tail = node
	}
	return &List{
		head: head,
		tail: tail,
	}
}

func (l *List) Add(value int) {
	node := &Node{value: value}
	if l.tail == nil {
		l.head = node
		l.tail = node
	} else {
		l.tail.next = node
		l.tail = node
	}
}

func (l *List) SplitAt(n int) (*List, *List) {
	if n < 0 {
		return nil, l
	}
	tail := l.head
	i := 0
	for i < n && tail != nil {
		tail = tail.next
		i++
	}
	var rest *List
	if tail != nil {
		rest = &List{
			head: tail.next,
			tail: l.tail,
		}
	}
	tail.next = nil
	return &List{
		head: l.head,
		tail: tail,
	}, rest
}

func (l *List) Skip(n int) *List {
	head := l.head
	tail := l.tail
	i := 0
	for i < n && head != nil {
		head = head.next
		i++
	}
	return &List{
		head: head,
		tail: tail,
	}
}

func (l *List) Append(other *List) *List {
	if other == nil || other.head == nil {
		return l
	}
	if l.tail == nil {
		return &List{
			head: other.head,
			tail: other.tail,
		}
	}
	l.tail.next = other.head
	return &List{
		head: l.head,
		tail: other.tail,
	}
}

func (l *List) ToArray() []int {
	res := []int{}
	ptr := l.head
	for ptr != nil {
		res = append(res, ptr.value)
		ptr = ptr.next
	}
	return res
}

func (l *List) ValueAt(n int) int {
	it := NewListIterator(l)
	i := 0
	var val int
	for it.HasNext() {
		val = it.Next()
		if i == n {
			return val
		}
		i++
	}
	return -1
}

type ListIterator struct {
	list *List
	ptr  *Node
}

func NewListIterator(list *List) *ListIterator {
	return &ListIterator{
		list: list,
		ptr:  list.head,
	}
}

func (it *ListIterator) HasNext() bool {
	return it.ptr != nil
}

func (it *ListIterator) Next() int {
	value := it.ptr.value
	it.ptr = it.ptr.next
	return value
}

func main() {
	//input := []int{3, 8, 9, 1, 2, 5, 4, 6, 7}
	input := []int{5, 3, 8, 9, 1, 4, 7, 6, 2}
	n := 10
	for len(input) < 1_000_000 {
		input = append(input, n)
		n++
	}
	solve2(input)
}

func solve() {
	//input := []int{3, 8, 9, 1, 2, 5, 4, 6, 7}
	input := []int{5, 3, 8, 9, 1, 4, 7, 6, 2}
	list := NewList(input)
	n := 10
	for n <= 1_000_000 {
		list.Add(n)
		n++
	}

	var head, pick *List
	for i := 0; i < 10_000_000; i++ {
		//for i := 0; i < 100; i++ {
		if i%1000 == 0 {
			log.Printf("========== round %d ==========", i)
		}
		head, list = list.SplitAt(0)
		pick, list = list.SplitAt(2)

		//log.Printf("head: %+v, pick: %+v, list: %+v", head.ToArray(), pick.ToArray(), list.ToArray())

		dest := -1
		max, min := -1, -1
		maxpos, minpos := -1, -1

		cur := head.ValueAt(0)
		it := NewListIterator(list)
		j := 0
		for it.HasNext() {
			val := it.Next()
			if val == cur-1 {
				dest = j
				break
			}
			if val < cur && min < val {
				min = val
				minpos = j
			}
			if max < val {
				max = val
				maxpos = j
			}
			j++
		}

		if dest < 0 {
			if minpos >= 0 {
				dest = minpos
			} else {
				dest = maxpos
			}
		}

		//log.Printf("dest: %d", list.ValueAt(dest))

		newlist := head
		ch1, ch2 := list.SplitAt(dest)
		//log.Printf("ch1: %+v, ch2: %+v", ch1.ToArray(), ch2.ToArray())
		newlist = newlist.Append(ch1)
		//log.Printf("newlist: %+v", newlist.ToArray())
		newlist = newlist.Append(pick)
		//log.Printf("newlist: %+v", newlist.ToArray())
		newlist = newlist.Append(ch2)
		//log.Printf("newlist: %+v", newlist.ToArray())
		ch1, ch2 = newlist.SplitAt(0)
		//log.Printf("ch1: %+v, ch2: %+v", ch1.ToArray(), ch2.ToArray())
		newlist = ch2.Append(ch1)
		//log.Printf("newlist: %+v", newlist.ToArray())
		list = newlist
		//log.Printf("new list: %+v", list.ToArray())
	}

	pos := -1
	it := NewListIterator(list)
	i := 0
	for it.HasNext() {
		if it.Next() == 1 {
			pos = i
			break
		}
		i++
	}
	head, tail := list.SplitAt(pos)
	//log.Printf("head: %+v, tail: %+v", head.ToArray(), tail.ToArray())
	head, _ = head.SplitAt(pos - 1)
	list = tail.Append(head)

	log.Printf("%d * %d", list.head.value, list.head.next.value)

	//log.Printf("the result is: %+v", list.ToArray())
	res := uint64(list.ValueAt(0)) * uint64(list.ValueAt(1))
	log.Printf("the result is: %d", res)
}

func vecToArr(vec []int) []int {
	res := make([]int, 0, len(vec)-1)
	ptr := vec[0]
	for ptr != 0 {
		res = append(res, ptr)
		ptr = vec[ptr]
	}
	return res
}

func vecToStr(vec []int) string {
	ptr := vec[0]
	var buf strings.Builder
	for ptr != 0 {
		if buf.Len() > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString(strconv.Itoa(ptr))
		ptr = vec[ptr]
	}
	return "[" + buf.String() + "]"
}

func makeVec(nums []int) ([]int, int) {
	vec := make([]int, len(nums)+1)
	vec[0] = nums[0]
	end := -1
	for i := 0; i < len(nums)-1; i++ {
		vec[nums[i]] = nums[i+1]
	}
	vec[nums[len(nums)-1]] = 0
	end = nums[len(nums)-1]
	return vec, end
}

func solve2(nums []int) uint64 {
	vec, end := makeVec(nums)

	log.Printf("nums: %+v", nums)
	log.Printf("vec: %+v", vecToStr(vec))

	for i := 0; i < 10_000_000; i++ {
		if i%1000 == 0 {
			log.Printf("========== round %d ==========", i)
		}
		//log.Printf("===== round %d =====", i)
		cur := vec[0]
		p1, p2, p3 := vec[cur], vec[vec[cur]], vec[vec[vec[cur]]]
		//log.Printf("pick: [%d %d %d]", p1, p2, p3)
		dest := -1
		for i := cur - 1; i > 0; i-- {
			if i == p1 || i == p2 || i == p3 {
				continue
			}
			dest = i
			break
		}
		if dest < 0 {
			for i := len(vec) - 1; i > cur; i-- {
				if i == p1 || i == p2 || i == p3 {
					continue
				}
				dest = i
				break
			}
		}

		//log.Printf("destination: %d", dest)
		//log.Printf("end: %d", end)

		vec[cur] = vec[p3]
		tmp := vec[dest]
		vec[dest] = p1
		vec[p3] = tmp

		//log.Printf("vec before shift: %+v", vecToArr(vec))

		if vec[p3] == 0 {
			// end has changed
			end = p3
		}

		tmp = vec[0]
		vec[0] = vec[vec[0]]
		vec[end] = tmp
		vec[tmp] = 0
		end = tmp

		//log.Printf("vec: %+v", vec)
		//log.Printf("new list: %s", vecToStr(vec))
	}

	list := vecToArr(vec)
	//log.Printf("final list: %+v", list)

	pos := -1
	for i := 0; i < len(list); i++ {
		if list[i] == 1 {
			pos = i
			break
		}
	}

	list = append(list[pos+1:], list[:pos]...)
	//log.Printf("the result is: %+v", list)
	res := uint64(list[0]) * uint64(list[1])
	log.Printf("the result is: %d", res)

	return 0
}
