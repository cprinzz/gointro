package jobsite

import (
	"math/rand"
	"time"
)

const GRID_LENGTH = 12
const GRID_WIDTH = 12

const NUM_WORKERS = 10

const LEFT = "LEFT"
const RIGHT = "RIGHT"
const UP = "UP"
const DOWN = "DOWN"

type Cone struct {
	LeftBound int
	RightBound int
	TopBound int
	BottomBound int
}

type Worker struct {
	WalkChannel chan Coords
	CurrentCoords []int
	Interval time.Duration
}

type Coords struct {
	X int
	Y int
}

// Helper function to pick a random direction and update the workers current coordinates
func pickRandomDirection(w *Worker) []int {
	randInt := rand.Intn(4)
	direction := []string{LEFT, RIGHT, UP, DOWN}[randInt]

	var newCoords []int
	if direction == LEFT{
		newCoords = []int{w.CurrentCoords[0] - 1, w.CurrentCoords[1]}
	}
	if direction == RIGHT{
		newCoords = []int{w.CurrentCoords[0] + 1, w.CurrentCoords[1]}
	}
	if direction == UP{
		newCoords = []int{w.CurrentCoords[0], w.CurrentCoords[1] - 1}
	}
	if direction == DOWN{
		newCoords = []int{w.CurrentCoords[0], w.CurrentCoords[1] + 1}
	}
	return newCoords
}

func (w *Worker) StartWalking(){
	for {
		// Worker sleeps for a random interval so that workers are concurrently walking at different rates
		time.Sleep(w.Interval)
		newCoords := pickRandomDirection(w)

		// If the new direction puts the worker on an edge, pick a different direction.
		for newCoords[0] <= 0 || newCoords[0] >= GRID_WIDTH || newCoords[1] <= 0 || newCoords[1] >= GRID_LENGTH {
			newCoords = pickRandomDirection(w)
		}
		w.CurrentCoords = newCoords

		// Send new coordinates to the workers walk channel
		w.WalkChannel <- Coords{newCoords[0], newCoords[1]}
	}
}