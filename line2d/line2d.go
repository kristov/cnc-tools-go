package line2d

import "math"

type Point [2]float64
type PointLine [4]float64

type Line struct {
    Vert bool
    Hori bool
    Slope float64
    Xint float64
    Yint float64
}

func PointLine2Line(l PointLine) Line {
    if (l[2] - l[0]) == 0 { // vertical
        return Line{true,false,0.0,l[2],0}
    }
    if (l[3] - l[1]) == 0 { // horizontal
        return Line{false,true,0.0,0,l[3]}
    }
    slope := (l[3] - l[1]) / (l[2] - l[0])
    yint := l[1] - (slope * l[0])
    return Line{false,false,slope,0,yint}
}

func PointLineAngle(l PointLine) float64 {
    var dx = l[2] - l[0]
    var dy = l[3] - l[1]
    angle := math.Atan(dy / dx)
    if dx < 0 {
        angle = math.Pi + angle
    } else if dy < 0 {
        angle = (math.Pi * 2) + angle
    }
    return angle
}

// Note: an Line can never give an angle greater than 90 or
// less than -90 because we lost the order of start and end points.
func LineAngle(al Line) float64 {
    if al.Vert {
        return math.Pi / 2
    }
    return math.Atan(al.Slope)
}

func PointInBbox(p Point, l PointLine) bool {
    var minx, maxx, miny, maxy float64
    if l[0] < l[2] {
        minx = l[0]
        maxx = l[2]
    } else {
        minx = l[2]
        maxx = l[0]
    }
    if l[1] < l[3] {
        miny = l[1]
        maxy = l[3]
    } else {
        miny = l[3]
        maxy = l[1]
    }
    if p[0] >= minx && p[0] <= maxx && p[1] >= miny && p[1] <= maxy {
        return true
    }
    return false
}

func PointLineIntersect(l1 PointLine, l2 PointLine) (Point, bool) {
    p, i := LineIntersect(PointLine2Line(l1), PointLine2Line(l2))
    if !i {
        return Point{0.0,0.0}, false
    }
    if PointInBbox(p, l1) && PointInBbox(p, l2) {
        return p, true
    }
    return Point{0.0,0.0}, false
}

func LineIntersect(l1 Line, l2 Line) (Point, bool) {
    if l1.Vert && l2.Vert {
        // Two vertical lines do not intersect
        return Point{0.0,0.0}, false
    }
    if l1.Hori && l2.Hori {
        // Two horizontal lines do not intersect
        return Point{0.0,0.0}, false
    }
    lines := [2]Line{l1,l2}
    var vi = 3
    var hi = 3
    if l1.Vert {
        vi = 0
    } else if l1.Hori {
        hi = 0
    }
    if l2.Vert {
        vi = 1
    } else if l2.Hori {
        hi = 1
    }
    if vi != 3 && hi != 3 {
        // We have a horizontal and a vertical line
        return Point{lines[vi].Xint,lines[hi].Yint}, true
    }
    var ri = 3
    if vi != 3 {
        // A vertical line and something else
        ri = 1 - vi
        var y = (lines[ri].Slope * lines[vi].Xint) + lines[ri].Yint
        return Point{lines[vi].Xint,y}, true
    }
    x := (lines[1].Yint - lines[0].Yint) / (lines[0].Slope - lines[1].Slope)
    y := (lines[0].Slope * x) + lines[0].Yint
    return Point{x,y}, true
}
