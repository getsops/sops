Forked from [Vault](https://github.com/hashicorp/vault/tree/master/shamir)

## How it works

We want to split a secret into parts.

Any two points on the cartesian plane define a line. Three points define a
parabola. Four points define a cubic curve, and so on. In general, `n` points
define an function of degree `(n - 1)`. If our secret was somehow an function of
degree `(n - 1)`, we could just compute `n` different points of that function
and give `n` different people one point each. In order to recover the secret,
then we'd need all the `n` points. If we wanted to, we could compute more than n
points, but even then still only `n` points out of our whole set of computed
points would be required to recover the function.

A concrete example: our secret is the function `y = 2x + 1`. This function is of
degree 1, so we need at least 2 points to define it. For example, let's set
`x = 1`, `x = 2` and `x = 3`. From this follows that `y = 3`, `y = 5` and
`y = 7`, respectively. Now, with the information that our secret is of degree 1,
we can use any 2 of the 3 points we computed to recover our original function.
For example, let's use the points `x = 1; y = 3` and `x = 2; y = 5`.
We know that first degree functions are lines, defined by their slope and their
intersection point with the y axis. We can easily compute the slope given our
two points: it's the change in `y` divided by the change in `x`:
`(5 - 3)/(2 - 1) = 2`. Now, knowing the slope we can compute the intersection
point with the `y` axis by "working our way back". We know that at `x = 1`,
`y` equals `3`, so naturally because the slope is `2`, at `x = 0`, `y` must be
`1`.

## Lagrange interpolation

The method we've used for this isn't very general: it only works for polynomials
of degree 1. Lagrange interpolation is a more general way that lets us obtain
the function of degree `(n - 1)` that passes through `n` arbitrary points.

Understanding how to perform Lagrange interpolation isn't really necessary to
understand Shamir's Secret Sharing: it's enough to know that there's only one
function of degree `(n - 1)` that passes through `n` given points and that
computing this function given the points is computationally efficient.

But for those interested, here's an explanation:

Let's say our points are `(x_0, y_0),...,(x_j, y_j),...,(x_(n-1), y_(n-1))`.
Then, the Lagrange polynomial `L(x)`, the polynomial we're looking for, is
defined as follows:

`L(x) = sum from j=0 to j=(n-1) of {y_j * l_j(x)}`

and `l_j(x) = product from m=0 to m=(n-1) except when m=j of {(x - x_m)/(x_j - x_m)}`

A concrete example, with 3 points:

```
x_0 = 1   y_0 = 1
x_1 = 2   y_1 = 4
x_2 = 3   y_2 = 9
```

Let's apply the formula:

```
L(x) =
        y_0 * l_0(x) +
        y_1 * l_1(x) +
        y_2 * l_2(x)
```

Substitute `y_j` for the actual value:

```
L(x) =
        1 * l_0(x) +
        4 * l_1(x) +
        9 * l_2(x)
```

Replace `l_j(x)`:

```
l_0(x) = (x - 2)/(1 - 2) * (x - 3)/(1 - 3) =  0.5x^2 - 2.5x + 3
l_1(x) = (x - 1)/(2 - 1) * (x - 3)/(2 - 3) = -   x^2 +   4x - 3
l_2(x) = (x - 1)/(3 - 1) * (x - 2)/(3 - 2) =  0.5x^2 - 1.5x + 1
```

```

L(x) =
        1 * ( 0.5x^2 - 2.5x + 3) +
        4 * (   -x^2 +   4x - 3) +
        9 * ( 0.5x^2 - 1.5x + 1)
```

```

L(x) =
        ( 0.5x^2 -  2.5x +  3) +
        (  -4x^2 +   16x - 12) +
        ( 4.5x^2 - 13.5x +  9)
     =       x^2 +    0x +  0
     = x^2
```

So the polynomial we were looking for is `y = x^2`.


## Splitting a secret

So we have the ability of splitting a function into parts, but in the context
of computing we generally want to split a number, not a function. For this,
let's define a function of degree `threshold`. `threshold` is the amount of
parts we want to require in order to recover the secret. Let's set the parameter
of degree zero to our secret `S` and make the rest of the parameters random:

`y = ax^(threshold) + bx^(threshold-1) + ... + zx^1 + S`

With `a, b, ...` random.

Then, we want to generate our parts. For this, we evaluate our function at as
many points as we want parts. For example, say our secret is 123, we want 5
parts and a threshold of 2. Because the threshold is 2, we're going to need a
polynomial of degree 2:

`y = ax^2 + bx + 123`

We randomly set `a = 7` and `b = 1`:

`y = 7x^2 + x + 123`

Because we want 5 parts, we need to compute 5 points:

```
x = 0 -> y = 123 # woops! This is the secret itself. Let's not use that one.
x = 1 -> y = 131
x = 2 -> y = 153
x = 3 -> y = 189
x = 4 -> y = 239
x = 5 -> y = 303
```

And that's it. Each of the computed points is one part of the secret.

## Combining a secret

Now that we have our parts, we have to define a way to recover them. Using
the example from the previous section, we only need any two points out of the
five we created to recover the secret, because we set the threshold to two.
So with any two of the five points we created, we can recover the original
polynomial, and because the secret is the free term in the polynomial, we can
recover the secret.

## Finite fields

In the previous examples we've only used integers, and this unfortunately has
a flaw. First of all, it's impossible to uniformly sample integers to get
random coefficients for our generated polynomial. Additionally, if we don't
operate in a finite field, information about the secret is leaked for every part
someone recovers.

For these reasons, Vault's implementation of Shamir's Secret Sharing uses finite
field arithmetic, specifically in GF(2^8), with 229 as the generator. GF(2^8)
has 256 elements, so using this we can only split one byte at a time. This is
not a problem, though, as we can just split each byte in our secret
independently. This implementation uses tables to speed up the execution of
finite field arithmetic.
