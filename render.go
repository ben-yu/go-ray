package main

import (
	"github.com/ben-yu/go-ray/primitives"
	"image"
	"image/color"
	"image/png"
	"log"
    "math"
	"os"
)

func Color(r primitives.Ray, world primitives.Hitable) primitives.Vector {
    var rec primitives.HitRecord
    if world.Hit(r, 0.0, math.MaxFloat64, &rec) {
        norm := rec.Normal
        return primitives.Vector{norm.X() + 1, norm.Y() + 1, norm.Z() + 1}.ScalarMul(0.5)
    }
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

    world := primitives.HitableList{
        []primitives.Hitable {
            primitives.Sphere{
                0.5, primitives.Vector{0.0, 0.0, -1.0},
            },
            primitives.Sphere{
                100.0, primitives.Vector{0.0, -100.5, -1.0},
            },
        },
    }

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			u := float64(x) / float64(width)
			v := float64(y) / float64(height)
			r := primitives.Ray{origin, lowerLeftCorner.Add(horizontal.ScalarMul(u)).Add(vertical.ScalarMul(v))}
			col := Color(r, world)


			img.Set(x, height-y, color.NRGBA{
				R: uint8(col.R() * 255.99),
				G: uint8(col.G() * 255.99),
				B: uint8(col.B() * 255.99),
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
