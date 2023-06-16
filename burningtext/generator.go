package burningtext

import (
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
	FontChain []string
}

func generateBurningTextImages(settings BurningTextSettings) []*image.RGBA {
	if settings.FontChain == nil {
		settings.FontChain = getDefaultFontChain()
	}

	burningText := NewBurningText(settings.Text, settings.FontChain[0])

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
