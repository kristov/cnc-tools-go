package main

import (
    "os"
    "math"
    "fmt"
    "bufio"
    "github.com/Succo/wkttoorb"
    "github.com/paulmach/orb"
    "github.com/paulmach/orb/encoding/wkt"
)

//type MultiPointLine struct {
//    sx int32
//    sy int32
//}

func main() {
    var things []interface{}
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        geo, err := wkttoorb.Scan(scanner.Text())
        if err != nil {
            panic(err)
        }
        things = append(things, geo)
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
    }

    for i := 0; i < len(things); i++ {
        switch t := things[i].(type) {
            case orb.Polygon:
                doPolygon(t)
            case orb.LineString:
                doLineString(t)
            default:
                fmt.Printf("skipping object of unknown type %T\n", t)
        }
    }
}

func doPolygon(poly orb.Polygon) {
    for i := 0; i < len(poly); i++ {
        doLineString(orb.LineString(poly[i]))
    }
}

func doLineString(ls orb.LineString) {
    if len(ls) < 2 {
        fmt.Println("this linestring does not have enough points")
        return
    }
    tpn := make([][2][2]float64, 0, len(ls) - 1)
    for i := 1; i < len(ls); i++ {
        sx := ls[i-1][0]
        sy := ls[i-1][1]
        ex := ls[i][0]
        ey := ls[i][1]
        dx := ex - sx
        dy := ey - sy
        angle := math.Atan(dy / dx)
        if dx < 0 {
            angle = angle + math.Pi
        }
        nangle := angle - (math.Pi / 2)
        nx := math.Cos(nangle) * 2
        ny := math.Sin(nangle) * 2
        tpn = append(tpn, [2][2]float64{{zify(sx+nx),zify(sy+ny)},{zify(ex+nx),zify(ey+ny)}})
    }
    var fx float64 = tpn[0][1][0]
    var fy float64 = tpn[0][1][1]
    fin := make(orb.LineString, 0, len(tpn) + 1)
    fin = append(fin, orb.Point{tpn[0][0][0],tpn[0][0][1]})
    for i := 1; i < len(tpn); i++ {
        sxa := tpn[i-1][0][0]
        sya := tpn[i-1][0][1]
        exa := tpn[i-1][1][0]
        eya := tpn[i-1][1][1]
        sxb := tpn[i][0][0]
        syb := tpn[i][0][1]
        exb := tpn[i][1][0]
        eyb := tpn[i][1][1]
        if sxa == exa {
            slb := (eyb - syb) / (exb - sxb)
            yib := syb - slb * sxb
            y := slb * exa + yib
            fin = append(fin, orb.Point{exa, y})
            fx = exb
            fy = eyb
            continue
        }
        if sxb == exb {
            sla := (eya - sya) / (exa - sxa)
            yia := sya - sla * sxa
            y := sla * sxb + yia
            fin = append(fin, orb.Point{sxb, y})
            fx = exb
            fy = eyb
            continue
        }
        sla := (eya - sya) / (exa - sxa)
        yia := sya - sla * sxa
        slb := (eyb - syb) / (exb - sxb)
        yib := syb - slb * sxb
        x := (yib - yia) / (sla - slb)
        y := sla * x + yia
        fin = append(fin, orb.Point{x, y})
        fx = exb
        fy = eyb
    }
    fin = append(fin, orb.Point{fx, fy})
    fmt.Println(wkt.MarshalString(fin))
}

func zify(value float64) float64 {
    if value < 0.000001 && value > -0.000001 {
        return 0.0
    }
    return value
}
