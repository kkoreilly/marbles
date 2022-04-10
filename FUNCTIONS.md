List of supported functions:

<details>
<summary>sin(v)</summary>


**Definition:**
Value is the sine of v.

**Examples:**
* sin(pi) will be equal to 0
* sin(3pi/2) will be equal to -1
</details>

<details>
<summary>cos(v)</summary>


**Definition:**
Value is the cosine of v.

**Examples:**
* cos(pi) will be equal to -1
* sin(3pi/2) will be equal to 0
</details>

<details>
<summary>tan(v)</summary>


**Definition:**
Value is the tangent of v.

**Examples:**
* tan(pi) will be equal to 0
* tan(3pi/4) will be equal to -1
</details>

<details>
<summary>sec(v)</summary>


**Definition:**
Value is the secant of v.

**Examples:**
* sec(pi) will be equal to -1
* sec(3pi/2) will be undefined
</details>

<details>
<summary>csc(v)</summary>


**Definition:**
Value is the cosecant of v.

**Examples:**
* csc(pi) will be undefined
* csc(3pi/2) will be equal to -1
</details>

<details>
<summary>cot(v)</summary>


**Definition:**
Value is the cotangent of v.

**Examples:**
* cot(pi/2) will be equal to 0
* cot(pi) will be undefined
</details>

<details>
<summary>arcsin(v)</summary>


**Definition:**
Value is the inverse/arc sine of v.

**Examples:**
* arcsin(1) will be equal to pi/2
* arcsin(-1) will be equal to -pi/2
</details>

<details>
<summary>arccos(v)</summary>


**Definition:**
Value is the inverse/arc cosine of v.

**Examples:**
* arccos(1) will be equal to pi
* arccos(0) will be equal to pi/2
</details>

<details>
<summary>arctan(v)</summary>


**Definition:**
Value is the inverse/arc tangent of v.

**Examples:**
* arctan(0) will be equal to 0
* arctan(1) will be equal to pi/4
</details>

<details>
<summary>arcsec/arccsc/arccot(v)</summary>


**Definition:**
Value is the inverse/arc secant/cosecant/cotangent of v.

**Examples:**
* arcsec(1) will be equal to 0
* arccsc(1) will be equal to pi/2
</details>

<details>
<summary>sinh/cosh/tanh/sech/csch/coth/arcsinh/arccosh/arctanh/arcsech/arccsch/arccoth(v)</summary>


**Definition:**
Value is the hyperbolic sine/cosine/tangent/secant/cosecant/cotangent/arcsine/arccosine/arctangent/arcsecant/arccosecant/arccotangent of v.

**Examples:**
* sinh(0) will be equal to 0
* cosh(0) will be equal to 1
</details>

<details>
<summary>pow(v1, v2)</summary>


**Definition:**
Value is v1 raised to the v2th power.

**Examples:**
* pow(2, 4) will be equal to 16
* pow(5, 3) will be equal to 125
</details>

<details>
<summary>abs(v)</summary>


**Definition:**
Value is the absolute value of v.

**Examples:**
* abs(-52) will be equal to 52
* abs(28) will be equal to 28
</details>

<details>
<summary>fact(v)</summary>


**Definition:**
Value is the factorial of v.

**Examples:**
* fact(3) will be equal to 6
* fact(4) will be equal to 24
</details>

<details>
<summary>ceil(v)</summary>


**Definition:**
Value is v, rounded up to the nearest integer.

**Examples:**
* ceil(4.001) will be equal to 5
* ceil(-6.8) will be equal to -6
</details>

<details>
<summary>floor(v)</summary>


**Definition:**
Value is v, rounded down to the nearest integer.

**Examples:**
* floor(4.9999) will be equal to 4
* floor(-7.2) will be equal to -8
</details>

<details>
<summary>round(v)</summary>


**Definition:**
Value is v, rounded to the nearest integer.

**Examples:**
* round(4.5) will be equal to 5
* round(3.1) will be equal to 3
</details>

<details>
<summary>min/max(v1, v2...)</summary>


**Definition:**
Value is the lowest/highest value of the provided values

**Examples:**
* max(3, 5) will be equal to 5
* min(1, -7.2, 3) will be equal to -7.2
</details>

<details>
<summary>avg(v1, v2...)</summary>


**Definition:**
Value is the average of the provided values.

**Examples:**
* avg(1, 2) will be equal to 1.5
* avg(0, 3, 6) will be equal to 3
</details>

<details>
<summary>mod(v1, v2)</summary>


**Definition:**
Value is the modulo operation using v1 and v2. The value is the remainder of v1/v2.

**Examples:**
* mod(1, 2) will be equal to 1
* mod(3, 3) will be equal to 0
</details>

<details>
<summary>sqrt(v)</summary>


**Definition:**
Value is the square root of v.

**Examples:**
* sqrt(4) will be equal to 2
* sqrt(1024) will be equal to 32
</details>

<details>
<summary>ln(v)</summary>


**Definition:**
Value is the natural logarithm of v.

**Examples:**
* ln(e) will be equal to 1
* ln(1) will be equal to 0
</details> 

<details>
<summary>rand(v)</summary>


**Definition:**
Value is a random number between 0 and 1 multiplied by v. The random number is determined on graph, so the value will be different each time the marbles are ran. The value of rand does not change while the marbles are being ran.

**Examples:**
* rand(1) will give a random value between 0 and 1.
* rand(pi) will give a random value between 0 and pi.
</details>

<details>
<summary>if(condition, v1, v2)</summary>


**Definition:**
If condition is true, value is v1. Otherwise value is v2

**Examples:**
* ife(0, 1, 1, 2) will be equal to 2
* ife(sin(pi), 0, 1, 2) will be equal to 1
</details>


<details>
<summary>funcName(v)</summary>


**Definition:**
Value is the value of the function specified evaluated at v. View the name of a function to the right of the equation.

**Examples:**
* If f(x) is x^2 + 1, then f(2) will be equal to 5.
* If b(x) is sqrt(x)+5, then b(16) will be equal to 9.
</details>

<details>
<summary>funcName'(v)</summary>


**Definition:**
Value is the derivative of the function specified at v. View the name of a function to the left of the equation.

**Examples:**
* If f(x) is x^2 + 1, then f'(2) will be equal to 4.
* If b(x) is sqrt(x)+5, then b'(16) will be equal to 1/8.
</details>

<details>
<summary>funcName''(v)</summary>


**Definition:**
Value is the second derivative of the function specified at v. View the name of a function to the left of the equation.

**Examples:**
* If f(x) is x^2 + 1, then f''(2) will be equal to 2.
* If b(x) is x^4, then b''(4) will be equal to 192.
</details>

<details>
<summary>funcNamei(a, b)</summary>


**Definition:**
Value is the integral of the function specified from a to b. View the name of a function to the left of the equation.

**Examples:**
* If f(x) is 2x, then fi(0, 1) will be equal to 1.
* If b(x) is x^2, then bi(1, 2) will be equal to 7/3.
</details>

<details>
<summary>FUNCNAME(v)</summary>


**Definition:**
Value is the antiderivative of the function specified at v. View the name of a function to the left of the equation.

**Examples:**
* If f(x) is x^2 + 1, then F(2) will be equal to 14/3.
* If b(x) is x, then B(1) will be equal to 1/2.
</details>