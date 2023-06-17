package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"flamingTextWebserver/burningtext"
	"flamingTextWebserver/palletising"

	"github.com/joho/godotenv"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	cpuProfile     = flag.String("cpuprofile", "", "write cpu profile to file")
	textToGenerate = flag.String("text", "", "Runs the generator once, generating the first 100 frames of the animation, and a sheet of 5 frames. Outputs into a folder named \"out\" if it exists. Stops the webserver from starting.")
	outputFormat   = flag.String("format", "gif", "Format of the output. Can be \"gif\" or \"neos\"")
	burningSpeed   = flag.String("speed", "fast", "Speed of the animation. Can be \"slow\" or \"fast\"")
	outputLocation = flag.String("output", "", "Location to output the generated files to. Defaults to \"out\"")
	fontChain      = flag.String("font", "agate", "Name of the font chain to use. Font chains can be specified in fonts.json.")

	flameColorFlag = flag.String("flame-color", "", "Color of the flames as a hex string.")
	textColorFlag  = flag.String("text-color", "", "Color of the text as a hex string.")
	randomColor    = flag.Bool("random-color", false, "Randomize the colors of the flames and text.")
)

func init() {
	flag.Parse()
	godotenv.Load()
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = pprof.StartCPUProfile(f)
		if err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	if *textToGenerate == "" {
		panic("no text provided to generate")
	}

	fontChain := *fontChain

	speed := burningtext.SpeedFast
	if (*burningSpeed) == "slow" {
		speed = burningtext.SpeedSlow
	}

	var flameColor color.Color = nil
	if *flameColorFlag != "" {
		var err error
		flameColor, err = colorful.Hex(*flameColorFlag)
		if err != nil {
			panic(err)
		}
	}

	var textColor color.Color = nil
	if *textColorFlag != "" {
		var err error
		textColor, err = colorful.Hex(*textColorFlag)
		if err != nil {
			panic(err)
		}
	}

	options := burningtext.BurningTextOptions{
		Text:      *textToGenerate,
		Speed:     speed,
		FontChain: fontChain,

		FlameColor:  flameColor,
		TextColor:   textColor,
		RandomColor: *randomColor,
	}

	outputBuffer := bytes.Buffer{}
	outputWriter := bufio.NewWriter(&outputBuffer)

	burningText := burningtext.New(&options)

	switch *outputFormat {
	case "gif":
		frames := burningText.GenerateBurningTextImages()
		delaySlice := make([]int, len(frames))
		disposal := make([]byte, len(frames))
		palettedFrames := make([]*image.Paletted, len(frames))
		pal := burningText.GeneratePalette()
		for i := 0; i < len(frames); i++ {
			if speed == burningtext.SpeedFast {
				delaySlice[i] = 5
			} else {
				delaySlice[i] = 2
			}

			disposal[i] = gif.DisposalPrevious
			palettedFrames[i] = palletising.Palletise(frames[i], pal)
		}

		// TODO: Implement palletization
		gif.EncodeAll(outputWriter, &gif.GIF{
			Image:    palettedFrames,
			Delay:    delaySlice,
			Disposal: disposal,
		})

	case "neos":
		frames := burningText.GenerateBurningTextImages()
		image := burningtext.StackImage(frames)
		png.Encode(outputWriter, image)
	default:
		panic(fmt.Sprintf("unknown output format %s", *outputFormat))
	}

	outputWriter.Flush()

	if *outputLocation != "" {
		file, err := os.Create(*outputLocation)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = file.Write(outputBuffer.Bytes())
		if err != nil {
			panic(err)
		}

		return
	}

	os.Stdout.Write(outputBuffer.Bytes())
}
