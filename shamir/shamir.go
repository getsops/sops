package shamir

// Some comments in this file were written by @autrilla
// The code was written by HashiCorp as part of Vault.

// This implementation of Shamir's Secret Sharing matches the definition
// of the scheme. Other tools used, such as GF(2^8) arithmetic, Lagrange
// interpolation and Horner's method also match their definitions and should
// therefore be correct.
// More information about Shamir's Secret Sharing and Lagrange interpolation
// can be found in README.md

import (
	"crypto/rand"
	"fmt"
)

const (
	// ShareOverhead is the byte size overhead of each share
	// when using Split on a secret. This is caused by appending
	// a one byte tag to the share.
	ShareOverhead = 1
)

// polynomial represents a polynomial of arbitrary degree
type polynomial struct {
	coefficients []uint8
}

// makePolynomial constructs a random polynomial of the given
// degree but with the provided intercept value.
func makePolynomial(intercept, degree uint8) (polynomial, error) {
	// Create a wrapper
	p := polynomial{
		coefficients: make([]byte, degree+1),
	}

	// Ensure the intercept is set
	p.coefficients[0] = intercept

	// Assign random co-efficients to the polynomial
	if _, err := rand.Read(p.coefficients[1:]); err != nil {
		return p, err
	}

	return p, nil
}

// evaluate returns the value of the polynomial for the given x
// Uses Horner's method <https://en.wikipedia.org/wiki/Horner%27s_method> to
// evaluate the polynomial at point x
func (p *polynomial) evaluate(x uint8) uint8 {
	// Special case the origin
	if x == 0 {
		return p.coefficients[0]
	}

	// Compute the polynomial value using Horner's method.
	degree := len(p.coefficients) - 1
	out := p.coefficients[degree]
	for i := degree - 1; i >= 0; i-- {
		coeff := p.coefficients[i]
		out = add(mult(out, x), coeff)
	}
	return out
}

// interpolatePolynomial takes N sample points and returns
// the value at a given x using a lagrange interpolation.
// An implementation of Lagrange interpolation
// <https://en.wikipedia.org/wiki/Lagrange_polynomial>
// For this particular implementation, x is always 0
func interpolatePolynomial(xSamples, ySamples []uint8, x uint8) uint8 {
	limit := len(xSamples)
	var result, basis uint8
	for i := 0; i < limit; i++ {
		basis = 1
		for j := 0; j < limit; j++ {
			if i == j {
				continue
			}
			num := add(x, xSamples[j])
			denom := add(xSamples[i], xSamples[j])
			term := div(num, denom)
			basis = mult(basis, term)
		}
		group := mult(ySamples[i], basis)
		result = add(result, group)
	}
	return result
}

// div divides two numbers in GF(2^8)
// GF(2^8) division using log/exp tables
func div(a, b uint8) uint8 {
	if b == 0 {
		// leaks some timing information but we don't care anyways as this
		// should never happen, hence the panic
		panic("divide by zero")
	}

	// a divided by b is the same as a multiplied by the inverse of b:
	return mult(a, inverse(b))
}

// inverse calculates the inverse of a number in GF(2^8)
// Note that a must be non-zero; otherwise 0 is returned
func inverse(a uint8) uint8 {
	// This makes use of Fermat's Little Theorem for finite groups:
	// If G is a finite group with n elements, and a any element of G,
	// then a raised to the power of n equals the neutral element of G.
	// (See https://en.wikipedia.org/wiki/Fermat%27s_little_theorem;
	// the generalization to finite groups follows from Lagrange's theorem:
	// https://en.wikipedia.org/wiki/Lagrange%27s_theorem_(group_theory))
	//
	// Here we use the multiplicative group of GF(2^8), which has
	// n = 2^8 - 1 elements (every element but zero). Thus raising a to
	// the (n - 1)th = 254th power gives a number x so that a*x = 1.
	//
	// If a happens to be 0, which is not part of the multiplicative group,
	// then a raised to the power of 254 is still 0.

	// (See also https://github.com/openbao/openbao/commit/a209a052024b70bc563d9674cde21a20b5106570)

	// In the comments, we use ^ to denote raising to the power:
	b := mult(a, a)   // b is now a^2
	c := mult(a, b)   // c is now a^3
	b = mult(c, c)    // b is now a^6
	b = mult(b, b)    // b is now a^12
	c = mult(b, c)    // c is now a^15
	b = mult(b, b)    // b is now a^24
	b = mult(b, b)    // b is now a^48
	b = mult(b, c)    // b is now a^63
	b = mult(b, b)    // b is now a^126
	b = mult(a, b)    // b is now a^127
	return mult(b, b) // result is a^254
}

