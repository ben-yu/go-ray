package primitives

import "math"

type HitRecord struct {
    t float64
    p, Normal Vector
}

type Hitable interface {
    Hit(r Ray, tMin float64, tMax float64, h *HitRecord) bool
}

type Sphere struct {
    Radius float64
    Center Vector
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
            h.p = r.PointAtParameter(h.t)
            h.Normal = (h.p.Sub(s.Center)).ScalarDiv(s.Radius)
            return true
        }
        temp = (-b + math.Sqrt((b * b) - (a * c))) / a
        if temp < tMax && temp > tMin {
            h.t = temp
            h.p = r.PointAtParameter(h.t)
            h.Normal = (h.p.Sub(s.Center)).ScalarDiv(s.Radius)
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
