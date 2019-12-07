package day07

func ReduceIntChannel(in <-chan int, out chan<- int, reducer func([] int) int) {
	arr := make([]int, 0)
	for n := range in {
		arr = append(arr, n)
	}
	out <- reducer(arr)
	close(out)
}

func FanOut(in <-chan int, targets ...chan<- int) {
	for i := range in {
		for _, c := range targets {
			c <- i
		}
	}
}
