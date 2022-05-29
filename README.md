# CNC Tools

A CNC tool for manipulating polygons (linestrings).

    cat cube.wkt | cnc --cmd=rotate --angle=70 | cnc --cmd=trans --dx 60 --dy 10

The `cnc` command takes WKT format as STDIN and produces WKT on STDOUT. For example, `cube.wkt`:

    LINESTRING(0 0, 20 0, 20 20, 0 20, 0 0)

## `rotate`

Rotate a shape about 0,0 by an angle in degrees:

    $ cat cube.wkt | cnc --cmd=rotate --angle=45
    LINESTRING(0 0,14.14 14.14,0 28.28,-14.14 14.14,0 0)

## `trans`

Translate a shape by delta X and delta Y:

    $ cat cube.wkt | cnc --cmd=trans --dx 20 --dy 20
    LINESTRING(20 20,40 20,40 40,20 40,20 20)

## `toolpath`

Generate a toolpath for cutting. This generates an "outline" around the shape, at a distance of the radius of the cutting tool:

    $ cat cube.wkt | ./cnc --cmd=toolpath

This is generated in a counter-clockwise direction for exterior shapes. If you want to cut out a shape, make sure the path of the shape is counter-clockwise. If you want to generate a hole make sure the path goes clockwise. This is due to the spindle rotation of the CNC being clockwise. You want the tool cutting direction to be opposite to the direction of the shape you are cutting out.

TODO: for closed linestrings it doesn't do the final intersection yet.

## Composition

Because of the STDIN/STDOUT thing, commands can be composed:

    $ cat cube.wkt | cnc --cmd=rotate --angle=45 | cnc --cmd=trans --dx 20 --dy 20
    LINESTRING(20 20,34.14 34.14,20 48.28,5.86 34.14,20 20)

## Viewing

At any point you can view the operation:

    $ cat cube.wkt | cnc --cmd=rotate --angle=45 | cnc-path-view

Which will open a window and show the shape on the screen.

## Generating GCode

TODO:

    $ cat cube.wkt | cnc --cmd=gcode

Should spit out GCode.

