package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"flamingTextWebserver/burningtext"
	"flamingTextWebserver/palletising"

	"github.com/joho/godotenv"
	"github.com/lucasb-eyer/go-colorful"
)

type CoolTextResponse struct {
	LogoID         int    `json:"logoID"`
	NewID          int    `json:"newID"`
	RenderLocation string `json:"renderLocation"`
	IsAnimated     bool   `json:"isAnimated"`
}

// textToGenerate = flag.String("text", "", "Runs the generator once, generating the first 100 frames of the animation, and a sheet of 5 frames. Outputs into a folder named \"out\" if it exists. Stops the webserver from starting.")
// outputFormat   = flag.String("format", "gif", "Format of the output. Can be \"gif\" or \"neos\"")
// burningSpeed   = flag.String("speed", "fast", "Speed of the animation. Can be \"slow\" or \"fast\"")
// outputLocation = flag.String("output", "", "Location to output the generated files to. Defaults to \"out\"")
// fontPath       = flag.String("font", "", "Path to the font to use. Use semicolon to separate multiple paths.")

func generateFlamingText(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	text := query.Get("text")
	if text == "" {
		w.WriteHeader(400)
		fmt.Fprintf(w, "No text provided")
		return
	}
	if len(text) > 50 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Text too long")
		return
	}

	format := strings.ToLower(query.Get("format"))
	if format == "" {
		format = "neos"
	}
	if format != "gif" && format != "neos" {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid format")
		return
	}

	speedString := strings.ToLower(query.Get("speed"))
	speed := burningtext.SpeedFast
	switch speedString {
	case "slow":
		speed = burningtext.SpeedSlow
	case "fast":
	case "":
		speed = burningtext.SpeedFast
	default:
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid speed")
		return
	}

	font := strings.ToLower(query.Get("font"))
	if font == "" {
		font = os.Getenv("DEFAULT_FONT_CHAIN")
	}
	if !burningtext.IsFontChainValid(font) {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid font")
		return
	}

	var flameColor color.Color = nil
	flameColorString := strings.ToLower(query.Get("flamecolor"))
	if flameColorString != "" {
		var err error
		flameColor, err = colorful.Hex("#" + flameColorString)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Invalid flame color")
			return
		}
	}

	var textColor color.Color = nil
	textColorString := strings.ToLower(query.Get("textcolor"))
	if textColorString != "" {
		var err error
		textColor, err = colorful.Hex("#" + textColorString)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Invalid text color")
			return
		}
	}

	randomColor := false
	randomColorString := strings.ToLower(query.Get("randomcolor"))
	if randomColorString != "" {
		if randomColorString == "true" {
			randomColor = true
		}
	}


	settings := burningtext.BurningTextOptions{
		Text:      text,
		Speed:     speed,
		FontChain: font,

		FlameColor:  flameColor,
		TextColor:   textColor,
		RandomColor: randomColor,
	}

	switch format {
	case "gif":
		generateGif(w, settings)
	case "neos":
		generateNeos(w, settings)
	default:
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid format")
		return
	}
}

func generateGif(w http.ResponseWriter, settings burningtext.BurningTextOptions) {
	burningText := burningtext.New(&settings)
	frames := burningText.GenerateBurningTextImages()
	delaySlice := make([]int, len(frames))
	disposal := make([]byte, len(frames))
	palettedFrames := make([]*image.Paletted, len(frames))
	pal := burningText.GeneratePalette()
	for i := 0; i < len(frames); i++ {
		if settings.Speed == burningtext.SpeedFast {
			delaySlice[i] = 5
		} else {
			delaySlice[i] = 2
		}

		disposal[i] = gif.DisposalPrevious
		palettedFrames[i] = palletising.Palletise(frames[i], pal)
	}

	w.Header().Set("Content-Type", "image/gif")

	gif.EncodeAll(w, &gif.GIF{
		Image:    palettedFrames,
		Delay:    delaySlice,
		Disposal: disposal,
	})
}

func generateNeos(w http.ResponseWriter, options burningtext.BurningTextOptions) {
	burningText := burningtext.New(&options)
	frames := burningText.GenerateBurningTextImages()
	image := burningtext.StackImage(frames)
	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, image)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	godotenv.Load()

	flag.Parse()

	http.HandleFunc("/", generateFlamingText)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = ":8090"
	}

	http.ListenAndServe(port, nil)
}
