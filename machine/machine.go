package machine

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// Font array size. 15 sets of 5 bytes
const fontCount = 80

// The loaded program code starts at memory location 0x200
const programStart = 0x200

// 4kb of memory
const memSize = 4096

// Machine configuration. Consists of memory, CPU, Display, input
// a debug flag. Also a flag to determine if the machine runs
// as specified by the orginginal specification
type Machine struct {
	Memory  [memSize]byte
	Cpu     Cpu
	Display Display
	Input   Input
	Debug   bool

	// Use the original machine implementation or not
	Original bool
}

// The font array. Contains the bytes that define the
// font bitmaps
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

// Create a new machine
func NewMachine() *Machine {
	machine := Machine{}

	// Setup the memory
	memory := [memSize]byte{}

	// Load the font into memory
	for i := 0; i < fontCount; i++ {
		memory[i] = font[i]
	}

	machine.Memory = memory

	// Set the program counter to the location where we expect the program
	// in  memory
	machine.Cpu.PC = 0x200

	// TODO: Make this configurable?
	machine.Original = false

	return &machine
}

// Method to start the machine.
// Turns on debug as specified. Sets up the display and starts
// the main loop. Controlled by ebiten.
func (m *Machine) StartMachine(program string, debugOn bool) {
	log.Println("Starting machine")

	m.Debug = debugOn
	m.loadProgram(program)

	// The Chip-8 machine has a 64x32 display. Which is tiny on large
	// computer screens. So we scale it up by some factor.

	const scaleFactor = 10

	screenWidth := m.Display.ScreenWidth() * scaleFactor
	screenHeight := m.Display.ScreenHeight() * scaleFactor

	// Create the display
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(m); err != nil {
		panic(err)
	}

}

// Handle any kepresses during a cycle.
// The Chip-8 Machine has 16 keys. Check if any are pressed
// and set the corresponding value in our input array accordingly
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

