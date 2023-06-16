package main

import (
	"flag"
	"flamingTextWebserver/burningtext"
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

	img := burningtext.GenerateNeosSpritesheet(text)

	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, img)
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

