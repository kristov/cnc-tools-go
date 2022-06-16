# line2d - Lines are funny things

Small library to try to represent lines in a fairly sensible way. Some issues with lines:

* Lines defined by a start and end point are different than abstract lines defined by a slope and and an intercept
* Line intersection calculations are best done using abstract lines
* A vertical line has an infinite slope which can not be represented in a computer (so we need a flag)
* A vertical line does not really have a Y intercept (unless Y=0 I suppose)
* Can you say two vertical lines with the same X intersect? If so, at what point?

## Angles

Now lets say we have a TwoPointLine and an AbstractLine and we want to know the angle of the line. The AbstractLine does not have a "direction" - it is not a vector. Therefore we can not know if the line is going "back" or "forward" in the X or Y direction. Given that, an AbstractLine can not have an angle greater than 90 degrees (Pi / 2) or less than -90 degrees.

However a TwoPointLine does have a direction, and therefore can have an angle greater than 90 degrees. In this implementation a TwoPointLine can actually *never* have a negative angle - the angle is always between 0 and 360 degrees (between 0 and TwoPi).
