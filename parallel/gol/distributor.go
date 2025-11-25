package gol

import (
	"fmt"
	"sync"
	"time"

	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
	keyPresses <-chan rune
}

func worker(p Params, world [][]byte, startY, height int, turn int, events chan<- Event) [][]byte {
	nextSlice := make([][]byte, height)
	for i := range nextSlice {
		nextSlice[i] = make([]byte, p.ImageWidth)
	}

	flipped := []util.Cell{}

	for y := 0; y < height; y++ {
		globalY := startY + y
		for x := 0; x < p.ImageWidth; x++ {
			sum := 0
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					nx := (x + dx + p.ImageWidth) % p.ImageWidth
					ny := (globalY + dy + p.ImageHeight) % p.ImageHeight
					if world[ny][nx] == 0xFF {
						sum++
					}
				}
			}

			current := world[globalY][x]
			next := byte(0x00)
			if current == 0xFF {
				if sum == 2 || sum == 3 {
					next = 0xFF
				}
			} else {
				if sum == 3 {
					next = 0xFF
				}
			}
			nextSlice[y][x] = next

			if current != next {
				flipped = append(flipped, util.Cell{X: x, Y: globalY})
			}
		}
	}

	if len(flipped) > 0 {
		events <- CellsFlipped{CompletedTurns: turn + 1, Cells: flipped}
	}

	return nextSlice
}

func savePGM(p Params, c distributorChannels, world [][]byte, turn int) {
	c.ioCommand <- ioOutput
	filename := fmt.Sprintf("%dx%dx%d", p.ImageWidth, p.ImageHeight, turn)
	c.ioFilename <- filename

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			c.ioOutput <- world[y][x]
		}
	}

	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- ImageOutputComplete{CompletedTurns: turn, Filename: filename}
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	// TODO: Create a 2D slice to store the world.
	c.ioCommand <- ioInput
	filename := fmt.Sprintf("%dx%d", p.ImageWidth, p.ImageHeight)
	c.ioFilename <- filename

	world := make([][]byte, p.ImageHeight)
	for i := range world {
		world[i] = make([]byte, p.ImageWidth)
	}

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			world[y][x] = <-c.ioInput
			if world[y][x] == 0xFF {
				c.events <- CellFlipped{CompletedTurns: 0, Cell: util.Cell{X: x, Y: y}}
			}
		}
	}

	turn := 0
	c.events <- StateChange{turn, Executing} // Initial state
	paused := false
	var mu sync.Mutex

	ticker := time.NewTicker(2 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				mu.Lock()
				count := 0
				for y := 0; y < p.ImageHeight; y++ {
					for x := 0; x < p.ImageWidth; x++ {
						if world[y][x] == 0xFF {
							count++
						}
					}
				}
				currentTurn := turn
				mu.Unlock()
				c.events <- AliveCellsCount{CompletedTurns: currentTurn, CellsCount: count}
			}
		}
	}()

	sliceHeight := p.ImageHeight / p.Threads
	results := make([]chan [][]byte, p.Threads)
	for i := 0; i < p.Threads; i++ {
		results[i] = make(chan [][]byte)
	}

	// TODO: Execute all turns of the Game of Life.
	for turn < p.Turns {
		select {
		case key := <-c.keyPresses:
			switch key {
			case 's':
				mu.Lock()
				worldCopy := make([][]byte, p.ImageHeight)
				for y := range world {
					worldCopy[y] = make([]byte, p.ImageWidth)
					copy(worldCopy[y], world[y])
				}
				currentTurn := turn
				mu.Unlock()
				savePGM(p, c, worldCopy, currentTurn)

			case 'q':
				ticker.Stop()
				done <- true

				mu.Lock()
				// TODO: Report final state and save for PGM
				alive := []util.Cell{}
				for y := 0; y < p.ImageHeight; y++ {
					for x := 0; x < p.ImageWidth; x++ {
						if world[y][x] == 0xFF {
							alive = append(alive, util.Cell{X: x, Y: y})
						}
					}
				}
				c.events <- FinalTurnComplete{CompletedTurns: turn, Alive: alive}

				worldCopy := make([][]byte, p.ImageHeight)
				for y := range world {
					worldCopy[y] = make([]byte, p.ImageWidth)
					copy(worldCopy[y], world[y])
				}
				currentTurn := turn
				mu.Unlock()

				savePGM(p, c, worldCopy, currentTurn)

				c.events <- StateChange{currentTurn, Quitting}
				close(c.events)
				return

			case 'p':
				paused = !paused
				if paused {
					mu.Lock()
					currentTurn := turn
					mu.Unlock()
					c.events <- StateChange{currentTurn, Paused}
				} else {
					mu.Lock()
					currentTurn := turn
					mu.Unlock()
					c.events <- StateChange{currentTurn, Executing}
				}
			}

		default:
			if !paused {
				mu.Lock()
				worldCopyForWorkers := make([][]byte, p.ImageHeight)
				for y := range world {
					worldCopyForWorkers[y] = make([]byte, p.ImageWidth)
					copy(worldCopyForWorkers[y], world[y])
				}
				currentTurn := turn
				mu.Unlock()

				nextWorld := make([][]byte, p.ImageHeight)
				for i := range nextWorld {
					nextWorld[i] = make([]byte, p.ImageWidth)
				}

				for w := 0; w < p.Threads; w++ {
					startY := w * sliceHeight
					height := sliceHeight
					if w == p.Threads-1 {
						height += p.ImageHeight % p.Threads
					}

					go func(w int, startY, height int) {
						nextSlice := worker(p, worldCopyForWorkers, startY, height, currentTurn, c.events)
						results[w] <- nextSlice
					}(w, startY, height)
				}

				for w := 0; w < p.Threads; w++ {
					nextSlice := <-results[w]
					startY := w * sliceHeight
					for y := 0; y < len(nextSlice); y++ {
						copy(nextWorld[startY+y], nextSlice[y])
					}
				}

				mu.Lock()
				world = nextWorld
				turn++
				mu.Unlock()

				c.events <- TurnComplete{CompletedTurns: currentTurn + 1}
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	ticker.Stop()
	done <- true

	mu.Lock()
	alive := []util.Cell{}
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == 0xFF {
				alive = append(alive, util.Cell{X: x, Y: y})
			}
		}
	}
	finalTurn := turn

	worldCopy := make([][]byte, p.ImageHeight)
	for y := range world {
		worldCopy[y] = make([]byte, p.ImageWidth)
		copy(worldCopy[y], world[y])
	}
	mu.Unlock()

	c.events <- FinalTurnComplete{CompletedTurns: finalTurn, Alive: alive}
	savePGM(p, c, worldCopy, finalTurn)

	c.events <- StateChange{finalTurn, Quitting}
	close(c.events)
}
