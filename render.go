package main

import (
	"github.com/ben-yu/go-ray/primitives"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
)

type Camera struct {
	Origin, LowerLeftCorner, Horizontal, Vertical primitives.Vector
}

func DefaultCamera() Camera {
	return Camera{
		Origin:          primitives.Vector{0.0, 0.0, 0.0},
		LowerLeftCorner: primitives.Vector{-2.0, -1.0, -1.0},
		Vertical:        primitives.Vector{0.0, 2.0, 0.0},
		Horizontal:      primitives.Vector{4.0, 0.0, 0.0},
	}
}

func (c Camera) GetRay(u float64, v float64) primitives.Ray {
	return primitives.Ray{c.Origin, c.LowerLeftCorner.Add(c.Horizontal.ScalarMul(u)).Add(c.Vertical.ScalarMul(v))}
}

func Color(r primitives.Ray, world primitives.Hitable, depth int) primitives.Vector {
	var rec primitives.HitRecord
	if world.Hit(r, 0.001, math.MaxFloat64, &rec) {
        var scattered primitives.Ray
        var attenuation primitives.Vector
        if depth < 10 && (&rec).Mat.Scatter(r, &rec, &attenuation, &scattered) {
            return attenuation.Mul(Color(scattered, world, depth + 1))
        } else {
            return primitives.Vector{0.0, 0.0, 0.0}
        }
	}
	unitDirection := r.Direction().Unit()
	t := 0.5 * (unitDirection.Y() + 1.0)
	return primitives.Vector{1.0, 1.0, 1.0}.ScalarMul(1.0 - t).Add(primitives.Vector{0.5, 0.7, 1.0}).ScalarMul(t)
}

func main() {
	const width, height, numOfSamples = 200, 100, 100

	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	camera := DefaultCamera()

	world := primitives.HitableList{
		[]primitives.Hitable{
			primitives.Sphere{
				0.5,
                primitives.Vector{0.0, 0.0, -1.0},
                primitives.Lambertian{primitives.Vector{0.8,0.3,0.3}},
			},
			primitives.Sphere{
				100.0,
                primitives.Vector{0.0, -100.5, -1.0},
                primitives.Lambertian{primitives.Vector{0.8,0.8,0.0}},
			},
			primitives.Sphere{
				0.5,
                primitives.Vector{1.0, 0.0, -1.0},
                primitives.Metal{primitives.Vector{0.8,0.6,0.2},0.3},
			},
			primitives.Sphere{
				0.5,
                primitives.Vector{-1.0, 0.0, -1.0},
                primitives.Metal{primitives.Vector{0.8,0.8,0.8},1.0},
			},

		},
	}

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			col := primitives.Vector{0.0, 0.0, 0.0}
			for s := 0; s < numOfSamples; s++ {
				u := (float64(x) + rand.Float64()) / float64(width)
				v := (float64(y) + rand.Float64()) / float64(height)
				r := camera.GetRay(u, v)
				col = col.Add(Color(r, world, 0))
			}
			col = col.ScalarDiv(numOfSamples)
            col = primitives.Vector{
                math.Sqrt(col.R()),
                math.Sqrt(col.G()),
                math.Sqrt(col.B()),
            }

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
