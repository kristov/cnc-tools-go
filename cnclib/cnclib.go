package cnclib

import (
        "fmt"
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

type PolyFillRaster struct {
    Sx uint32
    Sy uint32
    Conv float64
    Raster []uint8
}

type PolyFillVector struct {
    Sx uint32
    Sy uint32
    Yincr bool
}

func ToolPath(ls orb.LineString, toolrad float64) orb.LineString {
    if len(ls) < 2 {
        return orb.LineString{}
    }
    tpn := make([]line2d.PointLine, len(ls) - 1)
    var ep orb.Point = ls[0]
    for i := 1; i < len(ls); i++ {
        sx := ls[i-1][0]
        sy := ls[i-1][1]
        ex := ls[i][0]
        ey := ls[i][1]
        ep = ls[i]
        pl := line2d.PointLine{sx,sy,ex,ey}
        angle := line2d.PointLineAngle(pl)
        nangle := angle - (math.Pi / 2)
        nx := math.Cos(nangle) * toolrad
        ny := math.Sin(nangle) * toolrad
        tpn[i-1] = line2d.PointLine{sx+nx,sy+ny,ex+nx,ey+ny}
    }
    var end line2d.PointLine = tpn[0]
    points := make([]line2d.Point, len(ls))
    points[0] = line2d.Point{tpn[0][0],tpn[0][1]}
    for i := 1; i < len(tpn); i++ {
        end = tpn[i]
        p, _ := line2d.LineIntersect(line2d.PointLine2Line(tpn[i-1]), line2d.PointLine2Line(tpn[i]))
        points[i] = p
    }
    if (ls[0][0] == ep[0]) && (ls[0][1] == ep[1]) {
        p, _ := line2d.LineIntersect(line2d.PointLine2Line(end), line2d.PointLine2Line(tpn[0]))
        points[len(tpn)] = p
        points[0][0] = p[0]
        points[0][1] = p[1]
    } else {
        points[len(tpn)] = line2d.Point{end[2],end[3]}
    }
    fin := make(orb.LineString, len(ls))
    for i := 0; i < len(points); i++ {
        fin[i] = orb.Point{zify(points[i][0]),zify(points[i][1])}
    }
    return fin
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

func polyfillPoint(rst *PolyFillRaster, x, y uint32) orb.Point {
    return orb.Point{zify(float64(x) * rst.Conv),zify(float64(y) * rst.Conv)}
}

func polyfillCanMove(rst *PolyFillRaster, x, y uint32, dx, dy int8) bool {
    if dx < 0 {
        if x == 0 {
            return false
        }
        x = x - 1
    }
    if dy < 0 {
        if y == 0 {
            return false
        }
        y = y - 1
    }
    if dx > 0 {
        x = x + 1
        if x == rst.Sx {
            return false
        }
    }
    if dy > 0 {
        y = y + 1
        if y == rst.Sy {
            return false
        }
    }
    if rst.Raster[(y * rst.Sx) + x] != 1 {
        return false
    }
    return true
}

func polyfillMarkAppend(rst *PolyFillRaster, ls orb.LineString, x, y uint32) orb.LineString {
    ls = append(ls, polyfillPoint(rst, x, y))
    rst.Raster[(y * rst.Sx) + x] = 2
    return ls
}

func polyfillTracePath(rst *PolyFillRaster, pv *PolyFillVector) orb.LineString {
    path := make(orb.LineString, 0)
    path = polyfillMarkAppend(rst, path, pv.Sx, pv.Sy)
    for {
        if pv.Yincr && polyfillCanMove(rst, pv.Sx, pv.Sy, 0, 1) {
            // We are moving up and we can go to the next space up
            rst.Raster[(pv.Sy * rst.Sx) + pv.Sx] = 2
            pv.Sy = pv.Sy + 1
            continue
        }
        if !pv.Yincr && polyfillCanMove(rst, pv.Sx, pv.Sy, 0, -1) {
            // We are moving down and we can go to the next space down
            rst.Raster[(pv.Sy * rst.Sx) + pv.Sx] = 2
            pv.Sy = pv.Sy - 1
            continue
        }
        // We can no longer move up or down in this column
        path = polyfillMarkAppend(rst, path, pv.Sx, pv.Sy)
        if pv.Yincr {
            // We are moving up
            if polyfillCanMove(rst, pv.Sx, pv.Sy, 1, 1) {
                // We can we move diagonally up and to the right
                pv.Sx = pv.Sx + 1
                pv.Sy = pv.Sy + 1
                pv.Yincr = false
                break
            } else if polyfillCanMove(rst, pv.Sx, pv.Sy, 1, 0) {
                // We can move directly to the right
                pv.Sx = pv.Sx + 1
                pv.Yincr = false
                path = polyfillMarkAppend(rst, path, pv.Sx, pv.Sy)
                continue
            } else if polyfillCanMove(rst, pv.Sx, pv.Sy, 1, -1) {
                // We can we move diagonally down and to the right
                pv.Sx = pv.Sx + 1
                pv.Sy = pv.Sy - 1
                pv.Yincr = false
                break
            } else {
                // We can not move anywhere, return to rescan
                break
            }
        } else {
            // If we are moving down
            if polyfillCanMove(rst, pv.Sx, pv.Sy, 1, -1) {
                // We can we move diagonally down and to the right
                pv.Sx = pv.Sx + 1
                pv.Sy = pv.Sy - 1
                pv.Yincr = true
                break
            } else if polyfillCanMove(rst, pv.Sx, pv.Sy, 1, 0) {
                // We can move directly to the right
                pv.Sx = pv.Sx + 1
                pv.Yincr = true
                path = polyfillMarkAppend(rst, path, pv.Sx, pv.Sy)
                continue
            } else if polyfillCanMove(rst, pv.Sx, pv.Sy, 1, 1) {
                // We can we move diagonally up and to the right
                pv.Sx = pv.Sx + 1
                pv.Sy = pv.Sy + 1
                pv.Yincr = true
                break
            } else {
                // We can not move anywhere, return to rescan
                break
            }
        }
    }
    return path
}

func polyfillFindNext(rst *PolyFillRaster, pv *PolyFillVector) bool {
    for x := pv.Sx; x < rst.Sx; x++ {
        if pv.Yincr {
            for y := pv.Sy; y < rst.Sy; y++ {
                if rst.Raster[(y * rst.Sx) + x] == 1 {
                    pv.Sx = x
                    pv.Sy = y
                    return true
                }
            }
        } else {
            for y := pv.Sy; y > 0; y-- {
                if rst.Raster[(y * rst.Sx) + x] == 1 {
                    pv.Sx = x
                    pv.Sy = y
                    return true
                }
            }
        }
        pv.Yincr = !pv.Yincr
    }
    return false
}

func polyfillFindPaths(rst *PolyFillRaster) orb.MultiLineString {
    paths := make(orb.MultiLineString, 0)
    vector := PolyFillVector{0,0,true}
    for {
        if !polyfillFindNext(rst, &vector) {
            break
        }
        paths = append(paths, polyfillTracePath(rst, &vector))
    }
    return paths
}

func PolyFill(ls orb.LineString, toolrad float64) orb.MultiLineString {
    min, max := PolygonBounds(ls)
    lstr := Translate(ls, 0 - min[0], 0 - min[1])
    //fmt.Printf("min: %0.2f,%0.2f, max: %0.2f,%0.2f\n", min[0], min[1], max[0], max[1])
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
            if PointInPoly((float64(x) * line_sep), (float64(y) * line_sep), lstr) {
                rst.Raster[(y * rst.Sx) + x] = 1
            } else {
                rst.Raster[(y * rst.Sx) + x] = 0
            }
        }
    }
    paths := polyfillFindPaths(rst)
    fin := make(orb.MultiLineString, len(paths))
    for i := 0; i < len(paths); i++ {
        fin[i] = Translate(paths[i], min[0], min[1])
    }

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

    return fin
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

func BoundingBoxMulti(lss []orb.LineString) {
    var min, max orb.Point = orb.Point{math.MaxFloat64,math.MaxFloat64}, orb.Point{0,0}
    for i := 1; i < len(lss); i++ {
        lmin, lmax := PolygonBounds(lss[i])
        if lmax[0] > max[0] {
            max[0] = lmax[0]
        }
        if lmin[0] < min[0] {
            min[0] = lmin[0]
        }
        if lmax[1] > max[1] {
            max[1] = lmax[1]
        }
        if lmin[1] < min[1] {
            min[1] = lmin[1]
        }
    }
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
