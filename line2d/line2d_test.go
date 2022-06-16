package line2d

import (
    "math"
    "testing"
)

func Test45(t *testing.T) {
    L := PointLine{0.0,0.0,2.0,2.0}
    AL := PointLine2Line(L)
    if AL.Vert {
        t.Log("Incorrect vertical line flag")
        t.Fail()
    }
    if AL.Slope != 1.0 {
        t.Log("Slope of a 45 degree line should be 1.0")
        t.Fail()
    }
    if AL.Yint != 0.0 {
        t.Log("Incorrect y intercept")
        t.Fail()
    }
    lAngle := RoundP(PointLineAngle(L), 4)
    if lAngle != RoundP(math.Pi / 4, 4) {
        t.Log("Incorrect PointLineAngle")
        t.Fail()
    }
    alAngle := RoundP(LineAngle(AL), 4)
    if alAngle != RoundP(math.Pi / 4, 4) {
        t.Log("Incorrect LineAngle")
        t.Fail()
    }
}

func TestVert(t *testing.T) {
    L := PointLine{1.2,-1.0,1.2,1.0}
    AL := PointLine2Line(L)
    if !AL.Vert {
        t.Log("Vertical line flag not set")
        t.Fail()
    }
    if AL.Xint != 1.2 {
        t.Log("Incorrect x intercept")
        t.Fail()
    }
    lAngle := RoundP(PointLineAngle(L), 4)
    if lAngle != RoundP(math.Pi / 2, 4) {
        t.Log("Incorrect PointLineAngle")
        t.Fail()
    }
    alAngle := RoundP(LineAngle(AL), 4)
    if alAngle != RoundP(math.Pi / 2, 4) {
        t.Log("Incorrect LineAngle", alAngle)
        t.Fail()
    }
}

func TestHorizontal(t *testing.T) {
    L := PointLine{-1.0,2.1,1.0,2.1}
    AL := PointLine2Line(L)
    if AL.Vert {
        t.Log("Incorrect vertical line flag")
        t.Fail()
    }
    if AL.Slope != 0.0 {
        t.Log("Slope of a horizontal line should be zero")
        t.Fail()
    }
    if AL.Yint != 2.1 {
        t.Log("Incorrect y intercept")
        t.Fail()
    }
    lAngle := RoundP(PointLineAngle(L), 4)
    if lAngle != 0.0 {
        t.Log("Incorrect PointLineAngle")
        t.Fail()
    }
    alAngle := RoundP(LineAngle(AL), 4)
    if alAngle != 0.0 {
        t.Log("Incorrect LineAngle", alAngle)
        t.Fail()
    }
}

func TestYoff(t *testing.T) {
    L := PointLine{1.0,2.0,3.0,3.0}
    AL := PointLine2Line(L)
    if AL.Vert {
        t.Log("Incorrect vertical line flag")
        t.Fail()
    }
    if AL.Slope != 0.5 {
        t.Log("Slope of this line should be 0.5")
        t.Fail()
    }
    if AL.Yint != 1.5 {
        t.Log("Incorrect y intercept")
        t.Fail()
    }
    lAngle := RoundP(PointLineAngle(L), 4)
    if lAngle != 0.4636 {
        t.Log("Incorrect PointLineAngle")
        t.Fail()
    }
    alAngle := RoundP(LineAngle(AL), 4)
    if alAngle != 0.4636 {
        t.Log("Incorrect LineAngle", alAngle)
        t.Fail()
    }
}

func TestNegSlope(t *testing.T) {
    L := PointLine{2.0,-3.0,3.0,-6.0}
    AL := PointLine2Line(L)
    if AL.Vert {
        t.Log("Incorrect vertical line flag")
        t.Fail()
    }
    if AL.Slope != -3.0 {
        t.Log("Slope of this line should be -3.0")
        t.Fail()
    }
    if AL.Yint != 3.0 {
        t.Log("Incorrect y intercept")
        t.Fail()
    }
    lAngle := RoundP(PointLineAngle(L), 4)
    if lAngle != 5.0341 {
        t.Log("Incorrect PointLineAngle", lAngle)
        t.Fail()
    }
    alAngle := RoundP(LineAngle(AL), 4)
    if alAngle != -1.2490 {
        t.Log("Incorrect LineAngle", alAngle)
        t.Fail()
    }
}

