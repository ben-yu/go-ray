package main

import (
    "image"
    "image/png"
    "image/color"
    "log"
    "os"
    "go-ray/primitives"
)


func main() {
    const width, height = 200, 100

    img := image.NewNRGBA(image.Rect(0, 0, width, height))

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            col := primitives.Vector {
                E0: float64(x) / float64(width),
                E1: float64(y) / float64(height),
                E2: 0.2,
            }
            img.Set(x, height - y, color.NRGBA{
                R: uint8(col.R() * 255.9),
                G: uint8(col.G() * 255.9),
                B: uint8(col.B() * 255.9),
                A: 255,
            })
        }
    }

	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
