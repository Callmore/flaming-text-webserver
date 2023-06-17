package burningtext

import (
	"image/color"
	"testing"
)

func TestColorLerp(t *testing.T) {
	colorList := []color.Color{
		color.RGBA{R: 0, G: 0, B: 0, A: 255},
		color.RGBA{R: 255, G: 0, B: 0, A: 255},
		color.RGBA{R: 0, G: 255, B: 0, A: 255},
		color.RGBA{R: 0, G: 0, B: 255, A: 255},
	}

	t.Log(lerpColorList(colorList, 0));
	t.Log(lerpColorList(colorList, 0.333334));
	t.Log(lerpColorList(colorList, 0.666667));
	t.Log(lerpColorList(colorList, 1));
}
