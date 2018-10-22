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

func (m Metal) Scatter(rayIn Ray, h *HitRecord, attenuation *Vector, scattered *Ray) bool {
    reflected := Reflect(rayIn.Direction().Unit(), h.Normal)
    *scattered = Ray{h.P, reflected.Add(RandomInUnitSphere().ScalarMul(m.Fuzz))}
    *attenuation = m.Albedo
    return scattered.Direction().Dot(h.Normal) > 0
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
