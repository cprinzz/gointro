package jobsite

import (
	"math/rand"
	"time"
)

// Define consts for grid length and width

// Define const for num workers

// Define consts for directions

// Create Cone struct with L,R,T,B

//Create Worker struct with WalkChannel, CurrentCoords, Interval

// Private helper function to pick a random direction and update the workers current coordinates
func pickRandomDirection(w *Worker) []int {

}

// Create StartWalking function that will be a routine sending new coordinates at random intervals
func (w *Worker) StartWalking() {
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

	}
}
