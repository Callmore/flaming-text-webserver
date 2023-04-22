package main

import (
	"flamingTextWebserver/burningtext"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

func generateSingleText(text string) {
	img := generateNeosSpritesheet(text)
	file, err := os.Create("out/flaming.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}

	for i, img := range generateAnimatedFrames(text) {
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

func generateNeosSpritesheet(text string) *image.RGBA {
	burningText := burningtext.NewBurningText(text, filepath.Join(fontDirectory, defaultFont))

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

func generateAnimatedFrames(text string) []*image.RGBA {
	burningText := burningtext.NewBurningText(text, filepath.Join(fontDirectory, defaultFont))

	images := make([]*image.RGBA, 0)

	for i := 0; i < 100; i++ {
		burningText.Process()

		images = append(images, burningText.Draw())
	}

	return images
}