func TestFunkyColdMedina(t *testing.T) {
    L := PointLine{1.0,-2.0,-3.0,-1.0}
    AL := PointLine2Line(L)
    if AL.Vert {
        t.Log("Incorrect vertical line flag")
        t.Fail()
    }
    if AL.Slope != -0.25 {
        t.Log("Slope of this line should be -0.25")
        t.Fail()
    }
    if AL.Yint != -1.75 {
        t.Log("Incorrect y intercept")
        t.Fail()
    }
    lAngle := RoundP(PointLineAngle(L), 4)
    if lAngle != 2.8966 {
        t.Log("Incorrect PointLineAngle", lAngle)
        t.Fail()
    }
    alAngle := RoundP(LineAngle(AL), 4)
    if alAngle != -0.245 {
        t.Log("Incorrect LineAngle", alAngle)
        t.Fail()
    }
}

func TestQuadrantAngles(t *testing.T) {
    qpi := math.Pi / 4
    // Quadrants are areas on a cartesian plane that represent positive or
    // negative X or Y values. Quadrant 1 is positive X and positive Y,
    // Quadrant 2 is negative X and positive Y, Quadrant 3 is negative X and
    // negative Y and Quadrant 4 is positive X and negative Y.
    l1 := PointLine{0.0,0.0,1.0,1.0}   // Quadrant 1
    l2 := PointLine{0.0,0.0,-1.0,1.0}  // Quadrant 2
    l3 := PointLine{0.0,0.0,-1.0,-1.0} // Quadrant 3
    l4 := PointLine{0.0,0.0,1.0,-1.0}  // Quadrant 4
    la1 := RoundP(PointLineAngle(l1), 4)
    la2 := RoundP(PointLineAngle(l2), 4)
    la3 := RoundP(PointLineAngle(l3), 4)
    la4 := RoundP(PointLineAngle(l4), 4)
    if la1 != RoundP(qpi, 4) {
        t.Log("Incorrect PointLineAngle for Quadrant 1", la1)
    }
    if la2 != RoundP(math.Pi - qpi, 4) {
        t.Log("Incorrect PointLineAngle for Quadrant 2", la2)
    }
    if la3 != RoundP(math.Pi + qpi, 4) {
        t.Log("Incorrect PointLineAngle for Quadrant 3", la3)
    }
    if la4 != RoundP((math.Pi * 2) - qpi, 4) {
        t.Log("Incorrect PointLineAngle for Quadrant 4", la4, RoundP((math.Pi * 2) - qpi, 4))
    }
    laa1 := RoundP(LineAngle(PointLine2Line(l1)), 4)
    laa2 := RoundP(LineAngle(PointLine2Line(l2)), 4)
    laa3 := RoundP(LineAngle(PointLine2Line(l3)), 4)
    laa4 := RoundP(LineAngle(PointLine2Line(l4)), 4)
    if laa1 != RoundP(qpi, 4) {
        t.Log("Incorrect LineAngle for Quadrant 1", laa1)
    }
    if laa2 != RoundP(0 - qpi, 4) {
        t.Log("Incorrect LineAngle for Quadrant 2", laa2)
    }
    if laa3 != RoundP(qpi, 4) {
        t.Log("Incorrect LineAngle for Quadrant 3", laa3)
    }
    if laa4 != RoundP(0 - qpi, 4) {
        t.Log("Incorrect LineAngle for Quadrant 4", laa4)
    }
}

