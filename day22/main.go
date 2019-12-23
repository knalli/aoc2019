package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"
)

const AocDay = 22
const AocDayName = "day22"
const AocDayTitle = "Day 22"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		puzzle, _ := dl.ReadFileToArray(AocDayName + "/puzzle1.txt")
		const repeat = 1
		const length = 10007
		const card = 2019
		deck := shuffle(length, puzzle, repeat)
		position := 0
		for i, v := range deck {
			if v == 2019 {
				position = i
				break
			}
		}
		dl.PrintSolution(fmt.Sprintf("The position of card '%d' is #%d", card, position))

		// reworked using simple functions (prep for part2)
		dl.PrintSolution(fmt.Sprintf("The position of card '%d' (tracked) is #%d", card, trackCard(length, puzzle, repeat, card)))

		// first try for part2, left for the records
		dl.PrintSolution(fmt.Sprintf("Backtracing the final position #%d, card is '%d'", position, traceCardByPosition(length, puzzle, 1, int64(position))))

		{
			// Additional linear composite function: f(x)=a*x+b
			a, b := resolveFunctionParams(length, puzzle)
			f := func(i int64) int64 {
				return (a*i + b) % length
			}
			dl.PrintSolution(fmt.Sprintf("Given length = %d, f(x) = %d*x + %d ; so: f(%d) = %d", length, a, b, card, f(card)))
		}
	}

	{
		dl.PrintStepHeader(2)
		puzzle, _ := dl.ReadFileToArray(AocDayName + "/puzzle1.txt")
		const repeat = 101741582076661
		const length = 119315717514047
		const card = 2020

		dl.PrintSolution(fmt.Sprintf("%d", solve2(puzzle, length, repeat, card)))

		{
			a, b := resolveInverseFunctionParams(length, puzzle)

			r1 := big.NewInt(card)
			r1.Mul(r1, modpow2(a, big.NewInt(repeat), big.NewInt(length)))

			r2 := modpow2(a, big.NewInt(repeat), big.NewInt(length))
			r2.Add(r2, big.NewInt(length))
			r2.Sub(r2, big.NewInt(1))
			r2.Mul(r2, b)
			a1 := big.NewInt(-1)
			a1.Add(a1, a)
			r2.Mul(r2, modpow2(a1, big.NewInt(length-2), big.NewInt(length)))

			r1.Add(r1, r2)
			r1.Mod(r1, big.NewInt(length))
			// x := (card*modpow(a, repeat, length) + (modpow(a, repeat, length)+length-1)*b*modinv(a-1, length)) % length
			x := r1.Int64()
			dl.PrintSolution(fmt.Sprintf("%d", x))
		}
	}

}

func modpow(b int64, e int64, m int64) int64 {
	/*
		if e == 0 {
			return 1
		} else if e%2 == 0 {
			return modpow((b*b)%m, e/2, m)
		} else {
			return (b * modpow(b, e-1, m)) % m
		}
	*/
	r := big.NewInt(b)
	r.Exp(r, big.NewInt(e), big.NewInt(m))
	return r.Int64()
}

func modpow2(b *big.Int, e *big.Int, m *big.Int) *big.Int {
	/*
		if e == 0 {
			return 1
		} else if e%2 == 0 {
			return modpow((b*b)%m, e/2, m)
		} else {
			return (b * modpow(b, e-1, m)) % m
		}
	*/
	r := big.NewInt(0)
	r.Add(r, b)
	r.Exp(r, e, m)
	return r
}

// https://www.reddit.com/r/adventofcode/comments/ee0rqi/2019_day_22_solutions/fbqs5bk/
func solve(lines []string, c, n, p int64) int64 {
	o := int64(0)
	i := int64(1)
	inv := func(x int64) int64 {
		return modpow(x, c-2, c)
	}
	for _, line := range lines {
		s := strings.Split(line, " ")
		if s[0] == "cut" {
			o += i * int64(dl.ParseInt(s[len(s)-1]))
		} else if s[1] == "with" {
			i *= inv(int64(dl.ParseInt(s[len(s)-1])))
		} else if s[1] == "into" {
			o -= i
			i *= -1
		}
	}
	o *= inv(1 - i)
	i = modpow(i, n, c)
	return (p*i + (1-i)*o) % c
}

