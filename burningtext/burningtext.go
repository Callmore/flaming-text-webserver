package burningtext

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	colorful "github.com/lucasb-eyer/go-colorful"
	opensimplex "github.com/ojrac/opensimplex-go"
)

const pixelsToPt = 1.0 / 0.75

func lerpColorList(colors []color.Color, t float64) color.Color {
	if t <= 0 {
		return colors[0]
	}
	if t >= 1 {
		return colors[len(colors)-1]
	}

	t *= float64(len(colors) - 1)
	colorIndex := int(t)
	localT := t - float64(colorIndex)

	c1, _ := colorful.MakeColor(colors[colorIndex])
	c2, _ := colorful.MakeColor(colors[colorIndex+1])

	return c1.BlendLab(c2, localT).Clamped()
}

var (
	colorfulWhite = colorful.Color{R: 1, G: 1, B: 1}
	colorfulBlack = colorful.Color{R: 0, G: 0, B: 0}
)

// generateColorList generates a list of 4 colors that can be used to color the fire. The base colour is the 3rd colour in the list.
func generateColorList(base color.Color) []color.Color {
	baseColor, _ := colorful.MakeColor(base)

	baseH, baseS, baseL := baseColor.Hsl()

	brighter := colorful.Hsl(wrap(baseH+30, 360), baseS, math.Min(baseL+0.05, 1))

	darker := colorful.Hsl(wrap(baseH-30, 360), baseS, baseL)
	darker = darker.BlendLab(colorfulBlack, 0.5)

	darkest := colorful.Hsl(wrap(baseH-45, 360), baseS, baseL)
	darkest = darkest.BlendLab(colorfulBlack, 0.75)

	return []color.Color{
		darkest.Clamped(),
		darker.Clamped(),
		base,
		brighter.Clamped(),
	}
}

func (bt *BurningText) GeneratePalette() []color.Color {
	pal := color.Palette{}

	pal = append(pal, color.Alpha{A: 0})

	for i := 0; i < 254; i++ {
		if (float64(i) / 253) > 1 {
			panic("i is too big")
		}
		pal = append(pal, lerpColorList(bt.fireColors, float64(i)/253))
	}

	pal = append(pal, bt.textColor)

	return pal
}

func easeInCirc(x float64) float64 {
	return 1 - math.Sqrt(1-math.Pow(x, 2))
}

const noiseScale = 1 / 3.

var noise = opensimplex.NewNormalized(0)

func getNoiseValue(x, y, z int) float64 {
	return noise.Eval3(float64(x)*noiseScale, float64(y)*noiseScale, float64(z)*noiseScale)
}

func wrap(x float64, max float64) float64 {
	ret := math.Mod(x, max)
	if ret < 0 {
		ret += max
	}
	return ret
}

type BurningText struct {
	text string

	// width, height int
	rect image.Rectangle

	textMask   *image.Alpha
	fireBuffer *image.Alpha

	speed BurningSpeed

	t int

	fireColors []color.Color
	textColor  color.Color
}

func New(options *BurningTextOptions) *BurningText {
	fontPath := options.GetFontChain()[0]

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

	bounds, width := font.BoundString(face, options.Text)

	imageWidth := int((float64(bounds.Max.X.Ceil() - bounds.Min.X.Ceil())) + 20)
	imageHeight := int((float64(bounds.Max.Y.Ceil() - bounds.Min.Y.Ceil())) + 20)

	img := image.NewAlpha(image.Rect(0, 0, imageWidth, imageHeight))

	drawer := font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P((imageWidth/2)-(((width.Ceil())/2)+bounds.Min.X.Ceil()/2), -bounds.Min.Y.Ceil()+19),
	}

	drawer.DrawString(options.Text)

	rect := image.Rect(0, 0, int(imageWidth), int(imageHeight))

	var flameColor color.Color = color.RGBA{R: 255, G: 149, B: 0, A: 255}
	var textColor color.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}

	if options.RandomColor {
		baseH := rand.Float64() * 360
		flameColor = colorful.Hsv(wrap(baseH, 360), 1, 1).Clamped()
		textColor = colorful.Hsv(wrap(baseH-30, 360), 1, 1).Clamped()
	} else {
		if options.FlameColor != nil {
			flameColor = options.FlameColor
		}
		if options.TextColor != nil {
			textColor = options.TextColor
		}
	}

	return &BurningText{
		text:       options.Text,
		rect:       rect,
		textMask:   img,
		fireBuffer: image.NewAlpha(rect),

		speed:      options.Speed,
		fireColors: generateColorList(flameColor),
		textColor:  textColor,
	}
}

// Draw returns the current frame of the animation
func (bt *BurningText) Draw() *image.RGBA {
	img := image.NewRGBA(bt.rect)

	for x := 0; x < bt.rect.Dx(); x++ {
		for y := 0; y < bt.rect.Dy(); y++ {
			if bt.textMask.AlphaAt(x, y).A > 0 {
				img.Set(x, y, bt.textColor) // color.RGBA{R: 255, G: 0, B: 0, A: 255})
			} else if bt.fireBuffer.AlphaAt(x, y).A > 0 {
				fireColor := lerpColorList(bt.fireColors, float64(bt.fireBuffer.AlphaAt(x, y).A)/255)
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

func (bt *BurningText) GenerateBurningTextImages() []*image.RGBA {
	for i := 0; i < 50; i++ {
		bt.Process()
	}

	images := make([]*image.RGBA, 0)

	switch bt.speed {
	case SpeedFast:
		for i := 0; i < 5; i++ {
			for k := 0; k < 3; k++ {
				bt.Process()
			}

			images = append(images, bt.Draw())
		}

	case SpeedSlow:
		for i := 0; i < 50; i++ {
			bt.Process()

			images = append(images, bt.Draw())
		}
	}
	return images
}
