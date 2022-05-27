package cnc

import (
    "math"
    "github.com/paulmach/orb"
)

type TwoPointLine struct {
    sx float64
    sy float64
    ex float64
    ey float64
}

func LineStringCuttingPath(ls orb.LineString) orb.LineString {
    fin := make(orb.LineString, 0, len(ls))
    if len(ls) < 2 {
        return fin
    }
    tpn := make([]TwoPointLine, 0, len(ls) - 1)
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
        tpn = append(tpn, TwoPointLine{zify(sx+nx),zify(sy+ny),zify(ex+nx),zify(ey+ny)})
    }
    var fx float64 = tpn[0].ex
    var fy float64 = tpn[0].ey
    fin = append(fin, orb.Point{tpn[0].sx,tpn[0].sy})
    for i := 1; i < len(tpn); i++ {
        sxa := tpn[i-1].sx
        sya := tpn[i-1].sy
        exa := tpn[i-1].ex
        eya := tpn[i-1].ey
        sxb := tpn[i].sx
        syb := tpn[i].sy
        exb := tpn[i].ex
        eyb := tpn[i].ey
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
    return fin
}

func LineStringTranslate(ls orb.LineString, dx, dy float64) orb.LineString {
    fin := make(orb.LineString, 0, len(ls))
    for i := 0; i < len(ls); i++ {
        fin = append(fin, orb.Point{ls[i][0] + dx, ls[i][1] + dy})
    }
    return fin
}

func zify(value float64) float64 {
    if value < 0.000001 && value > -0.000001 {
        return 0.0
    }
    return value
}
