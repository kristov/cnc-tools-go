package cnclib

import (
    "testing"
    "github.com/paulmach/orb"
)

func TestPointInPoly(t *testing.T) {
    // Counter clockwise square
    ccwls := orb.LineString{orb.Point{0,0},orb.Point{10,0},orb.Point{10,10},orb.Point{0,10},orb.Point{0,0}}
    if PointInPoly(20, 20, ccwls) {
        t.Log("Point should not be in poly")
    }
    if !PointInPoly(5, 5, ccwls) {
        t.Log("Point should be in poly")
    }
    if !PointInPoly(10, 10, ccwls) {
        t.Log("Point should be in poly")
    }
    if !PointInPoly(0, 5, ccwls) {
        t.Log("Point should be in poly")
    }
    if !PointInPoly(5, 0, ccwls) {
        t.Log("Point should be in poly")
    }

    // Clockwise square
    cwls := orb.LineString{orb.Point{0,0},orb.Point{0,10},orb.Point{10,10},orb.Point{10,0},orb.Point{0,0}}
    if PointInPoly(20, 20, cwls) {
        t.Log("Point should not be in poly")
    }
    if !PointInPoly(5, 5, cwls) {
        t.Log("Point should be in poly")
    }
    if !PointInPoly(10, 10, cwls) {
        t.Log("Point should be in poly")
    }
    if !PointInPoly(0, 5, cwls) {
        t.Log("Point should be in poly")
    }
    if !PointInPoly(5, 0, cwls) {
        t.Log("Point should be in poly")
    }
}

func TestToolPath(t *testing.T) {
    // Counter clockwise square (outer toolpath)
    ccwls := orb.LineString{orb.Point{0,0},orb.Point{10,0},orb.Point{10,10},orb.Point{0,10},orb.Point{0,0}}
    tp1 := ToolPath(ccwls, 2.0)
    if tp1[0][0] != -2.0 && tp1[0][1] != -2.0 {
        t.Log("Toolpath point wrong")
    }
    if tp1[1][0] != 12.0 && tp1[1][1] != -2.0 {
        t.Log("Toolpath point wrong")
    }
    if tp1[2][0] != 12.0 && tp1[2][1] != 12.0 {
        t.Log("Toolpath point wrong")
    }
    if tp1[3][0] != -2.0 && tp1[3][1] != 12.0 {
        t.Log("Toolpath point wrong")
    }
    if tp1[4][0] != -2.0 && tp1[4][1] != -2.0 {
        t.Log("Toolpath point wrong")
    }

    // Clockwise square (inner toolpath)
    cwls := orb.LineString{orb.Point{0,0},orb.Point{0,10},orb.Point{10,10},orb.Point{10,0},orb.Point{0,0}}
    tp2 := ToolPath(cwls, 2.0)
    if tp2[0][0] != 2.0 && tp2[0][1] != 2.0 {
        t.Log("Toolpath point wrong")
    }
    if tp2[1][0] != 2.0 && tp2[1][1] != 8.0 {
        t.Log("Toolpath point wrong")
    }
    if tp2[2][0] != 8.0 && tp2[2][1] != 8.0 {
        t.Log("Toolpath point wrong")
    }
    if tp2[3][0] != 8.0 && tp2[3][1] != 2.0 {
        t.Log("Toolpath point wrong")
    }
    if tp2[4][0] != 2.0 && tp2[4][1] != 2.0 {
        t.Log("Toolpath point wrong")
    }
}

func TestPolyFillNew(t *testing.T) {
    //cwls := orb.LineString{orb.Point{0,0},orb.Point{0,10},orb.Point{10,10},orb.Point{10,0},orb.Point{0,0}}
    cwls := orb.LineString{orb.Point{7.07,0},orb.Point{0,7.07},orb.Point{7.07,14.14},orb.Point{14.14,7.07},orb.Point{7.07,0}}
    PolyFillNew(cwls, 0.8)
}
