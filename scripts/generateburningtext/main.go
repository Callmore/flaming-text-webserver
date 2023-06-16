package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"os"
	"runtime/pprof"
	"strings"

	"flamingTextWebserver/burningtext"
	"flamingTextWebserver/palletising"

	"github.com/joho/godotenv"
)

// Set here to avoid a windows defender false positive (i'm not kidding it's so dumb)
const gifSpeed = 2

var (
	cpuProfile     = flag.String("cpuprofile", "", "write cpu profile to file")
	textToGenerate = flag.String("text", "", "Runs the generator once, generating the first 100 frames of the animation, and a sheet of 5 frames. Outputs into a folder named \"out\" if it exists. Stops the webserver from starting.")
	outputFormat   = flag.String("format", "gif", "Format of the output. Can be \"gif\" or \"neos\"")
	burningSpeed   = flag.String("speed", "fast", "Speed of the animation. Can be \"slow\" or \"fast\"")
	outputLocation = flag.String("output", "", "Location to output the generated files to. Defaults to \"out\"")
	fontPath       = flag.String("font", "", "Path to the font to use. Use semicolon to separate multiple paths.")
)

func init() {
	flag.Parse()
	godotenv.Load()
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

	fontChain := strings.Split(*fontPath, ";")

	speed := burningtext.SpeedFast
	if (*burningSpeed) == "slow" {
		speed = burningtext.SpeedSlow
	}

	settings := burningtext.BurningTextSettings{
		Text:      *textToGenerate,
		Speed:     speed,
		FontChain: fontChain,
	}

	outputBuffer := bytes.Buffer{}
	outputWriter := bufio.NewWriter(&outputBuffer)

	switch *outputFormat {
	case "gif":
		pal := palletising.LoadJSONPalette("palette.json")
		frames := burningtext.GenerateAnimatedFrames(settings)
		delaySlice := make([]int, len(frames))
		disposal := make([]byte, len(frames))
		palettedFrames := make([]*image.Paletted, len(frames))
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
			Image: palettedFrames,
			Delay: delaySlice,
			Disposal: disposal,
		})

	case "neos":
		image := burningtext.GenerateNeosSpritesheet(settings)
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
