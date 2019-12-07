package day07

func ReduceMax(in <-chan int, out chan int) {
	m := 0
	for n := range in {
		if m < n {
			m = n
		}
	}
	out <- m
	close(out)
}

func FanOut(in chan int, targets ...chan int) {
	for i := range in {
		for _, c := range targets {
			c <- i
		}
	}
}
