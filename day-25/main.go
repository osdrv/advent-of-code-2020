package main

import "log"

const (
	MOD = 20201227
)

func findLoopNumber(pk int) int {
	subj := 7

	val := 1
	n := 0
	for {
		n++
		val = val * subj
		val %= MOD
		if val == pk {
			break
		}
	}

	return n
}

func getEncKey(subj, loop int) int {
	res := 1
	for i := 0; i < loop; i++ {
		res = res * subj
		res %= MOD
	}
	return res
}

func main() {

	input := []int{
		14788856,
		19316454,
	}

	//input := []int{
	//	5764801,
	//	17807724,
	//}

	cardLoop := findLoopNumber(input[0])
	doorLoop := findLoopNumber(input[1])

	log.Printf("card loop: %d door loop: %d", cardLoop, doorLoop)

	encKey := getEncKey(input[1], cardLoop)
	log.Printf("encryption key: %d", encKey)
}
