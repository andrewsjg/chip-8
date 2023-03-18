package machine

type Cpu struct {
	PC uint16
	I  uint16
	V  [16]byte

	Stack      [16]uint16
	SP         byte
	DelayTimer byte
	SoundTimer byte
}
