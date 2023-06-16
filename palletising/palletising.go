package palletising

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"os"
)

func Palletise(inImage *image.RGBA, palette color.Palette) *image.Paletted {
	bounds := inImage.Bounds()
	paletted := image.NewPaletted(bounds, palette)

	draw.Draw(paletted, bounds, inImage, image.Point{}, draw.Src)
	// draw.FloydSteinberg.Draw(paletted, bounds, inImage, image.Point{})

	return paletted
}

func LoadJSONPalette(path string) color.Palette {
	palette := color.Palette{}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var colors []color.RGBA
	err = decoder.Decode(&colors)

	if err != nil {
		panic(err)
	}

	for _, c := range colors {
		palette = append(palette, c)
	}

	return palette
}
