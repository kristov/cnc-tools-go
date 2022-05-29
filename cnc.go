package main

import (
    "os"
    "fmt"
    "bufio"
    "flag"
    "github.com/Succo/wkttoorb"
    "cnc-tools-go/cnclib"
    "github.com/paulmach/orb"
    "github.com/paulmach/orb/encoding/wkt"
)

func main() {
    var things []orb.Geometry
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

    var cmd string = "help"
    var dx float64 = 0.0
    var dy float64 = 0.0
    var radius float64 = 0.0
    var angle float64 = 0.0
    flag.StringVar(&cmd, "cmd", "help", "The command")
    flag.Float64Var(&dx, "dx", 0.0, "Delta X")
    flag.Float64Var(&dy, "dy", 0.0, "Delta Y")
    flag.Float64Var(&radius, "radius", 0.0, "Radius of cutting tool")
    flag.Float64Var(&angle, "angle", 0.0, "Angle of rotation")
    flag.Parse()

    switch {
        case cmd == "trans":
            doTranslate(things, dx, dy)
        case cmd == "toolpath":
            doToolpath(things, radius)
        case cmd == "rotate":
            doRotate(things, angle)
        default:
            printHelp()
    }
}

func printHelp() {
    fmt.Println("cat geometry.wkt | cnc [command] [arg1, arg2, arg3]")
    fmt.Println("")
    fmt.Println("  COMMANDS:")
    fmt.Println("")
    fmt.Println("    trans - cnc trans --dx 10.0 --dy 5.2")
}

func doRotate(things []orb.Geometry, angle float64) {
    lss := getLineStrings(things)
    for i := 0; i < len(lss); i++ {
        fin := cnclib.LineStringRotate(lss[i], angle / 57.29578)
        fmt.Println(wkt.MarshalString(fin))
    }
}
func doTranslate(things []orb.Geometry, dx, dy float64) {
    lss := getLineStrings(things)
    for i := 0; i < len(lss); i++ {
        fin := cnclib.LineStringTranslate(lss[i], dx, dy)
        fmt.Println(wkt.MarshalString(fin))
    }
}

func doToolpath(things []orb.Geometry, radius float64) {
    lss := getLineStrings(things)
    for i := 0; i < len(lss); i++ {
        fin := cnclib.LineStringCuttingPath(lss[i])
        fmt.Println(wkt.MarshalString(fin))
    }
}

func getLineStrings(things []orb.Geometry) []orb.LineString {
    var linestrings []orb.LineString
    for i := 0; i < len(things); i++ {
        switch t := things[i].(type) {
            case orb.LineString:
                linestrings = append(linestrings, orb.LineString(t))
            case orb.Polygon:
                for j := 0; j < len(t); j++ {
                    linestrings = append(linestrings, orb.LineString(t[j]))
                }
            default:
                fmt.Printf("skipping object of unknown type %T\n", t)
        }
    }
    return linestrings
}
