package main

import "fmt"

type Config uint8
type transition struct {
	pathLength           int
	previousConfig       Config
	cellToPreviousConfig int
	config               Config
}

type cellMask uint8

const (
	cell8 cellMask = 1 << iota // order matters here
	cell7
	cell6
	cell5
	cell4
	cell3
	cell2
	cell1
)

/*
	 Cell numbering (which should help explain values in changeMask):
		1 2 3
		4 x 5
		6 7 8
*/
var changeMask = [8]Config{
	Config(cell1 | cell2 | cell4),                 // cell 1
	Config(cell1 | cell2 | cell3 | cell4 | cell5), // cell 2
	Config(cell2 | cell3 | cell5),                 // cell 3
	Config(cell1 | cell2 | cell4 | cell6 | cell7), // cell 4
	Config(cell2 | cell3 | cell5 | cell7 | cell8), // cell 5
	Config(cell4 | cell6 | cell7),                 // cell 6
	Config(cell4 | cell5 | cell6 | cell7 | cell8), // cell 7
	Config(cell5 | cell7 | cell8)}                 // cell 8

func configFromCellMask(mask cellMask) (config Config) {
	config = Config(mask)
	return
}

var pathList [256]transition
var configQueue [256]Config
var queueHead, queueTail int
var queueLen int
var configsDone [256]bool

func main() {
	config := transition{1, 255, -1, 255}
	pathList[255] = config
	addConfigToGeneratePaths(255)
	for queueLen != 0 {
		generateNextSetOfPathSteps()
	}
	printPaths()
}

func printPaths() {
	for config := 0; config < 255; config++ {
		if configsDone[config] {
			printConfigPath(Config(config))
		}

		if config%8 == 7 {
			fmt.Printf("\n")
		}
	}
}

func printConfigPath(config Config) {
	var path [8]int
	first := true

	printConfigName(config)
	count := 0
	for config != 255 {
		path[count] = pathList[config].cellToPreviousConfig+1
		count++
		config = pathList[config].previousConfig
	}
	for count != 0 {
		count--
		if first {
			fmt.Printf(": %d", path[count])
			first = false
		} else {
			fmt.Printf(", %d", path[count])
		}
	}
	fmt.Printf("\n")
}

func printConfigName(config Config) {
	cells := []cellMask{cell1, cell2, cell3, cell4, cell5, cell6, cell7, cell8}
	for i := 0; i < 8; i++ {
		printConfigLetter(config & Config(cells[i]))
		if i == 2 || i == 4 {
			fmt.Printf("_")
		}
	}
}

func printConfigLetter(mask Config) {
	if mask == 0 {
		fmt.Printf("r")
	} else {
		fmt.Printf("b")
	}

}

func generateNextSetOfPathSteps() {
	if queueHead == queueTail {
		panic("empty queue?")
	}
	config := configQueue[queueHead]
	// fmt.Printf("Looking at transitions for config %d\n", config);
	queueHead++
	if queueHead > len(configQueue) {
		queueHead = 0
	}
	queueLen--

	newConfigPathStep := pathList[config]
	newConfigPathStep.pathLength++
	newConfigPathStep.previousConfig = config

	for cell := 0; cell < 8; cell++ {
		newConfig := config ^ changeMask[cell]
		newConfigPathStep.cellToPreviousConfig = cell
		potentiallyAddNewPathStep(newConfigPathStep, newConfig)
	}
}

func potentiallyAddNewPathStep(newStep transition, newConfig Config) {
	if pathList[newConfig].pathLength != 0 && pathList[newConfig].pathLength <= newStep.pathLength {
		return
	}
	pathList[newConfig] = newStep
	addConfigToGeneratePaths(newConfig)
}

func addConfigToGeneratePaths(config Config) {
	if configsDone[config] {
		fmt.Printf("config %d already done\n", config)
		return
	} else {
		// fmt.Printf("adding config %d to the queue\n", config)
	}
	configsDone[config] = true
	queueLen++
	configQueue[queueTail] = config
	// design choices assure that the queue will never be full, so
	// skipping check and handling logic, and just panicking if things
	// go wrong (lazy, but safe)
	queueTail++
	if queueTail > len(configQueue) {
		queueTail = 0
	}
	if queueTail == queueHead {
		panic("overflowed queue")
	}
}
