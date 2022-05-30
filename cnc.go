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
    flag.Float64Var(&angle, "angle", 0.0, "Angle of rotation in degrees")
    flag.Parse()

    lss := getLineStrings(things)
    switch {
        case cmd == "trans":
            for i := 0; i < len(lss); i++ {
                fin := cnclib.Translate(lss[i], dx, dy)
                fmt.Println(wkt.MarshalString(fin))
            }
        case cmd == "toolpath":
            for i := 0; i < len(lss); i++ {
                fin := cnclib.CuttingPath(lss[i])
                fmt.Println(wkt.MarshalString(fin))
            }
        case cmd == "rotate":
            for i := 0; i < len(lss); i++ {
                fin := cnclib.Rotate(lss[i], angle / 57.29578)
                fmt.Println(wkt.MarshalString(fin))
            }
        case cmd == "mirrory":
            for i := 0; i < len(lss); i++ {
                fin := cnclib.MirrorY(lss[i])
                fmt.Println(wkt.MarshalString(fin))
            }
        case cmd == "mirrorx":
            for i := 0; i < len(lss); i++ {
                fin := cnclib.MirrorX(lss[i])
                fmt.Println(wkt.MarshalString(fin))
            }
        default:
            printHelp()
    }
}

func printHelp() {
    fmt.Println("cat geometry.wkt | cnc [command] [arg1, arg2, arg3]")
    fmt.Println("")
    fmt.Println("    cnc --cmd=trans --dx=10.0 --dy=5.2")
    fmt.Println("    cnc --cmd=rotate --angle=45.0")
    fmt.Println("    cnc --cmd=mirrory")
    fmt.Println("    cnc --cmd=mirrorx")
    fmt.Println("    cnc --cmd=toolpath")
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
