package palletising

import (
	"image"
	"image/color"
	"image/draw"
)

func Palletise(inImage *image.RGBA, palette color.Palette) *image.Paletted {
	bounds := inImage.Bounds()
	paletted := image.NewPaletted(bounds, palette)

	draw.Draw(paletted, bounds, inImage, image.Point{}, draw.Src)

	return paletted
}
