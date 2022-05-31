package main

import (
    "os"
    "fmt"
    "bufio"
    "flag"
    "cnc-tools-go/cnclib"
    "github.com/Succo/wkttoorb"
    "github.com/paulmach/orb"
    "github.com/paulmach/orb/encoding/wkt"
)

func main() {
    if len(os.Args) == 1 {
        PrintHelp()
        os.Exit(0)
    }

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
    lss := cnclib.GeometryToLineStrings(things)

    trcmd := flag.NewFlagSet("translate", flag.ExitOnError)
    treco := trcmd.Bool("echo", false, "Echo the input geometry")
    trdxv := trcmd.Float64("dx", 0.0, "Delta X")
    trdyv := trcmd.Float64("dy", 0.0, "Delta Y")

    rocmd := flag.NewFlagSet("rotate", flag.ExitOnError)
    roeco := rocmd.Bool("echo", false, "Echo the input geometry")
    roang := rocmd.Float64("angle", 0.0, "Angle of rotation in degrees")

    //mxcmd := flag.NewFlagSet("mirrorx", flag.ExitOnError)
    //mycmd := flag.NewFlagSet("mirrory", flag.ExitOnError)

    tpcmd := flag.NewFlagSet("toolpath", flag.ExitOnError)
    tpeco := tpcmd.Bool("echo", false, "Echo the input geometry")
    tprad := tpcmd.Float64("radius", 1.5, "Radius of cutting tool")

    gccmd := flag.NewFlagSet("gcode", flag.ExitOnError)
    gceco := gccmd.Bool("echo", false, "Echo the input geometry (as a GCode comment)")
    gcclr := gccmd.Float64("clearance", 3.0, "Height tool is lifted to before rapid movement")
    gcdth := gccmd.Float64("depth", 1.0, "Height tool is dropped to before cutting")

    switch os.Args[1] {
        case "help":
            PrintHelp()
        case "translate":
            trcmd.Parse(os.Args[2:])
            for i := 0; i < len(lss); i++ {
                if *treco {
                    fmt.Println(wkt.MarshalString(lss[i]))
                }
                fin := cnclib.Translate(lss[i], *trdxv, *trdyv)
                fmt.Println(wkt.MarshalString(fin))
            }
        case "rotate":
            rocmd.Parse(os.Args[2:])
            for i := 0; i < len(lss); i++ {
                if *roeco {
                    fmt.Println(wkt.MarshalString(lss[i]))
                }
                fin := cnclib.Rotate(lss[i], *roang / 57.29578)
                fmt.Println(wkt.MarshalString(fin))
            }
        case "mirrorx":
            for i := 0; i < len(lss); i++ {
                fin := cnclib.MirrorX(lss[i])
                fmt.Println(wkt.MarshalString(fin))
            }
        case "mirrory":
            for i := 0; i < len(lss); i++ {
                fin := cnclib.MirrorY(lss[i])
                fmt.Println(wkt.MarshalString(fin))
            }
        case "reverse":
            for i := 0; i < len(lss); i++ {
                fin := cnclib.Reverse(lss[i])
                fmt.Println(wkt.MarshalString(fin))
            }
        case "toolpath":
            tpcmd.Parse(os.Args[2:])
            for i := 0; i < len(lss); i++ {
                if *tpeco {
                    fmt.Println(wkt.MarshalString(lss[i]))
                }
                fin := cnclib.ToolPath(lss[i], *tprad)
                fmt.Println(wkt.MarshalString(fin))
            }
        case "gcode":
            gccmd.Parse(os.Args[2:])
            for i := 0; i < len(lss); i++ {
                var gcodes []string
                if *gceco {
                    gcodes = append(gcodes, fmt.Sprintf("( WKT: '%s' )", wkt.MarshalString(lss[i])))
                }
                gcodes = append(gcodes, GToolUp(*gcclr))
                gcodes = append(gcodes, GMoveTo(lss[i][0]))
                gcodes = append(gcodes, GToolDown(*gcdth))
                for j := 1; j < len(lss[i]); j++ {
                    gcodes = append(gcodes, GCutTo(lss[i][j]))
                }
                gcodes = append(gcodes, GToolUp(*gcclr))
                for j := 0; j < len(gcodes); j++ {
                    fmt.Println(gcodes[j])
                }
            }
        default:
            fmt.Printf("unknown command '%s', choose one of: translate, rotate, mirrorx, mirrory, toolpath, help\n", os.Args[1])
            os.Exit(1)
    }
}

func PrintHelp() {
    fmt.Println("cat geometry.wkt | cnc [command] [arg1, arg2, arg3]")
    fmt.Println("")
    fmt.Println("    cnc translate --dx=10.0 --dy=5.2")
    fmt.Println("    cnc rotate --angle=45.0")
    fmt.Println("    cnc mirrory")
    fmt.Println("    cnc mirrorx")
    fmt.Println("    cnc toolpath")
}

func GMoveTo(point orb.Point) string {
    return fmt.Sprintf("G00 X%0.1f Y%0.1f", point[0], point[1])
}

func GCutTo(point orb.Point) string {
    return fmt.Sprintf("G01 X%0.1f Y%0.1f", point[0], point[1])
}

func GToolUp(clearance float64) string {
    return fmt.Sprintf("G00 Z%0.1f", clearance)
}

func GToolDown(depth float64) string {
    return fmt.Sprintf("G00 Z%0.1f", depth)
}
