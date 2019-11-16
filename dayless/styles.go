package dayless

import "github.com/fatih/color"

//noinspection GoUnhandledErrorResult
func PrintDayHeader(day int, title string) {
	c := color.New(color.Bold, color.FgGreen)
	c.Println()
	c.Printf("🎄 Advent of Code 2019 - Day %02d - %s\n", day, title)
	c.Println("================================================================")
	c.Println()
}

//noinspection GoUnhandledErrorResult
func PrintStepHeader(step int) {
	c := color.New(color.Bold, color.FgGreen)
	c.Println()
	switch step {
	case 1:
		c.Println("--- Part One ---")
		break
	case 2:
		c.Println("--- Part Two ---")
		break
	default:
		c.Println("--- Part ??? ---")
	}
}

//noinspection GoUnhandledErrorResult
func PrintSolution(result interface{}) {
	c := color.New(color.Bold, color.FgGreen)
	c.Printf("🎉 The result is: %s\n", result)
	c.Println()
}

//noinspection GoUnhandledErrorResult
func PrintError(err error) {
	c := color.New(color.Bold, color.FgRed)
	c.Printf("💥 Oh no, there is a show stopper: %s\n", err.Error())
	c.Println()
}
