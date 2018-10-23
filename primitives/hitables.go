package primitives

import "math"
import "math/rand"

type Material interface {
    Scatter(rayIn Ray, h *HitRecord, attenuation *Vector, scattered *Ray) bool
}

type Lambertian struct {
    Albedo Vector
}

func RandomInUnitSphere() Vector {
    var p Vector
    for notFound := true; notFound; notFound = (p.SquaredLength() >= 1.0) {
        p = Vector{rand.Float64(), rand.Float64(), rand.Float64()}.ScalarMul(2.0).Sub(Vector{1.0,1.0,1.0})
    }
    return p
}

func RandomInUnitDisk() Vector {
    var p Vector
    for notFound := true; notFound; notFound = (p.SquaredLength() >= 1.0) {
        p = Vector{rand.Float64(), rand.Float64(), 0.0}.ScalarMul(2.0).Sub(Vector{1.0,1.0,0.0})
    }
    return p
}

// Scatters ray in a random direction
func (l Lambertian) Scatter(rayIn Ray, h *HitRecord, attenuation *Vector, scattered *Ray) bool {
    target := h.P.Add(h.Normal).Add(RandomInUnitSphere())
    *scattered = Ray{h.P, target.Sub(h.P)}
    *attenuation = l.Albedo
    return true
}

type Metal struct {
    Albedo Vector
    Fuzz float64
}

func Reflect(v Vector, n Vector) Vector {
    return v.Sub(n.ScalarMul(v.Dot(n)).ScalarMul(2.0))
}

func Schlick(cosine float64, refIdx float64) float64 {
    r0 := (1-refIdx) / (1+refIdx)
    r0 = r0 * r0
    return r0 + (1-r0) * math.Pow(1.0-cosine, 5.0)
}

// Scatters ray by reflecting from point of intersection with some randomness
// from the reflection point
func (m Metal) Scatter(rayIn Ray, h *HitRecord, attenuation *Vector, scattered *Ray) bool {
    reflected := Reflect(rayIn.Direction().Unit(), h.Normal)
    *scattered = Ray{h.P, reflected.Add(RandomInUnitSphere().ScalarMul(m.Fuzz))}
    *attenuation = m.Albedo
    return scattered.Direction().Dot(h.Normal) > 0
}

type Dielectric struct {
    RefIdx float64
}

func (d Dielectric) Scatter(rayIn Ray, h *HitRecord, attenuation *Vector, scattered *Ray) bool {
    var outwardNormal Vector
    reflected := Reflect(rayIn.Direction(), h.Normal)
    var ratio, reflectProb, cosine float64
    var refracted Vector
    *attenuation = Vector{1.0, 1.0, 1.0}
    if rayIn.Direction().Dot(h.Normal) > 0 {
        outwardNormal = h.Normal.ScalarMul(-1)
        ratio = d.RefIdx
        cosine = d.RefIdx * rayIn.Direction().Dot(h.Normal) / rayIn.Direction().Length()
    } else {
        outwardNormal = h.Normal
        ratio = 1.0 / d.RefIdx
        cosine = -rayIn.Direction().Dot(h.Normal) / rayIn.Direction().Length()
    }

    if Refract(rayIn.Direction(), outwardNormal, ratio, &refracted) {
        reflectProb = Schlick(cosine, d.RefIdx)
    } else {
        *scattered = Ray{h.P, reflected}
        reflectProb = 1.0
    }

    if rand.Float64() < reflectProb {
        *scattered = Ray{h.P, reflected}
    } else {
        *scattered = Ray{h.P, refracted}
    }
    return true
}

// Implementation of Snell's Law
func Refract(v Vector, n Vector, ratio float64, refracted *Vector) bool {
    uv := v.Unit()
    dt := uv.Dot(n)
    discriminant := 1.0 - ratio * ratio * (1 - dt * dt)
    if discriminant > 0 {
        *refracted = uv.Sub(n.ScalarMul(dt)).ScalarMul(ratio).Sub(n.ScalarMul(math.Sqrt(discriminant)))
        return true
    }

    return false
}

type HitRecord struct {
    t float64
    P, Normal Vector
    Mat Material
}

type Hitable interface {
    Hit(r Ray, tMin float64, tMax float64, h *HitRecord) bool
}

type Sphere struct {
    Radius float64
    Center Vector
    Mat Material
}

func (s Sphere) Hit(r Ray, tMin float64, tMax float64, h *HitRecord) bool {
    oc := r.Origin().Sub(s.Center)
    a := r.Direction().Dot(r.Direction())
    b := oc.Dot(r.Direction())
    c := oc.Dot(oc) - (s.Radius * s.Radius)
    discriminant := (b * b) - (a * c)
    if discriminant > 0 {
        temp := (-b - math.Sqrt((b * b) - (a * c))) / a
        if temp < tMax && temp > tMin {
            h.t = temp
            h.P = r.PointAtParameter(h.t)
            h.Normal = (h.P.Sub(s.Center)).ScalarDiv(s.Radius)
            h.Mat = s.Mat
            return true
        }
        temp = (-b + math.Sqrt((b * b) - (a * c))) / a
        if temp < tMax && temp > tMin {
            h.t = temp
            h.P = r.PointAtParameter(h.t)
            h.Normal = (h.P.Sub(s.Center)).ScalarDiv(s.Radius)
            h.Mat = s.Mat
            return true
        }
    }
    return false
}

type HitableList struct {
    List []Hitable
}

func (hitableList HitableList) Hit(r Ray, tMin float64, tMax float64, h *HitRecord) bool {
    tempRecord := h
    hitAnything := false
    closestSoFar := tMax
    for _, hitableElem := range hitableList.List {
        if hitableElem.Hit(r, tMin, closestSoFar, tempRecord) {
            hitAnything = true
            closestSoFar = tempRecord.t
            h = tempRecord
        }
    }
    return hitAnything
}
