package dayless

func Permutations(arr []int) [][]int {
	result := make([][]int, 0)
	for p := make([]int, len(arr)); p[0] < len(p); nextPerm(p) {
		result = append(result, getPerm(arr, p))
	}
	return result
}

func nextPerm(p []int) {
	for i := len(p) - 1; i >= 0; i-- {
		if i == 0 || p[i] < len(p)-i-1 {
			p[i]++
			return
		}
		p[i] = 0
	}
}

func getPerm(orig, p []int) []int {
	result := append([]int{}, orig...)
	for i, v := range p {
		result[i], result[i+v] = result[i+v], result[i]
	}
	return result
}
