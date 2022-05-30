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

/*
    addcmd := flag.NewFlagSet("add", flag.ExitOnError)
    a_add := addcmd.Int("a", 0, "The value of a")
    b_add := addcmd.Int("b", 0, "The value of b")

    mulcmd := flag.NewFlagSet("mul", flag.ExitOnError)
    a_mul := mulcmd.Int("a", 0, "The value of a")
    b_mul := mulcmd.Int("b", 0, "The value of b")

    switch os.Args[1] {
    case "add":
        addcmd.Parse(os.Args[2:])
        fmt.Println(*a_add + *b_add)
    case "mul":
        mulcmd.Parse(os.Args[2:])
        fmt.Println(*(a_mul) * (*b_mul))
    default:
        fmt.Println("expected add or mul command")
        os.Exit(1)
    }
*/

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

    lss := cnclib.GeometryToLineStrings(things)
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