// like #solve() but using math.big.Int
func solve2(lines []string, c, n, p int64) int64 {
	c2 := big.NewInt(c)
	c2m2 := big.NewInt(c)
	c2m2.Sub(c2m2, big.NewInt(2))
	o := big.NewInt(int64(0))
	i := big.NewInt(int64(1))
	inv := func(x *big.Int) *big.Int {
		return modpow2(x, c2m2, c2)
	}
	for _, line := range lines {
		s := strings.Split(line, " ")
		if s[0] == "cut" {
			arg := big.NewInt(int64(dl.ParseInt(s[len(s)-1])))
			r := big.NewInt(0)
			r.Add(r, i)
			r.Mul(r, arg)
			o.Add(o, r)
		} else if s[1] == "with" {
			arg := big.NewInt(int64(dl.ParseInt(s[len(s)-1])))
			i.Mul(i, inv(arg))
		} else if s[1] == "into" {
			o.Sub(o, i)
			i.Neg(i)
		}
	}
	{
		r := big.NewInt(1)
		r.Sub(r, i)
		o.Mul(o, inv(r))
	}
	i = modpow2(i, big.NewInt(n), c2)

	{
		r := big.NewInt(0)
		r.Add(r, big.NewInt(p))
		r.Mul(r, i)

		r2 := big.NewInt(1)
		r2.Sub(r2, i)
		r2.Mul(r2, o)

		r.Add(r, r2)
		r.Mod(r, c2)

		return r.Int64()
	}
}

// modpow(n, m-2, m)
func modinv(n int64, m int64) int64 {
	return modpow(n, m-2, m)
	//x := big.NewInt(n)
	//x.ModInverse(x, big.NewInt(m))
	//return x.Int64()
}

func dealIntoStack(deck map[int]int) map[int]int {
	stack := make(map[int]int)
	l := len(deck)
	for i := 0; i < l; i++ {
		stack[l-i-1] = deck[i]
	}
	return stack
}

func cutCards(deck map[int]int, n int) map[int]int {
	l := len(deck)
	if n < 0 {
		n += l
	}
	stack := make(map[int]int)
	for i := n; i < l; i++ {
		stack[i-n] = deck[i]
	}
	for i := 0; i < n; i++ {
		stack[l-n+i] = deck[i]
	}
	return stack
}

func dealWithInc(deck map[int]int, n int) map[int]int {
	l := len(deck)
	stack := make(map[int]int)
	/*
		for i:=0;i<len(deck);i++ {
			stack[i] = -1
		}
	*/
	pos := 0
	for i := 0; i < l; i++ {
		/*
			if stack[pos] != -1 {
				panic("invalid deal with increment operation")
			}
		*/
		stack[pos] = deck[i]
		pos = (pos + n) % l
	}
	return stack
}

type Technique func(deck map[int]int) map[int]int

func parseTechniques(lines []string) []Technique {
	result := make([]Technique, 0)
	for _, line := range lines {
		if line[0:3] == "cut" {
			n := dl.ParseInt(line[4:])
			result = append(result, func(deck map[int]int) map[int]int {
				return cutCards(deck, n)
			})
		} else if line[0:6] == "deal w" {
			n := dl.ParseInt(line[len("deal with increment "):])
			result = append(result, func(deck map[int]int) map[int]int {
				return dealWithInc(deck, n)
			})
		} else if line[0:6] == "deal i" {
			result = append(result, func(deck map[int]int) map[int]int {
				return dealIntoStack(deck)
			})
		} else {
			panic("unknown technique")
		}
	}
	return result
}

func shuffle(factorySize int, techniques []string, repeats int) map[int]int {

	factory := make(map[int]int, factorySize)
	for i := 0; i < factorySize; i++ {
		factory[i] = i
	}
	fmt.Println("Factory deck created")

	deck := make(map[int]int, factorySize)
	for k, v := range factory {
		deck[k] = v
	}
	// copy(deck, factory)

	techniqueFns := parseTechniques(techniques)

	for r := 0; r < repeats; r++ {
		for _, fn := range techniqueFns {
			deck = fn(deck)
		}
		if reflect.DeepEqual(deck, factory) {
			left := repeats - r
			d := r + 1
			jump := (left / d) * d
			if jump > 0 {
				fmt.Printf("Found cycle at %d (cached %d, length %d)\n", r, 0, d)
				r += jump // auto floor in div
			}
		}
	}

	return deck
}

