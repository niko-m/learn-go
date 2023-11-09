package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	lissajous(w)
}

func lissajous(out io.Writer) {
	var palette = []color.Color{
		color.Black,
		color.White,
	}

	const (
		bgIndex      = 0
		primaryIndex = 1
	)

	const (
		cycles  = 5     // кол-во полных колебаний Х
		res     = 0.001 // угловое разрешение
		size    = 100   // канва изображения [size..+size]
		nframes = 64    // кол-во кадров анимации
		delay   = 8     // задержка между кадрами (единица - 10мс)
	)

	freq := rand.Float64() * 3.0 // относительная частота колебаний Y
	phase := 0.0                 // разность фаз
	anim := gif.GIF{LoopCount: nframes}

	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)

		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), primaryIndex)
		}

		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	gif.EncodeAll(out, &anim) // игнор ошибок
}
