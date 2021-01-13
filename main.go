package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/ben-yu/go-ray/primitives"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli"
)

type Camera struct {
	Origin, LowerLeftCorner, Horizontal, Vertical, U, V, W primitives.Vector
	LensRadius                                             float64
	StartTime                                              float64
	EndTime                                                float64
}

func DefaultCamera(lookFrom primitives.Vector,
	lookAt primitives.Vector,
	vUP primitives.Vector,
	vFOV float64,
	aspect float64,
	aperture float64,
	focusDist float64,
	startTime float64,
	endTime float64) Camera {

	lensRadius := aperture / 2.0
	theta := vFOV * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2)
	halfWidth := aspect * halfHeight

	w := lookFrom.Sub(lookAt).Unit()
	u := vUP.Cross(w).Unit()
	v := w.Cross(u)

	return Camera{
		Origin:          lookFrom,
		LowerLeftCorner: lookFrom.Sub(u.ScalarMul(halfWidth)).Sub(v.ScalarMul(halfHeight)).Sub(w),
		Horizontal:      u.ScalarMul(halfWidth * 2),
		Vertical:        v.ScalarMul(halfHeight * 2),
		U:               u,
		V:               v,
		W:               w,
		LensRadius:      lensRadius,
		StartTime:       startTime,
		EndTime:         endTime,
	}
}

func (c Camera) GetRay(u float64, v float64) primitives.Ray {
	return primitives.Ray{
		A:    c.Origin,
		B:    c.LowerLeftCorner.Add(c.Horizontal.ScalarMul(u)).Add(c.Vertical.ScalarMul(v)).Sub(c.Origin),
		Time: (rand.Float64() * (c.EndTime - c.StartTime)) + c.StartTime,
	}
}

func Color(r primitives.Ray, world primitives.Hitable, depth int) primitives.Vector {
	var rec primitives.HitRecord
	if world.Hit(r, 0.001, math.MaxFloat64, &rec) {
		var scattered primitives.Ray
		var attenuation primitives.Vector
		if depth < 50 && (&rec).Mat.Scatter(r, &rec, &attenuation, &scattered) {
			return attenuation.Mul(Color(scattered, world, depth+1))
		} else {
			return primitives.Vector{0.0, 0.0, 0.0}
		}
	}
	unitDirection := r.Direction().Unit()
	t := 0.5 * (unitDirection.Y() + 1.0)
	return primitives.Vector{1.0, 1.0, 1.0}.ScalarMul(1.0 - t).Add(primitives.Vector{0.5, 0.7, 1.0}).ScalarMul(t)
}

func RandomScene() primitives.HitableList {
	const n = 484
	var list []primitives.Hitable
	list = make([]primitives.Hitable, n+1, n+1)

	list[0] = primitives.Sphere{
		1000.0,
		primitives.Vector{0.0, -1000.0, 0.0},
		primitives.Lambertian{primitives.Vector{0.5, 0.5, 0.5}}}

	var i = 1
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand.Float64()
			center := primitives.Vector{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}
			if center.Sub(primitives.Vector{4.0, 0.2, 0.0}).Length() > 0.9 {
				if chooseMat < 0.8 {
					center2 := center.Add(primitives.Vector{0.0, rand.Float64() * 0.5, 0.0})
					// Diffuse
					list[i] = primitives.MovingSphere{
						0.2,
						0.0,
						1.0,
						center,
						center2,
						primitives.Lambertian{primitives.Vector{rand.Float64() * rand.Float64(), rand.Float64() * rand.Float64(), rand.Float64() * rand.Float64()}}}
				} else if chooseMat < 0.95 {
					// Metal
					list[i] = primitives.Sphere{
						0.2,
						center,
						primitives.Metal{primitives.Vector{0.5 * (1 + rand.Float64()), 0.5 * (1 + rand.Float64()), 0.5 * (1 + rand.Float64())}, 0.5 * rand.Float64()},
					}
				} else {
					// Glass
					list[i] = primitives.Sphere{
						0.2,
						center,
						primitives.Dielectric{1.5},
					}
				}
				i += 1
			}
		}
	}

	list[i] = primitives.Sphere{
		1.0,
		primitives.Vector{-4.0, 1.0, 0.0},
		primitives.Lambertian{primitives.Vector{0.4, 0.2, 0.1}},
	}
	i += 1
	list[i] = primitives.Sphere{
		1.0,
		primitives.Vector{4.0, 1.0, 0.0},
		primitives.Metal{primitives.Vector{0.7, 0.6, 0.5}, 0.0},
	}
	i += 1
	list[i] = primitives.Sphere{
		1.0,
		primitives.Vector{0.0, 1.0, 0.0},
		primitives.Dielectric{1.5},
	}
	return primitives.HitableList{list}
}

type RenderParams struct {
	Width      int
	Height     int
	Samples    int
	OutputFile string
}

func main() {
	var outputFile string
	var width, height, samples int
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "output",
				Value:       "output/output.png",
				Usage:       "Path of the output image - only supports png at the moment",
				Destination: &outputFile,
			},
			&cli.IntFlag{
				Name:        "width",
				Value:       400,
				Usage:       "Width of the image in pixels",
				Destination: &width,
			},
			&cli.IntFlag{
				Name:        "height",
				Value:       225,
				Usage:       "Height of the image in pixels",
				Destination: &height,
			},
			&cli.IntFlag{
				Name:        "samples",
				Value:       100,
				Usage:       "Number of samples to take per pixel",
				Destination: &samples,
			},
		},
		Action: func(c *cli.Context) error {

			params := RenderParams{}
			params.OutputFile = outputFile
			params.Width = width
			params.Height = height
			params.Samples = samples
			render(params)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func render(params RenderParams) {
	var width, height, numOfSamples = params.Width, params.Height, params.Samples

	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	lookFrom := primitives.Vector{13.0, 2.0, 3.0}
	lookAt := primitives.Vector{0.0, 0.0, 0.0}

	camera := DefaultCamera(
		lookFrom,
		lookAt,
		primitives.Vector{0.0, 1.0, 0.0},
		20,
		float64(width)/float64(height),
		0.1,
		10.00,
		0.0,
		1.0,
	)

	world := RandomScene()

	bar := progressbar.Default(int64(height * width))

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			col := primitives.Vector{0.0, 0.0, 0.0}
			for s := 0; s < numOfSamples; s++ {
				u := (float64(x) + rand.Float64()) / float64(width)
				v := (float64(y) + rand.Float64()) / float64(height)
				r := camera.GetRay(u, v)
				col = col.Add(Color(r, world, 0))
			}
			col = col.ScalarDiv(float64(numOfSamples))
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
			bar.Add(1)
		}
	}

	f, err := os.Create(params.OutputFile)
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
