# Marbles App

Graph equations and run marbles on them. Based on [desmos.com](https://desmos.com). Uses [goki/gi](https://github.com/goki/gi) for graphics, and [Knetic/govaluate](https://github.com/Knetic/govaluate) for evaluating equations.  

## Features

* For lines
  * You can set the domain/range of lines you graph (Min/MaxX, Min/MaxY)
  * You can set bounciness for a line. If it is 1, the marble will not gain or lose speed after hitting the line. If bounciness is less than 1, marbles will lose speed when they hit the line. If it is greater than one, marbles will gain speed.
  * You can set the color marbles that hit the line will change to (LineColors.ColorSwitch) - none means the marbles stay the same color.
  * You can also set the color for lines (LineColors.Color).
* For the whole graph
  * You can set the amount of marbles that spawn (NMarbles)
  * You can set the amount of steps the graph runs for (NSteps)
  * You can set the starting speed of the marbles (StartSpeed)
  * You can set the update rate of the marbles (UpdtRate)
  * You can set the gravity of the marbles (Gravity)
  * You can set the range in which the marbles can spawn, 0 makes them spawn in a straight vertical line (Width)
  * You can set the amount the variable t increases every step (TimeStep)
  * You can set the size of the graph (Min/MaxX, Min/MaxY)
* Controls
  * Open allows you to open a saved json file of a graph
  * Save allows you to save a graph to a json file
  * Open autosaved opens the last graph you graphed, helpful if the app crashes
  * Graph graphs all of the lines and resets the marbles to their starting positions
  * Run runs the marbles for NSteps
  * Stop stops the marbles
  * Step runs the marbles for one step
* Settings
  * You can customize the default line that will be added if a line is empty
  * You can customize the default graph parameters
  * You can change the size and color of marbles, default means the marbles will spawn in different random colors. Note that collision is not affected by changing the size of the marbles, so it is not recommended to change the size too much or things will look weird.
  * You can change the color of most things in the app
* Upload and Download graphs
  * You can upload a graph you make by clicking upload graph and giving it a name. Anyone else can download that graph if they want to. Your graph will be automatically deleted after a week.
  * You can view graphs other people have uploaded by clicking download graph and choosing a graph. It will open that graph and you can run it and make any changes to it you want. You can save the graph to your device and do anything you want with it.
* Important note for color and color switch: If you set color switch to white (RGBA(255, 255, 255, 255)), then the marbles will not change color upon hitting that line. It is strongly recommended that you have white as your default line color switch value in settings. Also, if you would like lines to automatically be in random different colors, keep the default color value for lines in settings at white. If you set it to something other than white, all new lines you make will be just that color by default.

## How to write equations

* Equations are written in the form y = ... for now, so no circles.
* There are 4 variables you can use: x, t, and a. x is the x value, h is the amount of times the line has been hit since the marbles were ran, t is the time passed since running the marbles (incremented by TimeStep), and a is a number that goes back and forth between -10 and 10 (defined as 10*sin(t)).
* Note: any value (including variables) can be used for the functions and operations below
* The following operations are supported:
  * +, -, \*, /
  * v1^v2 or v1\*\*v2 for exponents
  * v1%v2 for the modulo operation
* The following functions are supported:
  * sin(v), cos(v), tan(v) - sine, cosine, tangent
  * csc(v), sec(v), cot(v) - cosecant, secant, cotangent
  * asin(v), acos(v), atan(v) - inverse/arc sine, cosine, tangent
  * pow(v1, v2) - same as v1^v2
  * abs(v) - absolute value
  * fact(v) - factorial of a number
  * ceil(v), floor(v) - Round a number up/down
  * mod(v1, v2) - same as v1%v2
  * sqrt(v) - square root of a number
  * ln(v) - natural logarithm of a number
  * rand(v) - graphs a random number between 0 and 1 multiplied by the number specified for every value of x - creates quite a bit of lag and not recommended to use.
  * ife(v1,v2,v3,v4) - If v1 is equal to v2, returns v3. Otherwise returns v4. 
  * ifb(v1,v2,v3,v4,v5) - If v1 is between v2 and v3, returns v4. Otherwise returns v5.
  * f(index) - Value is the value of equation specified. You can see the index of a line to the left of the equation.
  * d(index) - Takes the derivative of the equation at that index. For example, if line 0 is y = sin(x), d(0) is the same as cos(x). If line 4 is y = (e^x)-x^2, d(4) will be equal to (e^x)-2x. You can view the index of a line to the left of its equation. Please note that this function just gets an approximation of the derivative, not the actual derivative.
  * sd(index) - Same as d except takes the second derivative of the equation.
  * i(index, a, b) - Takes the integral (area under the curve) of the equation specified by the index from a to b.

## Known Bugs

* Marbles will go sometimes through lines, especially if the lines are moving (using t or a)
* On functions like tan(x), where x is undefined at a point, the app will draw a vertical line. 

# Images

The app:
![Marbles app](https://github.com/kplat1/marblesInfo/raw/master/images/img1.png)

The app can support a wide variety of functions:
![Marbles app lot of functions](https://github.com/kplat1/marblesInfo/raw/master/images/img2.png)
