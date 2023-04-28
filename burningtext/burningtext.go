package burningtext

import (
	"image"
	"image/color"
	"math"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	opensimplex "github.com/ojrac/opensimplex-go"
)

const pixelsToPt = 1.0 / 0.75

// "#6a0500", "#c93900", "#ff9500", "#ffeb05"
var burningColors = []color.RGBA{
	{R: 106, G: 5, B: 0, A: 255},
	{R: 201, G: 57, B: 0, A: 255},
	{R: 255, G: 149, B: 0, A: 255},
	{R: 255, G: 235, B: 5, A: 255},
}

func lerpColorList(colors []color.RGBA, t float64) color.RGBA {
	if t <= 0 {
		return colors[0]
	}
	if t >= 1 {
		return colors[len(colors)-1]
	}

	t *= float64(len(colors) - 1)
	colorIndex := int(t)
	localT := t - float64(colorIndex)

	c1 := colors[colorIndex]
	c2 := colors[colorIndex+1]

	return color.RGBA{
		R: uint8(float64(c1.R)*(1-localT) + float64(c2.R)*localT),
		G: uint8(float64(c1.G)*(1-localT) + float64(c2.G)*localT),
		B: uint8(float64(c1.B)*(1-localT) + float64(c2.B)*localT),
		A: uint8(float64(c1.A)*(1-localT) + float64(c2.A)*localT),
	}
}

func easeInCirc(x float64) float64 {
	return 1 - math.Sqrt(1-math.Pow(x, 2))
}

const noiseScale = 1 / 3.

var noise = opensimplex.NewNormalized(0)

func getNoiseValue(x, y, z int) float64 {
	return noise.Eval3(float64(x)*noiseScale, float64(y)*noiseScale, float64(z)*noiseScale)
}

type BurningText struct {
	text string

	// width, height int
	rect image.Rectangle

	textMask   *image.Alpha
	fireBuffer *image.Alpha

	t int
}

func NewBurningText(text, fontPath string) *BurningText {
	fileBytes, err := os.ReadFile(fontPath)
	if err != nil {
		panic(err)
	}

	fontFace, err := opentype.Parse(fileBytes)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(fontFace, &opentype.FaceOptions{
		Size:    70 * pixelsToPt,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	width := font.MeasureString(face, text)

	imageWidth := int((float64(width) / 64) * 1.1)
	imageHeight := int(70 * 1.5)

	img := image.NewAlpha(image.Rect(0, 0, imageWidth, imageHeight))

	drawer := font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P((imageWidth/2)-(width.Round()/2), 70+16),
	}

	drawer.DrawString(text)

	rect := image.Rect(0, 0, int(imageWidth), int(imageHeight))

	return &BurningText{text: text, rect: rect, textMask: img, fireBuffer: image.NewAlpha(rect)}
}

// Draw returns the current frame of the animation
func (bt *BurningText) Draw() *image.RGBA {
	img := image.NewRGBA(bt.rect)

	for x := 0; x < bt.rect.Dx(); x++ {
		for y := 0; y < bt.rect.Dy(); y++ {
			if bt.textMask.AlphaAt(x, y).A > 0 {
				img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
			} else if bt.fireBuffer.AlphaAt(x, y).A > 0 {
				fireColor := lerpColorList(burningColors, float64(bt.fireBuffer.AlphaAt(x, y).A)/255)
				img.Set(x, y, fireColor)
			}
		}
	}

	return img
}

// Process updates the buffers for the next frame
func (bt *BurningText) Process() {
	newFireBuffer := image.NewAlpha(bt.rect)

	// Heat the text
	for x := 0; x < bt.rect.Dx(); x++ {
		for y := 0; y < bt.rect.Dy(); y++ {
			if bt.textMask.AlphaAt(x, y).A > 0 {
				newFireBuffer.SetAlpha(x, y, color.Alpha{A: byte((1-easeInCirc(getNoiseValue(x, y, bt.t)))*192) + 64})
			}
		}
	}

	// Move the fire up
	for x := 0; x < bt.rect.Dx(); x++ {
		for y := 0; y < bt.rect.Dy(); y++ {
			if bt.textMask.AlphaAt(x, y).A > 0 {
				continue
			}

			flameInfluence := float64(bt.fireBuffer.AlphaAt(x, y+1).A)*0.6 +
				float64(bt.fireBuffer.AlphaAt(x-1, y).A)*0.3 +
				float64(bt.fireBuffer.AlphaAt(x, y+1).A)*0.1

			if flameInfluence <= 1 {
				continue
			}

			value := flameInfluence - 1 - easeInCirc(getNoiseValue(x, y, bt.t))*64

			// Make sure the resulting value fits in a byte
			result := math.Round(math.Max(value, float64(newFireBuffer.AlphaAt(x, y).A)))

			newFireBuffer.SetAlpha(x, y, color.Alpha{A: byte(math.Min(math.Max(result, 0), 255))})
		}
	}

	bt.fireBuffer = newFireBuffer
	bt.t++
}