func TestLineIntersect(t *testing.T) {
    // A horizontal line and an angled line
    l1 := PointLine2Line(PointLine{1.0,1.0,5.0,1.0})
    l2 := PointLine2Line(PointLine{1.0,0.0,3.0,2.0})
    p1, int1 := LineIntersect(l1, l2)
    if !int1 {
        t.Log("Lines should intersect")
    }
    if !(p1[0] == 2.0 && p1[1] == 1.0) {
        t.Log("Intersect point wrong")
    }
    p2, int2 := LineIntersect(l2, l1)
    if !int2 {
        t.Log("Lines should intersect")
    }
    if !(p2[0] == 2.0 && p2[1] == 1.0) {
        t.Log("Intersect point wrong")
    }

    // A horizontal line and a vertical line
    l3 := PointLine2Line(PointLine{1.0,1.0,5.0,1.0})
    l4 := PointLine2Line(PointLine{2.0,-1.0,2.0,2.0})
    p3, int3 := LineIntersect(l3, l4)
    if !int3 {
        t.Log("Lines should intersect")
    }
    if !(p3[0] == 2.0 && p3[1] == 1.0) {
        t.Log("Intersect point wrong")
    }
    p4, int4 := LineIntersect(l4, l3)
    if !int4 {
        t.Log("Lines should intersect")
    }
    if !(p4[0] == 2.0 && p4[1] == 1.0) {
        t.Log("Intersect point wrong")
    }

    // Two horizontal lines
    l5 := PointLine2Line(PointLine{1.0,1.0,5.0,1.0})
    l6 := PointLine2Line(PointLine{1.0,2.0,5.0,2.0})
    _, int5 := LineIntersect(l5, l6)
    if int5 {
        t.Log("These lines should not intersect")
    }

    // Two vertical lines
    l7 := PointLine2Line(PointLine{2.0,-1.0,2.0,2.0})
    l8 := PointLine2Line(PointLine{3.0,-1.0,3.0,2.0})
    _, int6 := LineIntersect(l7, l8)
    if int6 {
        t.Log("These lines should not intersect")
    }

    // Two angled lines
    l9 := PointLine2Line(PointLine{1.0,1.0,-1.0,5.0})
    l10 := PointLine2Line(PointLine{-2.0,2.0,2.0,4.0})
    p7, int7 := LineIntersect(l9, l10)
    if !int7 {
        t.Log("Lines should intersect")
    }
    if !(p7[0] == 0.0 && p7[1] == 3.0) {
        t.Log("Intersect point wrong")
    }

    // Two lines far apart - but there is a virtual intersection
    l11 := PointLine2Line(PointLine{-4.0,-2.0,0.0,-1.0})
    l12 := PointLine2Line(PointLine{2.0,4.0,3.0,2.0})
    p9, int9 := LineIntersect(l11, l12)
    if !int9 {
        t.Log("Lines should intersect")
    }
    if !(p9[0] == 4.0 && p9[1] == 0.0) {
        t.Log("Intersect point wrong")
    }
}

func TestPointLineIntersect(t *testing.T) {
    // A horizontal line and an angled line
    l1 := PointLine{1.0,1.0,5.0,1.0}
    l2 := PointLine{1.0,0.0,3.0,2.0}
    p1, int1 := PointLineIntersect(l1, l2)
    if !int1 {
        t.Log("Lines should intersect")
    }
    if !(p1[0] == 2.0 && p1[1] == 1.0) {
        t.Log("Intersect point wrong")
    }
    p2, int2 := PointLineIntersect(l2, l1)
    if !int2 {
        t.Log("Lines should intersect")
    }
    if !(p2[0] == 2.0 && p2[1] == 1.0) {
        t.Log("Intersect point wrong")
    }

    // A horizontal line and a vertical line
    l3 := PointLine{1.0,1.0,5.0,1.0}
    l4 := PointLine{2.0,-1.0,2.0,2.0}
    p3, int3 := PointLineIntersect(l3, l4)
    if !int3 {
        t.Log("Lines should intersect")
    }
    if !(p3[0] == 2.0 && p3[1] == 1.0) {
        t.Log("Intersect point wrong")
    }
    p4, int4 := PointLineIntersect(l4, l3)
    if !int4 {
        t.Log("Lines should intersect")
    }
    if !(p4[0] == 2.0 && p4[1] == 1.0) {
        t.Log("Intersect point wrong")
    }

    // Two horizontal lines
    l5 := PointLine{1.0,1.0,5.0,1.0}
    l6 := PointLine{1.0,2.0,5.0,2.0}
    _, int5 := PointLineIntersect(l5, l6)
    if int5 {
        t.Log("These lines should not intersect")
    }

    // Two vertical lines
    l7 := PointLine{2.0,-1.0,2.0,2.0}
    l8 := PointLine{3.0,-1.0,3.0,2.0}
    _, int6 := PointLineIntersect(l7, l8)
    if int6 {
        t.Log("These lines should not intersect")
    }

    // Two angled lines
    l9 := PointLine{1.0,1.0,-1.0,5.0}
    l10 := PointLine{-2.0,2.0,2.0,4.0}
    p7, int7 := PointLineIntersect(l9, l10)
    if !int7 {
        t.Log("Lines should intersect")
    }
    if !(p7[0] == 0.0 && p7[1] == 3.0) {
        t.Log("Intersect point wrong")
    }

    // Two lines far apart - with PointLine they should not intersect
    l11 := PointLine{-4.0,-2.0,0.0,-1.0}
    l12 := PointLine{2.0,4.0,3.0,2.0}
    _, int9 := PointLineIntersect(l11, l12)
    if int9 {
        t.Log("Lines should not intersect")
    }
}

func RoundP(value float64, precision uint32) float64 {
    var ratio = math.Pow(10, float64(precision))
    return math.Round(value * ratio) / ratio
}
