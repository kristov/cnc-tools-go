package main

import (
    "os"
    "fmt"
    "bufio"
    "flag"
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

    var dx float64 = 0.0
    var dy float64 = 0.0
    flag.Float64Var(&dx, "dx", 0.0, "Delta X")
    flag.Float64Var(&dy, "dy", 0.0, "Delta Y")
    flag.Parse()

    for i := 0; i < len(things); i++ {
        switch t := things[i].(type) {
            case orb.LineString:
                doLineString(t, dx, dy)
            default:
                fmt.Printf("skipping object of unknown type %T\n", t)
        }
    }
}

func doLineString(ls orb.LineString, dx, dy float64) {
    fin := cnc.LineStringTranslate(ls, dx, dy)
    fmt.Println(wkt.MarshalString(fin))
}
