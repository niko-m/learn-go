package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

type Params struct {
	Cycles  float64 // кол-во полных колебаний Х
	Res     float64 // угловое разрешение
	Size    float64 // канва изображения [size..+size]
	Nframes float64 // кол-во кадров анимации
	Delay   float64 // задержка между кадрами (единица - 10мс)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := Params{5, 0.001, 100, 64, 8}

	q := r.URL.Query()

	if vals, exists := q["cycles"]; exists {
		v, err := strconv.ParseFloat(vals[len(vals)-1], 64)
		if err == nil {
			p.Cycles = v
		}
	}

	if vals, exists := q["res"]; exists {
		v, err := strconv.ParseFloat(vals[len(vals)-1], 64)
		if err == nil {
			p.Res = v
		}
	}

	if vals, exists := q["size"]; exists {
		v, err := strconv.ParseFloat(vals[len(vals)-1], 64)
		if err == nil {
			p.Size = v
		}
	}

	if vals, exists := q["nframes"]; exists {
		v, err := strconv.ParseFloat(vals[len(vals)-1], 64)
		if err == nil {
			p.Nframes = v
		}
	}

	if vals, exists := q["delay"]; exists {
		v, err := strconv.ParseFloat(vals[len(vals)-1], 64)
		if err == nil {
			p.Delay = v
		}
	}

	fmt.Println(p)

	lissajous(w, p)
}

func lissajous(out io.Writer, p Params) {
	var palette = []color.Color{
		color.Black,
		color.White,
	}

	const (
		bgIndex      = 0
		primaryIndex = 1
	)

	freq := rand.Float64() * 3.0 // относительная частота колебаний Y
	phase := 0.0                 // разность фаз
	anim := gif.GIF{LoopCount: int(p.Nframes)}

	for i := 0; i < int(p.Nframes); i++ {
		rect := image.Rect(0, 0, int(2*p.Size+1), int(2*p.Size+1))
		img := image.NewPaletted(rect, palette)

		for t := 0.0; t < p.Cycles*2*math.Pi; t += p.Res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(int(p.Size+x*p.Size+0.5), int(p.Size+y*p.Size+0.5), primaryIndex)
		}

		phase += 0.1
		anim.Delay = append(anim.Delay, int(p.Delay))
		anim.Image = append(anim.Image, img)
	}

	gif.EncodeAll(out, &anim) // игнор ошибок
}
