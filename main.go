package main

import (
	"fmt"
	"gointro/jobsite"
	"gointro/render"
	"math/rand"
	"time"
)

// This program simulates a Crane Safety solution
// Each worker has a "tag" that sends coordinates through a channel to a function (Render) that maps them onto a grid
// and determines how many workers are in the danger zone. This function sends the number of workers in danger through
// an alert channel. This value is received by an anonymous GoRoutine and printed to the screen.

func main() {
	// Create an array of channels for each worker to report it's coordinates. This is analogous to the RedPoint badges.
	var walkChannels []chan jobsite.Coords
	// Create an array of workers. Each worker has a walkChannel to notify channel receivers about new position.
	var workers []jobsite.Worker
	// Done channel is created so that program doesn't immediately exit starting last GoRoutine
	done := make(chan bool)

	for i := 0; i < jobsite.NUM_WORKERS; i++ {
		// Create a walk chan for each worker and append it to the walkChannels array
		walkChan := make(chan jobsite.Coords)
		walkChannels = append(walkChannels, walkChan)

		// Give each worker a random initial coordinate
		randX := rand.Intn(jobsite.GRID_WIDTH)
		randY := rand.Intn(jobsite.GRID_LENGTH)

		// Create a new worker. Interval defines how long the worker waits between each walk.
		worker:= jobsite.Worker{
			CurrentCoords: []int{randX,randY},
			WalkChannel:walkChan,
			Interval: time.Duration(rand.Intn(500) + 100) * time.Millisecond,
		}

		// Append worker to worker array
		workers = append(workers, worker)

		// Start a goroutine for the worker to start walking. Loop will start NUM_WORKERS number of GoRoutines.
		go worker.StartWalking()
	}

	// Create channel for Render to pass number of workers currently in cone.
	alertChan := make(chan int)
	// Create channel for Render to pass true if there are 0 in the cone
	allGoodChan := make(chan bool)

	// Start Render goroutine
	go render.Render(walkChannels, jobsite.Cone{4,8,4,8}, alertChan, allGoodChan)

	// Start anonymous routine to print numWorkers whenever a val is passed through alertChan
	go func() {
		for {
			select {
				case numWorkers := <- alertChan:
					fmt.Printf("\033[3;1H")
					fmt.Printf("           %d workers in the cone    \n", numWorkers)
				case <- allGoodChan:
					fmt.Printf("\033[3;1H")
					fmt.Printf("           Alright, alright, alright\n")
			}
		}

	}()
	// GoRoutines continue execution until done (which is never)
	// Quit using ctrl-C
	<- done
}
