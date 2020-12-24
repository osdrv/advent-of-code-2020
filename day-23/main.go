package main

import "log"

func main() {
	//input := []int{3, 8, 9, 1, 2, 5, 4, 6, 7}
	input := []uint32{5, 3, 8, 9, 1, 4, 7, 6, 2}
	//n := 10
	//for len(input) < 1_000_000 {
	//	input = append(input, uint32(n))
	//	n++
	//}

	for i := 0; i < 10; i++ {
		log.Printf("========== round %d ==========", i)
		cur := input[0]
		pick := input[1:4]
		dest := -1
		max, min := uint32(0), uint32(0)
		maxpos, minpos := -1, -1

		for j := 4; j < len(input); j++ {
			if input[j] == cur-1 {
				dest = j
				break
			}
			if input[j] < cur && min < input[j] {
				min = input[j]
				minpos = j
			}
			if max < input[j] {
				max = input[j]
				maxpos = j
			}
		}
		if dest < 0 {
			if minpos > 0 {
				dest = minpos
			} else {
				dest = maxpos
			}
		}
		var pp [3]uint32
		copy(pp[:], pick)
		newin := append(input[0:1], input[4:dest+1]...)
		log.Printf("%+v", newin)
		newin = append(newin, pp[0], pp[1], pp[2])
		log.Printf("%+v", newin)
		newin = append(newin, input[dest+1:]...)
		log.Printf("%+v", newin)
		newin = append(newin[1:], newin[0])
		input = newin
		log.Printf("new input: %+v", input)

		////log.Printf("pick is: %+v, dest is: %d", pick, input[dest])
		//newin := make([]int, len(input))
		//copy(newin[0:1], input[0:1])
		////log.Printf("%+v", newin)
		//copy(newin[1:], input[4:dest+1])
		////log.Printf("%+v", newin)
		//copy(newin[dest+1+1-4:], pick)
		////log.Printf("%+v", newin)
		//copy(newin[dest+1:], input[dest+1:])
	}

	pos := -1
	for i := 0; i < len(input); i++ {
		if input[i] == 1 {
			pos = i
			break
		}
	}

	input = append(input[pos+1:], input[:pos]...)
	log.Printf("the result is: %+v", input)
	res := input[0] * input[1]
	log.Printf("the result is: %d", res)
}
