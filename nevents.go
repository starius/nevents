package nevents

import "math/big"

// ProbabilityAtLeastN calculates probability of at least n of events happening.
// It accepts n and the list of probabilities of individual events. Events are
// independent. The code uses dynamic programming to calculate probabilities
// of exactly n events happening, then sums such probabilities for k>=n. It has
// complexity of O(N^2).
func ProbabilityAtLeastN(n int, probabilities []big.Float) big.Float {
	// exactProbs stores probabilities of exactly i events happening. It is
	// built dynamically, using dynamic programming approach.
	exactProbs := make([]big.Float, len(probabilities)+1)

	// We start DP for 0 events happening. Fill the element 0 of exactProbs
	// with 1, since with probability 1 there are 0 out of 0 events
	// happening.
	exactProbs[0].SetFloat64(1)

	// Go through all the events, adding them to solution.
	for i, pi := range probabilities {
		// Update exactProbs elements from (i+1) to 0. We do it in the
		// opposite direction, because each element depends on itself
		// and the previous element:
		//        exactProbs[5]      exactProbs[6]      exactProbs[7]
		// (i=5)  p55                p56                0
		// (i=6)  p55*(1-p6)+p54*p6  p56*(1-p6)+p55*p6  p56 * p6

		// Start with the last non-zero element, it is the probability
		// that all events up to i-th happen.
		exactProbs[i+1].Mul(&exactProbs[i], &pi)

		// Find probability of i-th event not happening:
		var qi big.Float
		qi.SetFloat64(1)
		qi.Sub(&qi, &pi)

		// Now visit all the elements in the middle (not last, not 0).
		for j := i; j > 0; j-- {
			var withoutI, withI big.Float
			// j events can happen two ways: either i-th event
			// doesn't happen and j events happened in the previous
			// generation:
			withoutI.Mul(&qi, &exactProbs[j])

			// ... or i-th event happens and j-1 events happened in
			// the previous generation:
			withI.Mul(&pi, &exactProbs[j-1])

			// Now add these probabilities:
			exactProbs[j].Add(&withoutI, &withI)
		}

		// Complete with updating the probability that 0 events happen.
		exactProbs[0].Mul(&qi, &exactProbs[0])
	}

	// Now we have probabilities of exactly i events happening among all the
	// events. To find the probability of >=n events happening, we need to
	// calculate the sum in this array from n to the end.
	var pn big.Float
	for _, p := range exactProbs[n:] {
		pn.Add(&pn, &p)
	}

	return pn
}