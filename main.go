package main

import (
	"fmt"
	"os"

	"github.com/andrewsjg/chip-8/machine"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Please pass a Chip 8 program on the command line\n\nUsage: chip-8 <program_file>\n")
		return
	}

	program := os.Args[1]
	chip8 := machine.NewMachine()
	chip8.StartMachine(program, false)
}
