package render

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gointro/jobsite"
	"html"
	"os"
	"reflect"
	"strconv"
	"sync"
)

type Row struct {
	RowNum int
	Row    [jobsite.GRID_WIDTH]string
}

// Receives new coordinate values from each worker, maps them onto a grid, prints the grid
func Render(walkChans []chan jobsite.Coords, coneBounds jobsite.Cone, alertChan chan int, allGoodChan chan bool) {
	// Define emojis for worker and warn
	worker := html.UnescapeString("&#" + strconv.Itoa(128513))
	warn := html.UnescapeString("&#" + strconv.Itoa(128520))

	// Create array of coordinates to map onto grid
	var workerLocs [jobsite.NUM_WORKERS]jobsite.Coords

	// Create array of cases for each worker chan
	cases := make([]reflect.SelectCase, len(walkChans))
	for i, ch := range walkChans {
		// Creates a case for each channel
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}

	for {
		// Whenever a value is sent to a walkChannel by a worker, it is received as ch (index of channel) and the value sent
		ch, value, _ := reflect.Select(cases)
		// Cast value to Coords
		workerLocs[ch] = value.Interface().(jobsite.Coords)

		var grid [jobsite.GRID_LENGTH][jobsite.GRID_WIDTH]string

		numWorkersInCone := 0

		rowChan := make(chan Row)
		// create wait group for row renderers
		var wg sync.WaitGroup
		wg.Add(jobsite.GRID_LENGTH)
		for r := 0; r < jobsite.GRID_LENGTH; r++ {
			// Handles row rendering in parallel
			go func(rowNum int, row [jobsite.GRID_WIDTH]string, rowChan chan Row) {
				defer wg.Done()
				for c := 0; c < jobsite.GRID_WIDTH; c++ {
					// Populates row with symbols
					row[c] = "   "
					isConeCoord := false
					if rowNum <= coneBounds.BottomBound && rowNum >= coneBounds.TopBound && c >= coneBounds.LeftBound && c < coneBounds.RightBound {
						row[c] = " X "
						isConeCoord = true
					}

					for _, val := range workerLocs {
						if c == val.X && rowNum == val.Y {
							row[c] = worker
						}
					}

					if row[c] == worker && isConeCoord {
						row[c] = warn
						numWorkersInCone += 1
					}
				}
				rowChan <- Row{RowNum: rowNum, Row: row}
			}(r, grid[r], rowChan)
		}
		// Assemble table
		go func() {
			for row := range rowChan {
				grid[row.RowNum] = row.Row
			}
		}()
		// Wait for all rows to render and arrange into grid
		wg.Wait()

		if numWorkersInCone == 0 {
			allGoodChan <- true
		} else {
			alertChan <- numWorkersInCone
		}

		table := tablewriter.NewWriter(os.Stdout)
		for _, v := range grid {
			table.Append(v[:])
		}
		var spaces [(jobsite.GRID_WIDTH / 2) - 6]string

		//Print grid
		cs := []string{" C ", " R ", " A ", " N ", " E ", " ", " S ", " A ", " F ", " E ", " T ", " Y "}
		header := append(spaces[:], cs...)
		header = append(header, spaces[:]...)
		table.SetHeader(header)
		table.SetRowLine(true)
		table.SetRowSeparator("-")
		table.SetAlignment(tablewriter.ALIGN_CENTER)

		fmt.Println()
		fmt.Printf("\033[4;1H")
		// Hackily fixing a print bug
		fmt.Print("                                                                                             ")
		fmt.Printf("\033[5;1H")
		table.Render()

	}
}
