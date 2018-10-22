package primitives

import "math"

type Vector struct {
    E0, E1, E2 float64
}

func (v Vector) X() float64 { return v.E0 }
func (v Vector) Y() float64 { return v.E1 }
func (v Vector) Z() float64 { return v.E2 }
func (v Vector) R() float64 { return v.E0 }
func (v Vector) G() float64 { return v.E1 }
func (v Vector) B() float64 { return v.E2 }

func (v Vector) Add(ov Vector) Vector {
    return Vector{ v.E0 + ov.E0, v.E1 + ov.E1, v.E2 + ov.E2 }
}

func (v Vector) Sub(ov Vector) Vector {
    return Vector{ v.E0 - ov.E0, v.E1 - ov.E1, v.E2 - ov.E2 }
}

func (v Vector) ScalarMul(m float64) Vector {
    return Vector{ m * v.E0, m * v.E1, m * v.E2 }
}

func (v Vector) Mul(ov Vector) Vector {
    return Vector{ v.E0 * ov.E0, v.E1 * ov.E1, v.E2 * ov.E2 }
}

func (v Vector) ScalarDiv(m float64) Vector {
    return Vector{ v.E0 / m, v.E1 / m, v.E2 / m }
}

func (v Vector) Div(ov Vector) Vector {
    return Vector{ v.E0 / ov.E0, v.E1 / ov.E1, v.E2 / ov.E2 }
}

func (v Vector) Dot(ov Vector) float64 {
    return v.E0 * ov.E0 + v.E1 * ov.E1 + v.E2 * ov.E2
}

func (v Vector) Cross(ov Vector) Vector {
    return Vector{ v.E1 * ov.E2 - v.E2 * ov.E1,
                  -(v.E0 * ov.E2 - v.E2 * ov.E0),
                  v.E0 * ov.E1 - v.E1 * ov.E0 }
}

func (v Vector) Length() float64 {
    return math.Sqrt( (v.E0 * v.E0) + (v.E1 * v.E1) + (v.E2 * v.E2) )
}

func (v Vector) SquaredLength() float64 {
    return (v.E0 * v.E0) + (v.E1 * v.E1) + (v.E2 * v.E2)
}


func (v Vector) Unit() Vector {
    return v.ScalarDiv(v.Length())
}
