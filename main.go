package main

import (
	"flag"
	"fmt"
	"image/png"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type CoolTextResponse struct {
	LogoID         int    `json:"logoID"`
	NewID          int    `json:"newID"`
	RenderLocation string `json:"renderLocation"`
	IsAnimated     bool   `json:"isAnimated"`
}

func generateFlamingText(w http.ResponseWriter, req *http.Request) {
	text := req.URL.Query().Get("text")
	
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

	img := generateNeosSpritesheet(text)

	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, img)
}

var textToGenerate = flag.String("text", "", "Runs the generator once, generating the first 100 frames of the animation, and a sheet of 5 frames. Outputs into a folder named \"out\" if it exists. Stops the webserver from starting.")

var fontDirectory string
var defaultFont string

func main() {
	godotenv.Load()

	flag.Parse()

	fontDirectory = os.Getenv("FONT_DIRECTORY")
	defaultFont = os.Getenv("DEFAULT_FONT")

	if *textToGenerate != "" {
		generateSingleText(*textToGenerate)
		return
	}

	http.HandleFunc("/", generateFlamingText)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = ":8090"
	}

	http.ListenAndServe(port, nil)
}

