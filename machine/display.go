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

func (d *Display) Clear() {

	for x := 0; x < screenWidth; x++ {
		for y := 0; y < screenHeight; y++ {
			d.Pixels[x][y] = 0
		}
	}

}

// From: https://github.com/szTheory/chip8go
func (d *Display) DrawSprite(x byte, y byte, row byte) bool {
	erased := false
	yIndex := y % screenHeight

	for i := x; i < x+8; i++ {
		xIndex := i % screenWidth

		wasSet := d.Pixels[xIndex][yIndex] == 1
		value := row >> (x + 8 - i - 1) & 1

		d.Pixels[xIndex][yIndex] ^= value

		if wasSet && d.Pixels[xIndex][yIndex] == 0 {
			erased = true
		}
	}

	return erased
}
