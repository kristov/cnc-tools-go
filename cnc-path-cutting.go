package main

import (
    "os"
    "fmt"
    "bufio"
    "github.com/Succo/wkttoorb"
    "cnc-tools-go/cnc"
    "github.com/paulmach/orb"
    "github.com/paulmach/orb/encoding/wkt"
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
    fin := cnc.LineStringCuttingPath(ls)
    fmt.Println(wkt.MarshalString(fin))
}