func trackCard(length int64, techniques []string, repeats int64, startPosition int64) int64 {

	factory := startPosition
	card := factory

	for r := int64(0); r < repeats; r++ {
		for _, line := range techniques {
			if line[0:3] == "cut" {
				n := int64(dl.ParseInt(line[len("cut "):]))
				card = (length + card - n) % length
			} else if line[0:6] == "deal w" {
				n := int64(dl.ParseInt(line[len("deal with increment "):]))
				card = (card * n) % length
			} else if line[0:6] == "deal i" {
				card = ((length - 1) - card) % length
			} else {
				panic("unknown technique")
			}
		}
		if card == factory {
			left := repeats - r
			d := r + 1
			jump := (left / d) * d
			if jump > 0 {
				fmt.Printf("Found cycle at %d (cached %d, length %d)\n", r, 0, d)
				r += jump // auto floor in div
			}
		}
	}

	return card
}

func resolveFunctionParams(length int64, techniques []string) (int64, int64) {

	a := int64(1)
	b := int64(0)

	for _, line := range techniques {
		if line[0:3] == "cut" {
			n := int64(dl.ParseInt(line[len("cut "):]))
			if n < 0 {
				n = length + n
			}
			b = b - n
		} else if line[0:6] == "deal w" {
			n := int64(dl.ParseInt(line[len("deal with increment "):]))
			a, b = (a*n)%length, (b*n)%length
		} else if line[0:6] == "deal i" {
			a, b = -a, length-1-b
		} else {
			panic("unknown technique")
		}
	}

	return a % length, b % length
}

func traceCardByPosition(length int64, techniques []string, repeats int64, searchFinalPosition int64) int64 {

	factory := searchFinalPosition

	card := big.NewInt(factory)
	bigLength := big.NewInt(length)

	for r := int64(0); r < repeats; r++ {
		for i := len(techniques) - 1; i >= 0; i-- {
			line := techniques[i]
			if line[0:3] == "cut" {
				n := int64(dl.ParseInt(line[len("cut "):]))
				//card = (length + card + n) % length
				card.Add(card, bigLength)
				card.Add(card, big.NewInt(n))
				card.Mod(card, bigLength)
			} else if line[0:6] == "deal w" {
				n := int64(dl.ParseInt(line[len("deal with increment "):]))
				// FIXME: not correct: card = card/n + (card%n)*n
				x := big.NewInt(n)
				x.ModInverse(x, bigLength)
				x.Mul(x, card)
				x.Mod(x, bigLength)
				card = x
			} else if line[0:6] == "deal i" {
				//card = ((length - 1) - card) % length
				card.Neg(card)
				card.Add(card, bigLength)
				card.Sub(card, big.NewInt(1))
				card.Mod(card, bigLength)
			} else {
				panic("unknown technique")
			}
		}
		if card.Int64() == factory {
			left := repeats - r
			d := r + 1
			jump := (left / d) * d
			if jump > 0 {
				fmt.Printf("Found cycle at %d (cached %d, length %d)\n", r, 0, d)
				r += jump // auto floor in div
			}
		}
	}

	return card.Int64()
}

func resolveInverseFunctionParams(length int64, techniques []string) (*big.Int, *big.Int) {

	a := big.NewInt(1)
	b := big.NewInt(0)

	for i := len(techniques) - 1; i >= 0; i-- {
		line := techniques[i]
		if line[0:3] == "cut" {
			n := big.NewInt(int64(dl.ParseInt(line[len("cut "):])))
			b.Add(b, n)
			// a, b = a, b+n
		} else if line[0:6] == "deal w" {
			n := big.NewInt(int64(dl.ParseInt(line[len("deal with increment "):])))
			ninv := modpow2(n, big.NewInt(length-2), big.NewInt(length))
			a.Mul(a, ninv)
			b.Mul(b, ninv)
			// a, b = ninv*a, ninv*b
		} else if line[0:6] == "deal i" {
			//card = ((length - 1) - card) % length
			// a, b = -a, length-1-b
			a.Neg(a)
			b.Neg(b)
			b.Sub(b, big.NewInt(1))
			b.Add(b, big.NewInt(length))
		} else {
			panic("unknown technique")
		}
	}

	a.Mod(a, big.NewInt(length))
	b.Mod(b, big.NewInt(length))

	return a, b
}
