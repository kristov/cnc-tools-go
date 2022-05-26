package main

import (
    "os"
//    "flag"
    "math"
    "fmt"
    "bufio"
    "github.com/Succo/wkttoorb"
    "github.com/paulmach/orb"
)

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
            //case orb.Polygon:
            //    doPolygon(t)
            case orb.LineString:
                doLineString(t)
            default:
                fmt.Printf("skipping object of unknown type %T\n", t)
        }
    }
}

func doLineString(ls orb.LineString) {
    if len(ls) < 2 {
        fmt.Println("this linestring does not have enough points")
        return
    }
    tp := make([][2][2]float64, 0, len(ls) - 1)
    for i := 1; i < len(ls); i++ {
        tp = append(tp, [2][2]float64{{ls[i-1][0], ls[i-1][1]},{ls[i][0], ls[i][1]}})
    }
    tpn := make([][2][2]float64, 0, len(ls) - 1)
    // [[[30 10] [10 30]] [[10 30] [40 40]]]
    for i := 0; i < len(tp); i++ {
        dy := tp[i][1][1] - tp[i][0][1]
        dx := tp[i][1][0] - tp[i][0][0]
        angle := math.Atan(dy / dx)
        nangle := angle - (math.Pi / 2)
        ny := math.Sin(nangle) * 2
        nx := math.Cos(nangle) * 2
        tpn = append(tpn, [2][2]float64{{tp[i][0][0] + nx, tp[i][0][1] + ny}, {tp[i][1][0] + nx, tp[i][1][1] + ny}})
    }
    fmt.Println(tpn)
}
