package machine

type Display struct {
	Pixels [screenWidth][screenHeight]byte
}

const screenWidth = 64
const screenHeight = 32

func (d *Display) ScreenWidth() int {
	return screenWidth
}

func (d *Display) ScreenHeight() int {
	return screenHeight
}
