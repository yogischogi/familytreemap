// Package provides statistical functions for advanced users.
package statistics

import (
	"math"
)

// AcceptableSampleSize returns the minimal sample size
// that is necessary to approximate the hypergeometric
// distribution by a normal distribution.
//
// The approximation will not be very good but we will accept
// it because the results should still make some sense.
//
// p: Calculated probability from a sample
// = (m/n: positive results/sample size),
// http://de.wikipedia.org/wiki/Normalverteilung
func AcceptableSampleSize(p float64) float64 {
	var n float64
	n1 := 4 / p
	n2 := 4 / (1 - p)
	if n1 >= n2 {
		n = n1
	} else {
		n = n2
	}
	return math.Floor(n + 1)
}

// MinSampleSize returns the minimal sample size
// that is necessary to approximate the hypergeometric
// distribution by a normal distribution.
//
// For more information:
// Marcus Hudec, Christian Neumann:
// Stichproben und Umfragen,
// http://www.stat4u.at/download/1423/stichpr.pdf
// Or WikiPedia:
// http://de.wikipedia.org/wiki/Normalverteilung
func MinSampleSize(p float64) float64 {
	return math.Floor(9/(p*(1-p)) + 1)
}

// Probability returns the probabilty and confidence
// interval (95%) for m positive samples out of a
// sample size of n individuals.
//
// More information:
// Marcus Hudec, Christian Neumann:
// Stichproben und Umfragen,
// http://www.stat4u.at/download/1423/stichpr.pdf
// Burt Gerstman: StatPrimer,
// http://www.sjsu.edu/faculty/gerstman/StatPrimer/
// http://www.sjsu.edu/faculty/gerstman/StatPrimer/estimation.pdf
func Probability(n float64, m float64) (p float64, s float64) {
	p = m / n
	// Quantil z = 1 - a/2 for 95% conficence.
	var z float64 = 1.96
	s = z * math.Sqrt(p*(1-p)/n)
	return p, s
}
