package burningtext

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

func GenerateSingleText(text string) {
	img := GenerateNeosSpritesheet(text)
	file, err := os.Create("out/flaming.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}

	for i, img := range GenerateAnimatedFrames(text) {
		file, err := os.Create(fmt.Sprintf("out/frames/img%d.png", i))
		if err != nil {
			panic(err)
		}
		defer file.Close()

		err = png.Encode(file, img)
		if err != nil {
			panic(err)
		}
	}
}

func GenerateNeosSpritesheet(text string) *image.RGBA {
	burningText := NewBurningText(text, filepath.Join(getFontDirectory(), getFontName()))

	for i := 0; i < 50; i++ {
		burningText.Process()
	}

	images := make([]*image.RGBA, 0)

	for i := 0; i < 5; i++ {
		for k := 0; k < 3; k++ {
			burningText.Process()
		}

		images = append(images, burningText.Draw())
	}

	return stackImage(images)
}

func stackImage(images []*image.RGBA) *image.RGBA {
	resultImage := image.NewRGBA(image.Rect(0, 0, images[0].Rect.Max.X, images[0].Rect.Max.Y*5))

	for i, img := range images {
		draw.Draw(resultImage, image.Rect(0, img.Rect.Max.Y*i, img.Rect.Max.X, img.Rect.Max.Y*(i+1)), img, image.Point{0, 0}, draw.Src)
	}

	return resultImage
}

func GenerateAnimatedFrames(text string) []*image.RGBA {
	burningText := NewBurningText(text, filepath.Join(getFontDirectory(), getFontName()))

	images := make([]*image.RGBA, 0)

	for i := 0; i < 100; i++ {
		burningText.Process()

		images = append(images, burningText.Draw())
	}

	return images
}

func getFontDirectory() string {
	value, ok := os.LookupEnv("FONT_DIRECTORY")
	if !ok {
		panic("FONT_DIRECTORY environment variable not set")
	}
	return value
}

func getFontName() string {
	value, ok := os.LookupEnv("FONT_NAME")
	if !ok {
		panic("FONT_NAME environment variable not set")
	}
	return value
}
