package main

import (
	"go-ray/primitives"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func Color(r primitives.Ray) primitives.Vector {
	unitDirection := r.Direction().Unit()
	t := 0.5 * (unitDirection.Y() + 1.0)
	return primitives.Vector{1.0, 1.0, 1.0}.ScalarMul(1.0 - t).Add(primitives.Vector{0.5, 0.7, 1.0}).ScalarMul(t)
}

func main() {
	const width, height = 200, 100

	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	lowerLeftCorner := primitives.Vector{-2.0, -1.0, -1.0}
	horizontal := primitives.Vector{4.0, 0.0, 0.0}
	vertical := primitives.Vector{0.0, 2.0, 0.0}
	origin := primitives.Vector{0.0, 0.0, 0.0}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			u := float64(x) / float64(width)
			v := float64(y) / float64(height)
			r := primitives.Ray{origin, lowerLeftCorner.Add(horizontal.ScalarMul(u)).Add(vertical.ScalarMul(v))}
			col := Color(r)

			img.Set(x, height-y, color.NRGBA{
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
