package dayless

import (
	"fmt"
	"strconv"
)

func ParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func BinarySearch(min int, max int, worker func(i int) bool, debug bool) int {

	cache := make(map[int]bool)
	cWorker := func(i int) bool {
		if _, exist := cache[i]; !exist {
			if debug {
				fmt.Printf("🤖 BinarySearch(invoke): Running %d\n", i)
			}
			cache[i] = worker(i)
			return cache[i]
		}
		if debug {
			fmt.Printf("🤖 BinarySearch(invoke): Running %d [CACHED]\n", i)
		}
		return cache[i]
	}

	b := 1
	for {
		if b*2 < min {
			b *= 2
		} else {
			break
		}
	}
	if debug {
		fmt.Printf("🤖 BinarySearch(init min): Lower range will be %d\n", b)
	}

	for b < max-min {
		b *= 2
		if cWorker(b) {
			if debug {
				fmt.Printf("🤖 BinarySearch(init max): Final range is [%d..%d]\n", min, b)
			}
			break
		}
	}

	for {
		if debug {
			fmt.Printf("🤖 BinarySearch(start): b=%d [min=%d, max=%d]\n", b, min, max)
		}
		for i := min; i <= max; i += b {
			if cWorker(i) {
				if b == 1 {
					if debug {
						fmt.Printf("🤖 BinarySearch(result): Found = %d\n", i)
					}
					return i
				}
				if debug {
					fmt.Printf("🤖 BinarySearch(loop): Reduce to range=[%d..%d]\n", i-b, i)
				}
				min = i - b
				max = i
				break
			}
		}

		if b == 1 {
			if debug {
				fmt.Printf("🤖 BinarySearch(result): NONE (-1)\n")
			}
			return -1
		} else {
			b /= 2
		}
	}
}

func MinInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func MaxIntArrayValue(arr []int) int {
	max := 0
	for _, n := range arr {
		max = MaxInt(max, n)
	}
	return max
}

func AbsInt(a int) int {
	if a < 0 {
		return -1 * a
	} else {
		return a
	}
}

// GGT
func GreatestCommonDivisor(a, b int) int {
	if a == 0 {
		return AbsInt(b)
	}
	if b == 0 {
		return AbsInt(a)
	}
	d := 0
	for b != 0 {
		d = a % b
		a = b
		b = d
	}
	return AbsInt(a)
}

// KGV
func LeastCommonMultiple(a, b int) int {
	return a * b / GreatestCommonDivisor(a, b)
}
