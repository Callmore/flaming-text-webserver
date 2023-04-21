package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"net/http"
	"net/url"
)

type CoolTextResponse struct {
	LogoID         int    `json:"logoID"`
	NewID          int    `json:"newID"`
	RenderLocation string `json:"renderLocation"`
	IsAnimated     bool   `json:"isAnimated"`
}

func getFlamingText(w http.ResponseWriter, req *http.Request) {
	text := req.URL.Query().Get("text")

	fourm := url.Values{}
	fourm.Add("LogoID", "4")
	fourm.Add("Text", text)
	fourm.Add("FontSize", "70")
	fourm.Add("Color1_color", "#FF0000")
	fourm.Add("Integer1", "15")
	fourm.Add("Boolean1", "on")
	fourm.Add("Integer9", "0")
	fourm.Add("Integer13", "on")
	fourm.Add("Integer12", "on")
	fourm.Add("BackgroundColor_color", "#FFFFFF")

	resp, err := http.PostForm("https://cooltext.com/PostChange", fourm)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	buf := bytes.Buffer{}
	buf.ReadFrom(resp.Body)

	var jsonResp CoolTextResponse
	err = json.Unmarshal(buf.Bytes(), &jsonResp)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// Get the GIF
	resp, err = http.Get(jsonResp.RenderLocation)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	buf = bytes.Buffer{}
	buf.ReadFrom(resp.Body)

	img, err := stackImage(buf.Bytes())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, img)
}

func stackImage(gifImageData []byte) (image.Image, error) {
	img, err := gif.DecodeAll(bytes.NewReader(gifImageData))
	if err != nil {
		return nil, err
	}

	resultImage := image.NewRGBA(image.Rect(0, 0, img.Image[0].Rect.Max.X, img.Image[0].Rect.Max.Y*5))

	for i, img := range img.Image {
		draw.Draw(resultImage, image.Rect(0, img.Rect.Max.Y*i, img.Rect.Max.X, img.Rect.Max.Y*(i+1)), img, image.Point{0, 0}, draw.Src)
	}

	return resultImage, nil
}

func main() {
	http.HandleFunc("/getFlamingText", getFlamingText)

	http.ListenAndServe(":8090", nil)
}
