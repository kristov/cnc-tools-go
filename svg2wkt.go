package main

import (
    "os"
    "fmt"
    "bufio"
//    "flag"
    "github.com/JoshVarga/svgparser"
    "github.com/JoshVarga/svgparser/utils"
    "github.com/paulmach/orb"
    "github.com/paulmach/orb/encoding/wkt"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    svg, err := svgparser.Parse(reader, true)
    if err != nil {
        panic(err)
    }

    for i := 0; i < len(svg.Children); i++ {
        elem := svg.Children[i]
        if elem.Name == "path" {
            if value, exists := elem.Attributes["d"]; exists {
                path, err := utils.PathParser(value)
                if err != nil {
                    panic(err)
                }
                mls := PathToLineString(path)
                fmt.Println(wkt.MarshalString(mls))
            }
        }
    }
}

func PathToLineString(path *utils.Path) orb.MultiLineString {
    mls := make(orb.MultiLineString, 0, 0)
    for j := 0; j < len(path.Subpaths); j++ {
        ls := make(orb.LineString, 0, 0)
        for k := 0; k < len(path.Subpaths[j].Commands); k++ {
            cmd := path.Subpaths[j].Commands[k]
            if cmd.Symbol == "z" {
                first := ls[0]
                ls = append(ls, first)
                continue
            }
            if (cmd.Symbol == "M") || (cmd.Symbol == "L") {
                ls = append(ls, orb.Point{cmd.Params[0], cmd.Params[1]})
                continue
            }
            panic(fmt.Sprintf("unsupported Symbol: %s", cmd.Symbol))
        }
        mls = append(mls, ls)
    }
    return mls
}
