package palletising

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"os"
)

var palette color.Palette = nil

func Palletise(inImage *image.RGBA) *image.Paletted {
	if palette == nil {
		palette = LoadJSONPalette(os.Getenv("PALETTE_PATH"))
	}

	bounds := inImage.Bounds()
	paletted := image.NewPaletted(bounds, palette)

	draw.Draw(paletted, bounds, inImage, image.Point{}, draw.Src)

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
