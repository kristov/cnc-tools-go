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
