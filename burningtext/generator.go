package burningtext

import (
	"encoding/json"
	"image"
	"image/draw"
	"os"
	"strings"
)

type BurningSpeed int

const (
	SpeedFast BurningSpeed = iota
	SpeedSlow
)

var fontChains map[string][]string = nil

func loadFontChainsData() {
	file, err := os.Open(os.Getenv("FONTCHAINS_PATH"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&fontChains)

	if err != nil {
		panic(err)
	}
}

func IsFontChainValid(fontChain string) bool {
	if fontChains == nil {
		loadFontChainsData()
	}

	_, ok := fontChains[fontChain]
	return ok
}

func StackImage(images []*image.RGBA) *image.RGBA {
	resultImage := image.NewRGBA(image.Rect(0, 0, images[0].Rect.Max.X, images[0].Rect.Max.Y*len(images)))

	for i, img := range images {
		draw.Draw(resultImage, image.Rect(0, img.Rect.Max.Y*i, img.Rect.Max.X, img.Rect.Max.Y*(i+1)), img, image.Point{0, 0}, draw.Src)
	}

	return resultImage
}

func getDefaultFontChain() []string {
	value, ok := os.LookupEnv("FONT")
	if !ok {
		panic("FONT_NAME environment variable not set")
	}
	return strings.Split(value, ";")
}
