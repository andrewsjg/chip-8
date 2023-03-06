package machine

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Machine struct {
	Memory  [4096]byte
	Cpu     Cpu
	Display Display
	Input   Input
}

const fontCount = 80

var font [fontCount]uint8 = [fontCount]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func NewMachine() *Machine {
	machine := Machine{}

	memory := [4096]byte{}

	// Load the font into memory
	for i := 0; i < fontCount; i++ {
		memory[i] = font[i]
	}

	machine.Memory = memory

	// Set the program counter to the location of the start of the program
	// in  memory

	machine.Cpu.PC = 0x200

	return &machine
}

func (m *Machine) StartMachine() {
	log.Println("Starting machine")

	const scaleFactor = 10

	screenWidth := m.Display.ScreenWidth() * scaleFactor
	screenHeight := m.Display.ScreenHeight() * scaleFactor

	// Create the display
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(m); err != nil {
		panic(err)
	}

}

func (m *Machine) handleInput() {
	m.Input = [16]bool{
		ebiten.IsKeyPressed(ebiten.KeyX),
		ebiten.IsKeyPressed(ebiten.Key1),
		ebiten.IsKeyPressed(ebiten.Key2),
		ebiten.IsKeyPressed(ebiten.Key3),
		ebiten.IsKeyPressed(ebiten.KeyQ),
		ebiten.IsKeyPressed(ebiten.KeyW),
		ebiten.IsKeyPressed(ebiten.KeyE),
		ebiten.IsKeyPressed(ebiten.KeyA),
		ebiten.IsKeyPressed(ebiten.KeyS),
		ebiten.IsKeyPressed(ebiten.KeyD),
		ebiten.IsKeyPressed(ebiten.KeyZ),
		ebiten.IsKeyPressed(ebiten.KeyC),
		ebiten.IsKeyPressed(ebiten.Key4),
		ebiten.IsKeyPressed(ebiten.KeyR),
		ebiten.IsKeyPressed(ebiten.KeyF),
		ebiten.IsKeyPressed(ebiten.KeyV),
	}
}

func (m *Machine) fetch() (instruction uint16) {

	// A Chip-8 nstruction is two bytes. To fetch an instruction from memory
	// we need to read the byte the program counter is pointing to and the byte
	// following and combine into a single 16 bit instruction.

	// The line below reads the byte at the progranm counter into a 16 bit uint and shifts it 8 bits
	// to the left. This creats a uint with the first 8 bits set to the value read at the memory location
	// pointed to by they program counter followed by 8 zero bits. This is then logically ORed with
	// the value at the memory location one higher than the program counter which has been read into
	// another 16 bit uint but not shifted. This means the second byte read in will be a uint with
	// 8 zero bits, followed by the bits read in. OR'ing this with the previous value will result in
	// a 16 bit instruction from the compined bytes

	instruction = uint16(m.Memory[m.Cpu.PC])<<8 | uint16(m.Memory[m.Cpu.PC+1])

	// Increment the program counter by 2
	m.Cpu.PC += 2

	return instruction
}

func (m *Machine) decodeAndExecute(instriction uint16) {

}

func (m *Machine) machineCycle() {

	m.decodeAndExecute(m.fetch())
}

// Instruction Set

// Methods required by ebiten
func (m *Machine) Update() error {

	m.handleInput()
	m.machineCycle()

	return nil
}

func (m *Machine) Draw(screen *ebiten.Image) {

}

func (m *Machine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return m.Display.ScreenHeight(), m.Display.ScreenWidth()
}
