package main

import (
	day09 "de.knallisworld/aoc/aoc2019/day09/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"sync"
	"time"
)

const AocDay = 23
const AocDayName = "day23"
const AocDayTitle = "Day 23"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		puzzle := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")

		switchIn := make(chan Packet, 100)
		natIn := make(chan Packet, 10)

		switchOuts := make([]chan Packet, 50)
		for nicId := 0; nicId < 50; nicId++ {
			switchOuts[nicId] = make(chan Packet, 100)
		}

		// dispatch
		go func() {
			var lastNatSendingPacket Packet
			for packet := range switchIn {
				//fmt.Printf("Packet %d -> %d with (%d/%d)\n", packet.Source, packet.Target, packet.X, packet.Y)

				if packet.Target == 255 {
					dl.PrintSolution(fmt.Sprintf("(PACKET INSPECTION) Solution 1 is Y = %d", packet.Y))
					natIn <- packet
				} else if packet.Target == 0 {
					if packet.Y == lastNatSendingPacket.Y {
						dl.PrintSolution(fmt.Sprintf("(PACKET INSPECTION) Solution 2 is Y = %d (wait)", packet.Y))
					}
					lastNatSendingPacket = packet
				}

				if packet.Target < len(switchOuts) {
					receiver := switchOuts[packet.Target]
					receiver <- packet
				}
			}
		}()

		// NAT
		go func() {
			var lastPacket Packet
			var lastNetworkIdle bool
			for {
				select {
				case packet := <-natIn:
					lastPacket = packet
				default:
					networkIdle := true
					for i := 0; i < 50; i++ {
						if len(switchOuts[i]) > 0 {
							networkIdle = false
							break
						}
					}
					if networkIdle && lastNetworkIdle && lastPacket.Y != 0 {
						// send
						packet := Packet{Source: 255, Target: 0, X: lastPacket.X, Y: lastPacket.Y}
						switchIn <- packet
					} else if networkIdle {
						lastNetworkIdle = true
					} else if lastNetworkIdle {
						lastNetworkIdle = false
					}
				}
				time.Sleep(100 * time.Millisecond)
			}
		}()

		var wg sync.WaitGroup
		wg.Add(50)
		for nicId := 0; nicId < 50; nicId++ {
			go func(nicId int) {
				//fmt.Printf("NIC #%d starting...\n", nicId)
				bootNic(puzzle, switchOuts[nicId], switchIn, nicId)
				fmt.Printf("NIC #%d stopped\n", nicId)
				wg.Done()
			}(nicId)
		}
		wg.Wait()

	}

}

type Packet struct {
	Source int
	Target int
	X      int
	Y      int
}

func bootNic(program []int, bridgeIn <-chan Packet, bridgeOut chan<- Packet, nicId int) {

	in := make(chan int, 100)  // program stdin
	out := make(chan int, 100) // program stdout
	fin := make(chan bool)     // game end
	halt := make(chan error)   // program halt

	go func() {
		halt <- day09.ExecutionInstructions(program, in, out, false)
	}()

	in <- nicId

	go func() {
		for {
			select {
			case <-halt:
				fin <- true
				return
			case received := <-bridgeIn:
				//fmt.Printf("nic #%d receiving %d/%d\n", nicId, received.X, received.Y)
				in <- received.X
				in <- received.Y
			case t := <-out:
				x := <-out
				y := <-out
				packet := Packet{Source: nicId, Target: t, X: x, Y: y}
				//fmt.Printf("nic #%d sending %d/%d to %d\n", packet.Source, packet.X, packet.Y, packet.Target)
				bridgeOut <- packet
			default:
				in <- -1
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	<-fin
	close(in)
	close(halt)
	close(fin)
	close(bridgeOut)
}
