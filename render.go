package main

import (
	"go-ray/primitives"
	"image"
	"image/color"
	"image/png"
	"log"
    "math"
	"os"
)

func HitSphere(center primitives.Vector, radius float64, r primitives.Ray) float64 {
    oc := r.Origin().Sub(center)
    a := r.Direction().Dot(r.Direction())
    b := oc.Dot(r.Direction()) * 2.0
    c := oc.Dot(oc) - (radius * radius)
    discriminant := (b * b) - (4 * a * c)
    if discriminant < 0 {
        return -1.0
    } else {
        return (-b - math.Sqrt(discriminant)) / (2.0 * a)
    }
}

func Color(r primitives.Ray) primitives.Vector {
    t_sphere := HitSphere(primitives.Vector{0.0, 0.0, -1.0}, 0.5, r)
    if t_sphere > 0.0 {
        norm := r.PointAtParameter(t_sphere).Sub(primitives.Vector{ 0.0, 0.0, -1.0}).Unit()
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

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			u := float64(x) / float64(width)
			v := float64(y) / float64(height)
			r := primitives.Ray{origin, lowerLeftCorner.Add(horizontal.ScalarMul(u)).Add(vertical.ScalarMul(v))}
			col := Color(r)


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
