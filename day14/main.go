package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"strings"
	"time"
)

const AocDay = 14
const AocDayName = "day14"
const AocDayTitle = "Day 14"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		puzzle, _ := dl.ReadFileToArray(AocDayName + "/puzzle1.txt")
		reactions := newReactions(puzzle)
		minRequiredAmountOfOre := react(reactions, Material{1, "FUEL"}, "ORE")
		dl.PrintSolution(fmt.Sprintf("Mininum mount of 'ORE' is %d", minRequiredAmountOfOre.Amount))
	}

	dl.PrintStepHeader(2)
	{
		puzzle, _ := dl.ReadFileToArray(AocDayName + "/puzzle1.txt")
		reactions := newReactions(puzzle)
		final := dl.BinarySearch(2, 10000000, func(i int) bool {
			result := react(reactions, Material{i, "FUEL"}, "ORE")
			return result != nil && result.Amount > 1_000_000_000_000
		}, true)
		dl.PrintSolution(fmt.Sprintf("Maximum mount of 'FUEL' is %d", final-1))
	}

}

func findReaction(reactions []Reaction, outputChemical string) Reaction {
	//fmt.Printf("Looking for %s\n", outputChemical)
	for _, reaction := range reactions {
		if reaction.output.Chemical == outputChemical {
			return reaction
		}
	}
	panic("invalid reaction output required: " + outputChemical)
}

func react(reactions []Reaction, targetMaterial Material, baseChemical string) *Material {

	type Amount struct {
		required, supplied int
	}

	// build map for lookup performance
	reactionMap := make(map[string]Reaction)
	for _, r := range reactions {
		reactionMap[r.output.Chemical] = r
	}

	pot := make(map[string]*Amount)
	pot[targetMaterial.Chemical] = &Amount{required: targetMaterial.Amount, supplied: 0}

	// flatten (reacting)
	for {
		// finished if only ORE is required
		finished := true
		for m, i := range pot {
			if m != "ORE" && i.required > i.supplied {
				finished = false
			}
		}
		if finished {
			break
		}

		for inputChemical := range pot {
			requiredMaterial := pot[inputChemical]
			if inputChemical == "ORE" {
				continue
			} else if requiredMaterial.required <= requiredMaterial.supplied {
				continue
			}
			delta := requiredMaterial.required - requiredMaterial.supplied
			//reaction := findReaction(reactions, inputChemical)
			reaction := reactionMap[inputChemical]
			mul := 1
			for reaction.output.Amount*mul < delta {
				mul++
			}
			requiredMaterial.supplied += mul * reaction.output.Amount
			for _, input := range reaction.inputs {
				if _, exist := pot[input.Chemical]; !exist {
					pot[input.Chemical] = &Amount{0, 0,}
				}
				required := pot[input.Chemical]
				required.required += mul * input.Amount
			}
		}
	}

	if m, exist := pot[baseChemical]; exist {
		return &Material{Amount: m.required, Chemical: baseChemical}
	} else {
		return nil
	}
}

type Material struct {
	Amount   int
	Chemical string
}

type Reaction struct {
	inputs []Material
	output Material
}

func newReaction(line string) Reaction {
	ioParts := strings.Split(line, " => ")
	inputs := make([]Material, 0)
	{
		for _, part := range strings.Split(ioParts[0], ", ") {
			w := strings.Index(part, " ")
			inputs = append(inputs, Material{Amount: dl.ParseInt(part[0:w]), Chemical: strings.TrimSpace(part[w+1:])})
		}
	}
	var output Material
	{
		w := strings.Index(ioParts[1], " ")
		output = Material{Amount: dl.ParseInt(ioParts[1][0:w]), Chemical: strings.TrimSpace(ioParts[1][w+1:])}
	}
	return Reaction{inputs: inputs, output: output}
}

func newReactions(lines []string) []Reaction {
	result := make([]Reaction, len(lines))
	for i, line := range lines {
		result[i] = newReaction(line)
	}
	return result
}
