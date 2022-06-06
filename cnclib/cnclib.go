package cnclib

import (
    "math"
    "github.com/paulmach/orb"
    "github.com/go-gl/mathgl/mgl32"
)

type TwoPointLine struct {
    Sx float64
    Sy float64
    Ex float64
    Ey float64
}

func ToolPath(ls orb.LineString, toolrad float64) orb.LineString {
    fin := make(orb.LineString, 0, len(ls))
    if len(ls) < 2 {
        return fin
    }
    tpn := make([]TwoPointLine, 0, len(ls) - 1)
    var ep orb.Point = ls[0]
    for i := 1; i < len(ls); i++ {
        sx := ls[i-1][0]
        sy := ls[i-1][1]
        ex := ls[i][0]
        ey := ls[i][1]
        ep = ls[i]
        dx := ex - sx
        dy := ey - sy
        angle := math.Atan(dy / dx)
        if dx < 0 {
            angle = angle + math.Pi
        }
        nangle := angle - (math.Pi / 2)
        nx := math.Cos(nangle) * toolrad
        ny := math.Sin(nangle) * toolrad
        tpn = append(tpn, TwoPointLine{zify(sx+nx),zify(sy+ny),zify(ex+nx),zify(ey+ny)})
    }
    var end TwoPointLine = tpn[0]
    fin = append(fin, orb.Point{zify(tpn[0].Sx), zify(tpn[0].Sy)})
    for i := 1; i < len(tpn); i++ {
        end = tpn[i]
        p := line_intersect_point(tpn[i-1], tpn[i])
        fin = append(fin, p)
    }
    if (ls[0][0] == ep[0]) && (ls[0][1] == ep[1]) {
        p := line_intersect_point(end, tpn[0])
        fin = append(fin, p)
        fin[0][0] = p[0]
        fin[0][1] = p[1]
    } else {
        fin = append(fin, orb.Point{zify(end.Ex), zify(end.Ey)})
    }
    return fin
}

func LineStringToTwoPointLines(ls orb.LineString) []TwoPointLine {
    tpl := make([]TwoPointLine, 0, len(ls) - 1)
    for i := 1; i < len(ls); i++ {
        tpl = append(tpl, TwoPointLine{ls[i-1][0],ls[i-1][1],ls[i][0],ls[i][1]})
    }
    return tpl
}

func line_intersect_point(a TwoPointLine, b TwoPointLine) orb.Point {
    sxa := a.Sx
    sya := a.Sy
    exa := a.Ex
    eya := a.Ey
    sxb := b.Sx
    syb := b.Sy
    exb := b.Ex
    eyb := b.Ey
    if sxa == exa {
        slb := (eyb - syb) / (exb - sxb)
        yib := syb - slb * sxb
        y := slb * exa + yib
        return orb.Point{zify(exa), zify(y)}
    }
    if sxb == exb {
        sla := (eya - sya) / (exa - sxa)
        yia := sya - sla * sxa
        y := sla * sxb + yia
        return orb.Point{zify(sxb), zify(y)}
    }
    sla := (eya - sya) / (exa - sxa)
    yia := sya - sla * sxa
    slb := (eyb - syb) / (exb - sxb)
    yib := syb - slb * sxb
    x := (yib - yia) / (sla - slb)
    y := sla * x + yia
    return orb.Point{zify(x), zify(y)}
}

func Translate(ls orb.LineString, dx, dy float64) orb.LineString {
    fin := make(orb.LineString, 0, len(ls))
    for i := 0; i < len(ls); i++ {
        fin = append(fin, orb.Point{zify(ls[i][0] + dx), zify(ls[i][1] + dy)})
    }
    return fin
}

func Rotate(ls orb.LineString, radians float64) orb.LineString {
    fin := make(orb.LineString, 0, len(ls))
    mat := mgl32.Rotate2D(float32(radians))
    for i := 0; i < len(ls); i++ {
        p := mat.Mul2x1(mgl32.Vec2{float32(ls[i][0]), float32(ls[i][1])})
        fin = append(fin, orb.Point{zify(float64(p[0])), zify(float64(p[1]))})
    }
    return fin
}

func MirrorY(ls orb.LineString) orb.LineString {
    fin := make(orb.LineString, 0, len(ls))
    for i := 0; i < len(ls); i++ {
        fin = append(fin, orb.Point{zify(0 - ls[i][0]), zify(ls[i][1])})
    }
    return fin
}

func MirrorX(ls orb.LineString) orb.LineString {
    fin := make(orb.LineString, 0, len(ls))
    for i := 0; i < len(ls); i++ {
        fin = append(fin, orb.Point{zify(ls[i][0]), zify(0 - ls[i][1])})
    }
    return fin
}

func Reverse(ls orb.LineString) orb.LineString {
    fin := ls.Clone()
    fin.Reverse()
    return fin
}

func zify(value float64) float64 {
    if value < 0.000001 && value > -0.000001 {
        return 0.0
    }
    return math.Round(value * 100) / 100
}

func GeometryToLineStrings(geos []orb.Geometry) []orb.LineString {
    var linestrings []orb.LineString
    for i := 0; i < len(geos); i++ {
        switch t := geos[i].(type) {
            case orb.LineString:
                linestrings = append(linestrings, orb.LineString(t))
            case orb.Polygon:
                for j := 0; j < len(t); j++ {
                    linestrings = append(linestrings, orb.LineString(t[j]))
                }
            case orb.MultiLineString:
                for j := 0; j < len(t); j++ {
                    linestrings = append(linestrings, orb.LineString(t[j]))
                }
        }
    }
    return linestrings
}
