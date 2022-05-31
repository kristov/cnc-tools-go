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

![Diagram explaining cutting direction](https://github.com/kristov/cnc-tools-go/blob/master/tool_direction.png?raw=true)

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

### `reverse`

Reverses the direction of the points in a LINESTRING. Hard to demonstrate with the cube:

    $ cat cube.wkt | cnc reverse --echo
    LINESTRING(0 0,20 0,20 20,0 20,0 0)
    LINESTRING(0 0,0 20,20 20,20 0,0 0)

(see "echo" below for why I did that)

### `toolpath`

Generate a toolpath for cutting. This generates an "outline" around the shape, at a distance of the radius of the cutting tool:

    $ cat cube.wkt | cnc toolpath --radius=1.5

This is generated in a counter-clockwise direction for exterior shapes. If you want to cut out a shape, make sure the path of the shape is counter-clockwise. If you want to generate a hole make sure the path goes clockwise. This is due to the spindle rotation of the CNC being clockwise. You want the tool cutting direction to be opposite to the direction of the shape you are cutting out.

WARNING! The "toolpath" code is broken for negative (hole) cuts containing complex corners. To illustrate, see this:

    cat resources/moon.svg | svg2wkt | cnc translate --dx=30 --dy=30 | cnc toolpath --echo | cnc-view2d

See how the sharp corners are all messed up? Needs work.

## `gcode`

The `gcode` command will generate very basic GCode for each LINESTRING in the input. The `clearance` and `depth` parameters control how the tool is lifted while moving (clearance) and to what height the Z axis will move to when cutting (depth). The default values are meant to be relatively safe (no cutting will take place), assuming Z=0 is where the cutting tool is barely in contact with the work:

    $ cat cube.wkt | cnc gcode --clearance=2.0 --depth=1.0

This means when you actually want to make a real cut you will need to provide an appropriate `depth` value (negative).

Note: The way I do my cuts is to clamp the work piece to the bed of the CNC machine, with a sacrificial block underneath it. This is because I want the cutting tool to go all the way though the work so I need the tool to be cutting to a depth below the underneath of the work. Without a sacrificial block this means the tool will cut into the metal bed of the CNC machine - *bad mojo*. I then turn on the spindle, jog the X and Y to the corner of the work, and then jog the Z until the spinning tool is *just* touching the work. I then zero the machine coordinates.

I now know that any negative Z position is cutting down into the work, and any positive Z position is above the work. If my work is 2mm thick, I will set the `depth` parameter to `--depth=-2.1` meaning the tool will barely break the bottom surface of the work and not chew up too much of the sacrificial block.

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

## The `cnc-view2d` tool

At any point you can view the result of an operation:

    $ cat cube.wkt | cnc rotate --angle=45 | cnc-view2d

Which will open a window and show the shape on the screen. The mouse wheel will zoom in and out. Pan is not available. The view will show a grey dashed box representing the maximum X and Y travel for the machine. These values default to 200x290 which are the defaults for my machine, so you might need to change them. The `maxx` and `maxy` parameters allow control over the size of this box. It is just a visual aid - the tool does not check if any shapes go outside this box. If your machine has a 100x100 maximum cutting area:

    $ cat cube.wkt | cnc rotate --angle=45 | cnc-view2d --maxx=100 --maxy=100

It is actually better to set `maxx` and `maxy` to the size of the *work* piece, which will be on a case-by-case basis. There are also `width` and `height` options for the default size of the window, but resizing the window manually also works.

## The `svg2wkt` tool

This tool takes an SVG file as input on STDIN and produces WKT as output:

    $ cat resources/moon.svg | svg2wkt | cnc translate --dx=30 --dy=30 | cnc-view2d

It is very basic and will only work with `<path d="...">` elements in the root `<svg>` parent. It does not handle elements like `<circle>`, `<line>`, `<polygon>` etc. Furthermore, the paths processed must use absolute coordinates ("M", "L") not their relative versions ("m", "l"). It was tested based on what OpenSCAD exports.

I also noticed that OpenSCAD seems to be generating SVG with clockwise winding, so I have to "reverse" the output of `svg2wkt` to compensate:

    $ cat resources/moon.svg | svg2wkt | cnc reverse | cnc translate --dx=30 --dy=30 | cnc-view2d

## Cookbook

Common recipies. There is only 1 so far, but more to come I guess.

### 1) Cut out a circle from 2mm thick ply

#### Generate the geometry

Open OpenSCAD and create a circle:

    circle(r=20);

This will generate a circle of a diameter of 40mm. Render it using `F6`. The export it to an SVG file called `circle.svg`.

#### Convert the SVG to WKT

This command will turn the SVG into a MULTILINESTRING WKT command. Make sure the output is as you are expecting.

    $ cat circle.svg | svg2wkt

#### Translate the circle into the work area

By default, OpenSCAD generates a circle around point 0,0. This will mean your CNC machine will try to move in the negative X,Y direction, and this may hit limit switches or break your machine. So you either need to translate the circle in OpenSCAD, or use the "translate" command of the `cnc` tool to move it into positive X,Y coordinates:

    $ cat circle.svg | svg2wkt | cnc translate --dx=20 --dy=20

#### Verify the circle is translated correctly

Now you can view the result and make sure the circle is completely visible:

    $ cat circle.svg | svg2wkt | cnc translate --dx=20 --dy=20 | cnc-view2d

#### Generate and test GCode

You can then generate some test GCode that does *not* plunge the cutting tool into the work. It will trace out the path 1mm above the top of the work. Note: this assumes the Z axis of the CNC machine has been set up to the tool is barely touching the surface of the work at Z=0.0 - this is *very important*. If you have not set up the machine like this, you may need to adjust your clearance and depth values to match, but I don't recommend it. I think it's better to set your Z axis to be zeroed so the cutting tool is just touching the top surface of the work. It means any negative Z value is inside the work, and positive Z values are above the work:

    $ cat circle.svg | svg2wkt | cnc translate --dx=20 --dy=20 | cnc gcode --clearance=2.0 --depth=1.0 > circle_test.gcode

The GCode will lift the tool 2mm above the work for travelling to the start point. It will then descend to 1.0mm *above* the work to simulate the cutting. This allows you to run this GCode in a real machine to validate it will not hit X or Y limits.

#### Generate real GCode

    $ cat circle.svg | svg2wkt | cnc translate --dx=20 --dy=20 | cnc gcode --clearance=2.0 --depth=-2.01 > circle.gcode

Now the tool will descend to -2.1mm below the surface of the work to begin cutting. Obviously this will mean the tool passes through the work, so there needs to be a sacrificial block underneath the work otherwise the tool will cut 0.1mm deep into the surface of the CNC table - *not good*. Always have a sacrificial block under the work.

It is called "Computer Numerical Control" but it is very important that *you* are in control of the machine. Be ready to hit the emergency stop button, or kill power to the machine if you think things are going the wrong way. Always keep the spindle running so you don't accidentially crash a stationary spindle into the work and break it. The computer is executing commands completely blindly, and the only safeguards are the limit switches (if you have them). But do not rely on them - maintain control at all times. Keep a close eye on everything and wear glasses so you don't get chips in your eyes, forcing you to look away.
