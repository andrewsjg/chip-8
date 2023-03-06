package machine

type Cpu struct {
	PC uint16
	I  uint16
	V  [16]uint16

	Stack [16]uint16

	DelayTimer byte
	SoundTimer byte
}
