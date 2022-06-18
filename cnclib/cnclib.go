package cnclib

import (
//        "fmt"
    "math"
    "cnc-tools-go/line2d"
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

func PointInPoly(x, y float64, ls orb.LineString) bool {
    plen := len(ls)
    j := plen - 1
    c := false
    for i := 0; i < len(ls); i++ {
        spt := ls[j]
        ept := ls[i]
        if (x == ept[0]) && (y == ept[1]) {
            return true
        }
        if (ept[1] > y) != (spt[1] > y) {
            dx := spt[0] - ept[0]
            dy := spt[1] - ept[1]
            slope := ((x - ept[0]) * dy) - ((y - ept[1]) * dx)
            if slope == 0 {
                return true
            }
            if (slope < 0) != (spt[1] < ept[1]) {
                c = !c
            }
        }
        j = i
    }
    return c
}

type PolyFillRaster struct {
    Sx uint32
    Sy uint32
    Conv float64
    Raster []uint8
}

func polyfillPoint(rst *PolyFillRaster, x, y uint32) orb.Point {
    return orb.Point{zify(float64(x) * rst.Conv),zify(float64(y) * rst.Conv)}
}

func polyfillCanPoint(rst *PolyFillRaster, sx, sy uint32, yincr int8) bool {
    if (yincr < 0) && (sy == 0) {
        return false
    }
    if (yincr > 0) && (sy == (rst.Sy - 1)) {
        return false
    }
    if sx == (rst.Sx - 1) {
        return false
    }
    if rst.Raster[(sy * rst.Sx) + sx] != 1 {
        return false
    }
    return true
}

func polyfillTracePath(rst *PolyFillRaster, sx, sy uint32) orb.LineString {
    var yincr bool = true
    path := make(orb.LineString, 0)
    path = append(path, polyfillPoint(rst, sx, sy))
    rst.Raster[(sy * rst.Sx) + sx] = 2
    for {
        var canVert bool = true
        if yincr {
            if ((sy + 1) == rst.Sy) || (rst.Raster[((sy + 1) * rst.Sx) + sx] != 1) {
                canVert = false
            }
        } else {
            if (sy == 0) || (rst.Raster[((sy - 1) * rst.Sx) + sx] != 1) {
                canVert = false
            }
        }
        if !canVert {
            if ((sx + 1) == rst.Sx) || (rst.Raster[(sy * rst.Sx) + (sx + 1)] != 1) {
                path = append(path, polyfillPoint(rst, sx, sy))
                break
            }
            sx = sx + 1
            yincr = !yincr
        } else {
            if yincr {
                sy = sy + 1
            } else {
                sy = sy - 1
            }
        }
        path = append(path, polyfillPoint(rst, sx, sy))
        rst.Raster[(sy * rst.Sx) + sx] = 2
    }
    return path
}

func polyfillFindPath(rst *PolyFillRaster) orb.LineString {
    var y, x uint32
    for y = 0; y < rst.Sy; y++ {
        for x = 0; x < rst.Sx; x++ {
            if rst.Raster[(y * rst.Sx) + x] == 1 {
                return polyfillTracePath(rst, x, y)
            }
        }
    }
    return orb.LineString{}
}

func PolyFill(ls orb.LineString, toolrad float64) orb.MultiLineString {
    min, max := PolygonBounds(ls)
    var line_sep = (toolrad * 2) * 0.9;
    rxdim := uint32(math.Round((max[0] - min[0]) / line_sep))
    rydim := uint32(math.Round((max[1] - min[1]) / line_sep))
    rst := new(PolyFillRaster)
    rst.Sx = rxdim
    rst.Sy = rydim
    rst.Conv = line_sep
    rst.Raster = make([]uint8, rxdim * rydim)
    var y, x uint32
    for y = 0; y < rst.Sy; y++ {
        for x = 0; x < rst.Sx; x++ {
            if PointInPoly(float64(x) * line_sep, float64(y) * line_sep, ls) {
                rst.Raster[(y * rst.Sx) + x] = 1
            } else {
                rst.Raster[(y * rst.Sx) + x] = 0
            }
        }
    }
    paths := make(orb.MultiLineString, 0)
    for {
        path := polyfillFindPath(rst)
        if len(path) == 0 {
            break
        }
        paths = append(paths, path)
    }
/*
    for y = 0; y < rst.Sy; y++ {
        for x = 0; x < rst.Sx; x++ {
            if rst.Raster[(y * rst.Sx) + x] == 1 {
                fmt.Print("#")
            } else if rst.Raster[(y * rst.Sx) + x] == 2 {
                fmt.Print("@")
            } else {
                fmt.Print(".")
            }
        }
        fmt.Print("\n")
    }
*/
    // 1) Create a 2d array of resolution 90% of the tool diameter
    // 2) Scan X,Y at this 90% of the tool diameter
    // 3) Call PointInPoly and if true set value to 1 in 2d array
    // 4) Now have rasterized grid where all 1s are inside the poly
    // 5) Scan 2d array and search for first nr 1
    // 6) From that point walk down until next point is not 1
    // 7) Check right and if 1 walk up from there until next point is not 1
    // 8) Each visit mark point as 2
    // 9) When finished append that linestring
    // 10) Repeat from 4) until no more 1s are found
    return paths
}

func LineString2PointLines(ls orb.LineString) []line2d.PointLine {
    tpl := make([]line2d.PointLine, 0, len(ls) - 1)
    for i := 1; i < len(ls); i++ {
        tpl = append(tpl, line2d.PointLine{ls[i-1][0],ls[i-1][1],ls[i][0],ls[i][1]})
    }
    return tpl
}

func BoundingBox(ls orb.LineString) orb.LineString {
    var bb = make(orb.LineString, 5)
    min, max := PolygonBounds(ls)
    bb[0] = orb.Point{min[0], min[1]}
    bb[1] = orb.Point{max[0], min[1]}
    bb[2] = orb.Point{max[0], max[1]}
    bb[3] = orb.Point{min[0], max[1]}
    bb[4] = orb.Point{min[0], min[1]}
    return bb
}

func PolygonBounds(ls orb.LineString) (orb.Point, orb.Point) {
    var minx, miny, maxx, maxy float64 = math.MaxFloat64, math.MaxFloat64, 0, 0
    for i := 0; i < len(ls); i++ {
        if ls[i][0] > maxx {
            maxx = ls[i][0]
        }
        if ls[i][0] < minx {
            minx = ls[i][0]
        }
        if ls[i][1] > maxy {
            maxy = ls[i][1]
        }
        if ls[i][1] < miny {
            miny = ls[i][1]
        }
    }
    return orb.Point{minx, miny}, orb.Point{maxx, maxy}
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