// mult multiplies two numbers in GF(2^8)
// GF(2^8) multiplication using log/exp tables
func mult(a, b uint8) (out uint8) {
	// This computes a * b in GF(2^8), which is defined as GF(2)[X] / <X^8 + X^4 + X^3 + X + 1>.
	// This finite field is known as Rijndael's finite field. (Rijndael is the algorithm that
	// was standardized as AES.)
	// (See https://en.wikipedia.org/wiki/Finite_field_arithmetic#Rijndael's_(AES)_finite_field)
	//
	// We identify elements in GF(2^8) with polynomials of degree < 8. The i-th bit of a field
	// element is the coefficient of X^i in that polynomial.
	//
	// To multiply a and b in this finite field, we use something similar to Russian peasant
	// multiplication. We iterate over b's bits, starting from the highest to the lowest.
	// i denotes the bit we're currently processing (7, 6, 5, 4, 3, 2, 1, 0).
	// The accumulator is set to 0; every iteration, we multiply the accumulator
	// by X modulo X^8+X^4+X^3+X+1, and then add a to the accumulator in case b's i-th bit is 1.
	var accumulator uint8 = 0
	var i uint8 = 8

	for i > 0 {
		i--
		// Get the i-th bit of b; bitOfB is either 0 or 1.
		bitOfB := b >> i & 1
		// aOrZero is 0 if the i-th bit of b is 0, and a if the i-th bit of b is 1. This is
		// what we later add to the accumulator.
		aOrZero := -bitOfB & a
		// zeroOr1B is 0 if the 7th bit of the accumulator is 0, and 0x1B = 11011_2 if the
		// 7th bit of accumulator is 1
		zeroOr1B := -(accumulator >> 7) & 0x1B
		// accumulatorMultipliedByX equals accumulator multiplied by X modulo X^8+X^4+X^3+X+1
		// In the expression, accumulator + accumulator equals accumulator << 1, which would be
		// the accumulator multiplied by X modulo X^8.
		// By XORing (addition and subtraction in GF(2^8)) with zeroOr1B, we turn this into
		// accumulator multiplied by X modulo X^8 + X^4 + X^3 + X + 1.
		accumulatorMultipliedByX := zeroOr1B ^ (accumulator + accumulator)
		// We can now compute the next value of the accumulator as the sum (in GF(2^8)) of aOrZero
		// and accumulatorMultipliedByX.
		accumulator = aOrZero ^ accumulatorMultipliedByX
	}

	return accumulator
}

// add combines two numbers in GF(2^8)
// This can also be used for subtraction since it is symmetric.
func add(a, b uint8) uint8 {
	// Addition in GF(2^8) equals XOR:
	return a ^ b
}

// Split takes an arbitrarily long secret and generates a `parts`
// number of shares, `threshold` of which are required to reconstruct
// the secret. The parts and threshold must be at least 2, and less
// than 256. The returned shares are each one byte longer than the secret
// as they attach a tag used to reconstruct the secret.
func Split(secret []byte, parts, threshold int) ([][]byte, error) {
	// Sanity check the input
	if parts < threshold {
		return nil, fmt.Errorf("parts cannot be less than threshold")
	}
	if parts > 255 {
		return nil, fmt.Errorf("parts cannot exceed 255")
	}
	if threshold < 2 {
		return nil, fmt.Errorf("threshold must be at least 2")
	}
	if threshold > 255 {
		return nil, fmt.Errorf("threshold cannot exceed 255")
	}
	if len(secret) == 0 {
		return nil, fmt.Errorf("cannot split an empty secret")
	}

	// Allocate the output array, initialize the final byte
	// of the output with the offset. The representation of each
	// output is {y1, y2, .., yN, x}.
	out := make([][]byte, parts)
	for idx := range out {
		// Store the x coordinate for each part as its last byte
		// Add 1 to the xCoordinate because if the x coordinate is 0,
		// then the result of evaluating the polynomial at that point
		// will be our secret
		out[idx] = make([]byte, len(secret)+1)
		out[idx][len(secret)] = uint8(idx) + 1
	}

	// Construct a random polynomial for each byte of the secret.
	// Because we are using a field of size 256, we can only represent
	// a single byte as the intercept of the polynomial, so we must
	// use a new polynomial for each byte.
	for idx, val := range secret {
		// Create a random polynomial for each point.
		// This polynomial crosses the y axis at `val`.
		p, err := makePolynomial(val, uint8(threshold-1))
		if err != nil {
			return nil, fmt.Errorf("failed to generate polynomial: %w", err)
		}

		// Generate a `parts` number of (x,y) pairs
		// We cheat by encoding the x value once as the final index,
		// so that it only needs to be stored once.
		for i := 0; i < parts; i++ {
			// Add 1 to the xCoordinate because if it's 0,
			// then the result of p.evaluate(x) will be our secret
			x := uint8(i) + 1
			// Evaluate the polynomial at x
			y := p.evaluate(x)
			out[i][idx] = y
		}
	}

	// Return the encoded secrets
	return out, nil
}

// Combine is used to reverse a Split and reconstruct a secret
// once a `threshold` number of parts are available.
func Combine(parts [][]byte) ([]byte, error) {
	// Verify enough parts provided
	if len(parts) < 2 {
		return nil, fmt.Errorf("less than two parts cannot be used to reconstruct the secret")
	}

	// Verify the parts are all the same length
	firstPartLen := len(parts[0])
	if firstPartLen < 2 {
		return nil, fmt.Errorf("parts must be at least two bytes")
	}
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) != firstPartLen {
			return nil, fmt.Errorf("all parts must be the same length")
		}
	}

	// Create a buffer to store the reconstructed secret
	secret := make([]byte, firstPartLen-1)

	// Buffer to store the samples
	xSamples := make([]uint8, len(parts))
	ySamples := make([]uint8, len(parts))

	// Set the x value for each sample and ensure no x_sample values are the same,
	// otherwise div() can be unhappy
	// Check that we don't have any duplicate parts, that is, two or
	// more parts with the same x coordinate.
	checkMap := map[byte]bool{}
	for i, part := range parts {
		samp := part[firstPartLen-1]
		if exists := checkMap[samp]; exists {
			return nil, fmt.Errorf("duplicate part detected")
		}
		checkMap[samp] = true
		xSamples[i] = samp
	}

	// Reconstruct each byte
	for idx := range secret {
		// Set the y value for each sample
		for i, part := range parts {
			ySamples[i] = part[idx]
		}

		// Use Lagrange interpolation to retrieve the free term
		// of the original polynomial
		val := interpolatePolynomial(xSamples, ySamples, 0)

		// Evaluate the 0th value to get the intercept
		secret[idx] = val
	}
	return secret, nil
}
