# CNC Tools

This is a collection of command line tools for building CNC workflows for an XY router table. These tools are suitable for building GCode for cutting out shapes from flat plates of material, and they are not suitable for cutting out 3D profiles - for that look into the [LinuxCNC](https://www.linuxcnc.org/) suite.

The tools are built to take advantage of Unix shell pipes, so they accept input from STDIN and produce output to STDOUT. Input is generally in the [WKT](https://en.wikipedia.org/wiki/Well-known_text_representation_of_geometry) format unless otherside stated. WKT is a nice and simple text representation for lines, polygons etc - shapes built using straight lines. If you need curves I suggest generating them in a program like OpenSCAD where you can choose how many segments (resolution) you want, and export to DXF or SVG. You can then use these tools to convert from those formats to WKT.

Programs accept multiple WKT objects, one per line. This means a WKT shape can not be broken up into multiple lines. The order of the points in a WKT LINESTRING matters for these tools. The wikipedia page explains it best:

    The OGC standard definition requires a polygon to be topologically closed.
    It also states that if the exterior linear ring of a polygon is defined in
    a counterclockwise direction, then it will be seen from the "top". Any
    interior linear rings should be defined in opposite fashion compared to the
    exterior ring, in this case, clockwise.

This lines up very well with spindle cutting tools rotating in a clockwise direction when looking down onto the work. The cutting tool is a cylinder of blades, but looking down from above the tool is a circle rotating in a clockwise direction. Imagine the tool cutting a straight line through the work from bottom to top (or from front to back: moving in a positive Y direction). On the left hand side the blades of the tool are cutting into the work in the same direction of travel as the tool is moving through the work. On the right hand side the blades are cutting into the work in the opposite direction of travel of the tool. If you are cutting through thermoplastic you will get a clean cut on the left hand side, but on the right the edge of the cut may have deposits of melted plastic. Therefore, if you are cutting out the outside edge of a shape you want the path to follow the shape in a counterclockwise direction. If you are cutting out a void from inside a shape you want the path to follow the hole in a clockwise direction. This diagram illustrates this point:

![Diagram explaining cutting direction](https://github.com/kristov/cnc-tools-go/blob/master/tool_direction?raw=true)

In short: Make your exterior outline as a counterclockwise path, and any interior holes to be cut as clockwise paths. You can use the `cnc` tool to reverse the direction of a path (LINESTRING) if you need to.

## The `cnc` tool

This is the most useful of the tools, used for performing various transformations on WKT geometries. You can translate, rotate and mirror LINESTRINGS. You can generate a cutting "toolpath", which is an outline of a shape at a given radius distance from the shape. And you can generate GCode from one or more LINESTRINGS.

A CNC tool for manipulating polygons (linestrings).

    $ echo "LINESTRING(0 0, 20 0, 20 20, 0 20, 0 0)" | cnc translate --dx 20 --dy 20
    LINESTRING(20 20,40 20,40 40,20 40,20 20)

Or if you have this same LINESTRING in a WKT file:

    $ cat cube.wkt | cnc translate --dx 20 --dy 20
    LINESTRING(20 20,40 20,40 40,20 40,20 20)

### `translate`

Translate a shape by delta X and delta Y:

    $ cat cube.wkt | cnc translate --dx 20 --dy 20
    LINESTRING(20 20,40 20,40 40,20 40,20 20)

### `rotate`

Rotate a shape about 0,0 by an angle in degrees:

    $ cat cube.wkt | cnc rotate --angle=45
    LINESTRING(0 0,14.14 14.14,0 28.28,-14.14 14.14,0 0)

### `mirrory`

This mirrors the shape in the Y axis (produces a mirror to the left)

    $ cat cube.wkt | cnc mirrory
    LINESTRING(0 0,-20 0,-20 20,0 20,0 0)

Note: this makes the geometry go negative in the X direction, so you most likely want to `translate` it back into positive X coordinate space after using composition (see "Composition" below).

### `mirrorx`

This mirrors the shape in the X axis (produces a mirror below).

    $ cat cube.wkt | cnc mirrorx
    LINESTRING(0 0,20 0,20 -20,0 -20,0 0)

### `toolpath`

Generate a toolpath for cutting. This generates an "outline" around the shape, at a distance of the radius of the cutting tool:

    $ cat cube.wkt | cnc toolpath --radius=1.5

This is generated in a counter-clockwise direction for exterior shapes. If you want to cut out a shape, make sure the path of the shape is counter-clockwise. If you want to generate a hole make sure the path goes clockwise. This is due to the spindle rotation of the CNC being clockwise. You want the tool cutting direction to be opposite to the direction of the shape you are cutting out.

TODO: for closed linestrings it doesn't do the final intersection yet.

### The `echo` option

It can be useful to echo the original input for viewing purposes:

    $ cat cube.wkt | cnc toolpath --radius=1.5 --echo
    LINESTRING(0 0,20 0,20 20,0 20,0 0)
    LINESTRING(-1.5 -1.5,21.5 -1.5,21.5 21.5,-1.5 21.5,-1.5 -1.5)

While probably not useful for generating GCode this can be used to verify some translation by piping this "doubled output" to `cnc-view2d` so you can see the original plus the translation (see "Viewing" below).

### Composition

Because of the STDIN/STDOUT thing, commands can be composed:

    $ cat cube.wkt | cnc rotate --angle=45 | cnc translate --dx 20 --dy 20
    LINESTRING(20 20,34.14 34.14,20 48.28,5.86 34.14,20 20)

Beware of multiple `--echo` flags in a chain of commands. The input is not aware that a shape from a previous output was an echo and will perform an operation on the echo like it was just another path.

## Viewing

At any point you can view the result of an operation:

    $ cat cube.wkt | cnc rotate --angle=45 | cnc-view2d

Which will open a window and show the shape on the screen. The mouse wheel will zoom in and out. Pan is not available. If you need more complex viewing functionality, consider outputting to SVG.

## Generating GCode

TODO:

    $ cat cube.wkt | cnc gcode

Should spit out GCode.

