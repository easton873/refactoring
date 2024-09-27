package heatingsim

import (
	"fmt"
	"log"
	"os"
)

func Run() {
	my2DArray := [10][10]float64{}

	fmt.Println("Initializing heating simulator")
	for i := 0; i < 10; i++ {
		for k := 0; k < 10; k++ {
			if (k != 0 && k != 9) && (i == 0 || i == 9) {
				my2DArray[i][k] = 100.0
			}
		}
	}

	for i := 0; i < 10; i++ {
		for k := 0; k < 10; k++ {
			fmt.Printf("%10.6f", my2DArray[i][k])
			if k != 9 {
				fmt.Printf(", ")
			}
		}
		fmt.Println()
	}

	copyPlate := my2DArray

	for i := 0; i < 10; i++ {
		for k := 0; k < 10; k++ {
			if i != 0 && i != 9 && k != 0 && k != 9 {
				copyPlate[i][k] = (my2DArray[i+1][k] + my2DArray[i-1][k] + my2DArray[i][k+1] + my2DArray[i][k-1]) / 4
			}
		}
	}
	my2DArray = copyPlate

	fmt.Println()
	fmt.Println("heating simulator after one iteration")
	for i := 0; i < 10; i++ {
		for k := 0; k < 10; k++ {
			fmt.Printf("%10.6f", my2DArray[i][k])
			if k != 9 {
				fmt.Printf(", ")
			}
		}
		fmt.Println()
	}

	for i := 0; i < 10; i++ {
		for k := 0; k < 10; k++ {
			if i != 0 && i != 9 && k != 0 && k != 9 {
				copyPlate[i][k] = (my2DArray[i+1][k] + my2DArray[i-1][k] + my2DArray[i][k+1] + my2DArray[i][k-1]) / 4
			}
		}
	}
	my2DArray = copyPlate

	fmt.Println()
	fmt.Println("heating simulator after second iteration")
	for i := 0; i < 10; i++ {
		for k := 0; k < 10; k++ {
			fmt.Printf("%10.6f", my2DArray[i][k])
			if k != 9 {
				fmt.Printf(", ")
			}
		}
		fmt.Println()
	}

	done := true
	for done {
		biggestChange := 0.0

		for i := 0; i < 10; i++ {
			for k := 0; k < 10; k++ {
				if i != 0 && i != 9 && k != 0 && k != 9 {
					copyPlate[i][k] = (my2DArray[i+1][k] + my2DArray[i-1][k] + my2DArray[i][k+1] + my2DArray[i][k-1]) / 4
				}
			}
		}

		for i := 0; i < 10; i++ {
			for k := 0; k < 10; k++ {
				if i != 0 && i != 9 && k != 0 && k != 9 {
					if currChange := copyPlate[i][k] - my2DArray[i][k]; currChange > biggestChange {
						biggestChange = currChange
					}
				}
			}
		}
		my2DArray = copyPlate

		if biggestChange < 0.1 {
			done = false
		}
	}

	fmt.Println()
	fmt.Println("heating simulator in stable state")
	for i := 0; i < 10; i++ {
		for k := 0; k < 10; k++ {
			fmt.Printf("%10.6f", my2DArray[i][k])
			if k != 9 {
				fmt.Printf(", ")
			}
		}
		fmt.Println()
	}

	fp, err := os.Create("heating simulator.csv")
	defer fp.Close()
	if err == nil {
		fmt.Println()
		fmt.Println("Outputting to file \"heating simulator.csv\"")
		for i := 0; i < 10; i++ {
			for k := 0; k < 10; k++ {
				_, err = fp.WriteString(fmt.Sprintf("%10.6f", my2DArray[i][k]))
				if err != nil {
					log.Println("unable to write to file: ", err)
				}
				if k != 9 {
					_, err = fp.WriteString(", ")
					if err != nil {
						log.Println("unable to write to file: ", err)
					}
				}
			}
			_, err = fp.WriteString("\n")
			if err != nil {
				log.Println("unable to write to file: ", err)
			}
		}
	}
}
