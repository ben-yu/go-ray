package main

import (
    "image"
    "image/png"
    "image/color"
	"log"
    "os"
)

func main() {
    const width, height = 200, 100

    img := image.NewNRGBA(image.Rect(0, 0, width, height))

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            r, g, b := float32(x) / float32(width), float32(y) / float32(height), 0.2
            img.Set(x, height - y, color.NRGBA{
                R: uint8(r * 255.9),
                G: uint8(g * 255.9),
                B: uint8(b * 255.9),
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