// Load a program into memory. This simply reads a file byte by byte
// and loads each byte (instruction) into memory
func (m *Machine) loadProgram(fileName string) {

	program, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	for i, instruction := range program {
		location := programStart + i
		if location > memSize {
			panic("Program wont fit into memory")
		}

		m.Memory[location] = instruction
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

// The decode and execute function. These are done in one function rather
// than two separate functions. This is the way it makes most sense to do
// and is recommended by the HOWTO post TODO: Add refreence here
func (m *Machine) decodeAndExecute(instruction uint16) {

	// TODO: Annotate this
	nnn := instruction & 0xFFF
	n := byte(instruction & 0xF)
	x := byte(instruction & 0xF00 >> 8)
	y := byte(instruction & 0xF0 >> 4)
	nn := byte(instruction & 0xFF)

	if instruction != 0 {
		switch instruction {

		// These two opcodes have no parameters so no further decoding required
		case 0x00E0:
			if m.Debug {
				fmt.Println("Clear Screen")
			}

			m.Display.Clear()

		// Return from a subroutine sets the program counter
		// to the location of ths stack pointer and decrements the stack pointer
		case 0x00EE:
			if m.Debug {
				log.Println("Return from Subroutine")
			}
			m.Cpu.PC = uint16(m.Cpu.Stack[m.Cpu.SP])
			m.Cpu.SP--

		// Fall into the default case for further decoding
		// This set of opcodes require decoding the first nibble (half byte)
		default:
			// Mask the instruction and shift right extract the first nibble
			switch byte(instruction & 0xF000 >> 12) {
			case 0x1:
				if m.Debug {
					log.Printf("Jump to nnn: %d\n", nnn)
				}
				// Set the program counter to loction nnn
				m.Cpu.PC = nnn

			case 0x2:
				if m.Debug {
					log.Println("Call subroutine at nnn")
				}
				// Increment the stack pointer
				m.Cpu.SP++

				// Put the current program counter location onto the stack
				m.Cpu.Stack[m.Cpu.SP] = m.Cpu.PC

				// The the program counter to nnn
				m.Cpu.PC = nnn

			case 0x3:
				if m.Debug {
					log.Println("Skip instruction if V[x] == nn")
				}

				// Check the Vx register against nn. If it is equal, skip
				// over an instruction by incrementgin the program counter by 2
				if m.Cpu.V[x] == nn {
					m.Cpu.PC += 2
				}
			case 0x4:
				if m.Debug {
					log.Println("Skip instruction if V[x] != nn")
				}
				// As above but check if Vx is not equal to nn
				if m.Cpu.V[x] != nn {
					m.Cpu.PC += 2
				}
			case 0x5:
				if m.Debug {
					log.Println("Skip instruction if V[x] == V[y]")
				}

				// Skip over an instruction by incrementing the program counter by
				// by 2 if register Vx is equal to Vy
				if m.Cpu.V[x] == m.Cpu.V[y] {
					m.Cpu.PC += 2
				}

			case 0x6:
				if m.Debug {
					log.Println("Set VX to nn")
				}
				// The Vx to the value of nn
				m.Cpu.V[x] = nn

			case 0x7:
				if m.Debug {
					log.Println("Add nn to vx")
				}

				// Add nn to the value of register Vx
				m.Cpu.V[x] += nn

			// The next block of instructions are the
			// Logical and arithmetic instructions and require further extraction
			case 0x8:
				// Mask to extract the last nibble
				switch instruction & 0xF {
				case 0x0:
					if m.Debug {
						log.Println("Set V[X] to V[Y]")
					}
					m.Cpu.V[x] = m.Cpu.V[y]

				case 0x1:
					if m.Debug {
						log.Println("Set V[X] to bitwise logical disjunction (OR) of V[X] and V[Y]")
					}
					m.Cpu.V[x] |= m.Cpu.V[y]

				case 0x2:
					if m.Debug {
						log.Println("Set V[X] to bitwise logical conjunction (AND) of V[X] and V[Y]")
					}
					m.Cpu.V[x] &= m.Cpu.V[y]

				case 0x3:
					if m.Debug {
						log.Println("Set V[X] to bitwise exclusive or (XOR)) of V[X] and V[Y]")
					}
					m.Cpu.V[x] ^= m.Cpu.V[y]

				case 0x4:
					if m.Debug {
						log.Println("Set V[X] to the value of V[X] + V[Y]")
					}

					result := uint16(m.Cpu.V[x]) + uint16(m.Cpu.V[y])

					// Need to set overflow flag if required
					// If the result of the addition operation is greater than 255
					// set the overflow flag to 1

					// TODO: Document the overflow flag
					m.Cpu.V[0xf] = 0
					if result > 0xFFFF {
						m.Cpu.V[0xf] = 1
					}

					m.Cpu.V[x] = byte(result)

				case 0x5:
					if m.Debug {
						log.Println("Set V[X] to the result of V[X] - V[Y]")
					}
					// Set the overflow flag. We will unset conditionally below
					m.Cpu.V[0xf] = 1

					//TODO: Document this
					minuend := uint16(m.Cpu.V[x])
					subtrahend := uint16(m.Cpu.V[y])

					if subtrahend > minuend {
						m.Cpu.V[0xf] = 0
					}

					result := minuend - subtrahend
					m.Cpu.V[x] = byte(result)

				case 0x6:
					if m.Debug {
						log.Println("Right Shift V[X]]")
					}

					// If the machine is setup to follow the original
					// spec, set the Vx register to Vy
					if m.Original {
						m.Cpu.V[x] = m.Cpu.V[y]
					}

					// Set the overflow bit
					m.Cpu.V[0xF] = (m.Cpu.V[x] & 0xF)

					// Shift V[X] right by 1
					m.Cpu.V[x] >>= 1

				case 0x7:
					if m.Debug {
						log.Println("Set V[X] to the result of V[Y] - V[X]")
					}

					// The opposite of the above. Vy - Vx as opposed to Vx - Vy
					m.Cpu.V[0xf] = 1

					minuend := uint16(m.Cpu.V[y])
					subtrahend := uint16(m.Cpu.V[x])

					if subtrahend > minuend {
						m.Cpu.V[0xf] = 0
					}

					result := minuend - subtrahend
					m.Cpu.V[x] = byte(result)

				case 0xE:
					if m.Debug {
						log.Println("Left Shift V[X]]")
					}

					// The opposite of the above. Left shift Vx as opposed to a
					// right shift

					if m.Original {
						m.Cpu.V[x] = m.Cpu.V[y]
					}

					// Set the overflow bit
					m.Cpu.V[0xF] = (m.Cpu.V[x] >> 7)

					// Shift V[X] left by 1
					m.Cpu.V[x] <<= 1

				// If we get something unexpected simply output an error.
				// This case should never be hit.
				default:
					log.Printf("Instruction not implemented %s0x\n", fmt.Sprintf("%X", instruction))
				}
			// The next block dont require extracting the last nibble
			case 0x9:
				if m.Debug {
					log.Println("Skip instruction if V[x] != V[y]")
				}
				if m.Cpu.V[x] != m.Cpu.V[y] {
					m.Cpu.PC += 2
				}

			case 0xA:
				if m.Debug {
					log.Println("Set index register")
				}
				m.Cpu.I = nnn

			case 0xB:
				if m.Debug {
					log.Println("Jump with offset")
				}

				if m.Original {
					m.Cpu.PC = nnn + uint16(m.Cpu.V[0x0])
				} else {
					m.Cpu.PC = nnn + uint16(m.Cpu.V[x])
				}

			case 0xC:
				if m.Debug {
					log.Println("Generate a random number and AND it with nn and store in V[X]")
				}
				// Generate random byte between 0 and 255
				rnd := byte(rand.Uint32() % 255)
				m.Cpu.V[x] = rnd & nn

			case 0xD:
				// From: https://github.com/szTheory/chip8go
				if m.Debug {
					log.Println("Display / Draw")
				}

				xVal := m.Cpu.V[x]
				yVal := m.Cpu.V[y]

				m.Cpu.V[0xF] = 0

				var i byte = 0
				for ; i < n; i++ {
					row := m.Memory[m.Cpu.I+uint16(i)]

					if erased := m.Display.DrawSprite(xVal, yVal+i, row); erased {
						m.Cpu.V[0xF] = 1
					}
				}
			case 0xE:
				// Mask to extract the last nibble
				switch instruction & 0xFF {
				case 0x9E:
					if m.Debug {
						log.Println("Skip if key correspoding to V[X] is pressed")
					}
					if m.Input[m.Cpu.V[x]] {
						m.Cpu.PC += 2
					}
				case 0xA1:
					if m.Debug {
						log.Println("Skip if key correspoding to V[X] is Not pressed")
					}
					if !m.Input[m.Cpu.V[x]] {
						m.Cpu.PC += 2
					}

				default:
					log.Printf("Instruction not implemented %s0x\n", fmt.Sprintf("%X", instruction))
				}

			case 0xF:
				// Mask to extract the last nibble
				switch instruction & 0xFF {
				case 0x07:
					if m.Debug {
						log.Println("Set V[X] to the value of the delay timer")
					}
					m.Cpu.V[x] = m.Cpu.DelayTimer

				case 0x15:
					if m.Debug {
						log.Println("Set the delay timer to the value of V[X]")
					}
					m.Cpu.DelayTimer = m.Cpu.V[x]

				case 0x18:
					if m.Debug {
						log.Println("Set the sound timer to the value of V[X]")
					}
					m.Cpu.SoundTimer = m.Cpu.V[x]

				case 0x1E:
					if m.Debug {
						log.Println("Add V[X] to the value of the index register and store in the index register")
					}
					m.Cpu.I += uint16(m.Cpu.V[x])

					// Do it the way the Amiga emulator did and set the overflow flag if the value is over 1000
					if m.Cpu.I > 1000 {
						m.Cpu.V[0xF] = 1
					}

				case 0x0A:
					if m.Debug {
						log.Println("Block excution of a key is pressed")
					}

					keyPressed := false
					for i := uint8(0); i < 16; i++ {
						if m.Input[i] {
							m.Cpu.V[x] = i
							keyPressed = true
						}
					}
					if !keyPressed {
						m.Input.enableWait(true)
					}

				case 0x29:
					if m.Debug {
						log.Println("Font Character. Set the index register to the address of the character in V[X]")
					}

					m.Cpu.I = uint16(m.Cpu.V[x] * 0x5)

				case 0x33:
					if m.Debug {
						log.Println("Binary coded decimal conversion")
					}
					//TODO: Anotate this
					m.Memory[m.Cpu.I] = m.Cpu.V[x] / 100
					m.Memory[m.Cpu.I+1] = (m.Cpu.V[x] / 10) % 10
					m.Memory[m.Cpu.I+2] = (m.Cpu.V[x] % 100) % 10

				case 0x55:
					if m.Debug {
						log.Println("Copy V registers to memory")
					}
					for i := byte(0); i < x; i++ {
						m.Memory[m.Cpu.I+uint16(i)] = m.Cpu.V[i]

						if m.Original {
							m.Cpu.I++
						}
					}

				case 0x65:
					if m.Debug {
						fmt.Println("Copy memory to V registers")
					}

					for i := byte(0); i < x; i++ {
						m.Cpu.V[i] = m.Memory[m.Cpu.I+uint16(i)]

						if m.Original {
							m.Cpu.I++
						}
					}
				default:
					log.Printf("Instruction not implemented %s0x\n", fmt.Sprintf("%X", instruction))
				}
			default:
				log.Printf("Instruction not implemented %s0x\n", fmt.Sprintf("%X", instruction))
			}

		}
	}
}

func (m *Machine) machineCycle() {

	// Update the timers
	if m.Cpu.DelayTimer > 0 {
		m.Cpu.DelayTimer--
	}

	if m.Cpu.SoundTimer > 0 {
		m.Cpu.SoundTimer--
	}

	// If waiting for input, skip the cycle
	if m.Input.wait() {
		return
	}

	// Fetch, Decode and Execute
	m.decodeAndExecute(m.fetch())

}

// The methods below are required by ebiten
// to satisfy the interface. They are called by
// the ebiten program loop
func (m *Machine) Update() error {

	m.handleInput()
	m.machineCycle()

	return nil
}

// From: https://github.com/szTheory/chip8go
func (m *Machine) Draw(screen *ebiten.Image) {

	canvas := ebiten.NewImage(m.Display.ScreenWidth(), m.Display.ScreenHeight())

	for x := 0; x < m.Display.ScreenWidth(); x++ {
		for y := 0; y < m.Display.ScreenHeight(); y++ {
			setColor := color.Black
			if m.Display.Pixels[x][y] == 1 {
				setColor = color.White
			}
			if setColor != canvas.At(x, y) {
				canvas.Set(x, y, setColor)
			}

		}
	}

	geometry := ebiten.GeoM{}
	screen.DrawImage(canvas, &ebiten.DrawImageOptions{GeoM: geometry})

}

func (m *Machine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return m.Display.ScreenWidth(), m.Display.ScreenHeight()
}
