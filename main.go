package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"net/http"
	"os"
	"strings"

	"flamingTextWebserver/burningtext"
	"flamingTextWebserver/palletising"

	"github.com/joho/godotenv"
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
		font = "agate"
	}
	if !burningtext.IsFontChainValid(font) {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid font")
		return
	}

	settings := burningtext.BurningTextSettings{
		Text:      text,
		Speed:     speed,
		FontChain: font,
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

func generateGif(w http.ResponseWriter, settings burningtext.BurningTextSettings) {
	frames := burningtext.GenerateAnimatedFrames(settings)
	delaySlice := make([]int, len(frames))
	disposal := make([]byte, len(frames))
	palettedFrames := make([]*image.Paletted, len(frames))
	for i := 0; i < len(frames); i++ {
		if settings.Speed == burningtext.SpeedFast {
			delaySlice[i] = 5
		} else {
			delaySlice[i] = 2
		}

		disposal[i] = gif.DisposalPrevious
		palettedFrames[i] = palletising.Palletise(frames[i])
	}

	w.Header().Set("Content-Type", "image/gif")

	gif.EncodeAll(w, &gif.GIF{
		Image:    palettedFrames,
		Delay:    delaySlice,
		Disposal: disposal,
	})
}

func generateNeos(w http.ResponseWriter, settings burningtext.BurningTextSettings) {
	image := burningtext.GenerateNeosSpritesheet(settings)
	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, image)
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
