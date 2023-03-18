package main

import "github.com/andrewsjg/chip-8/machine"

func main() {

	machine := machine.NewMachine()

	//machine.StartMachine("programs/IBM Logo.ch8")
	machine.StartMachine("programs/outlaw.ch8", false)
}
