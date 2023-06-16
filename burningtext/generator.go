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

type BurningTextSettings struct {
	Text      string
	Speed     BurningSpeed
	FontChain string
}

var fontChains map[string][]string = nil

func (settings BurningTextSettings) GetFontChain() []string {
	if settings.FontChain == "" {
		return getDefaultFontChain()
	}

	if fontChains == nil {
		loadFontChainsData()
	}

	return fontChains[settings.FontChain]
}

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

func generateBurningTextImages(settings BurningTextSettings) []*image.RGBA {
	burningText := NewBurningText(settings.Text, settings.GetFontChain()[0])

	for i := 0; i < 50; i++ {
		burningText.Process()
	}

	images := make([]*image.RGBA, 0)

	switch settings.Speed {
	case SpeedFast:
		for i := 0; i < 5; i++ {
			for k := 0; k < 3; k++ {
				burningText.Process()
			}

			images = append(images, burningText.Draw())
		}

	case SpeedSlow:
		for i := 0; i < 50; i++ {
			burningText.Process()

			images = append(images, burningText.Draw())
		}
	}
	return images
}

func GenerateNeosSpritesheet(settings BurningTextSettings) *image.RGBA {
	images := generateBurningTextImages(settings)

	return stackImage(images)
}

func GenerateAnimatedFrames(settings BurningTextSettings) []*image.RGBA {
	images := generateBurningTextImages(settings)

	return images
}

func stackImage(images []*image.RGBA) *image.RGBA {
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
