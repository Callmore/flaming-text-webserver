package main

import (
	"flag"
	"flamingTextWebserver/burningtext"

	"github.com/joho/godotenv"
)

var textToGenerate = flag.String("text", "", "Runs the generator once, generating the first 100 frames of the animation, and a sheet of 5 frames. Outputs into a folder named \"out\" if it exists. Stops the webserver from starting.")

func init() {
	flag.Parse()
	godotenv.Load()
}

func main() {
	if *textToGenerate != "" {
		burningtext.GenerateSingleText(*textToGenerate)
		return
	}
}
