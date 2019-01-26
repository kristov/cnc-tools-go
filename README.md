# CNC Tools

Various CNC tools written in Go.

## cnc-stl-view

Simple STL viewer. Mouse movement rotates the part. Mouse or pad scroll to zoom in and out.

  make deps
  make
  bin/cnc-stl-view --stl=resources/cube_10x10x10.stl

Note: change the value of the $(GO) variable.
