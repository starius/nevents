package nevents

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProbabilityAtLeastN(t *testing.T) {
	cases := []struct {
		probabilities []string
		n             int
		want          string
	}{
		{
			probabilities: []string{"0.1", "0.1"},
			n:             1,
			want:          "0.19",
		},
		{
			probabilities: []string{"0.1", "0.1"},
			n:             0,
			want:          "1",
		},
		{
			probabilities: []string{"0.1", "0.1"},
			n:             2,
			want:          "0.01",
		},
		{
			probabilities: []string{"0.1", "0.1", "0.2"},
			n:             1,
			want:          "0.352",
		},
		{
			probabilities: []string{"0.1", "0.1", "0.2"},
			n:             2,
			want:          "0.046",
		},
		{
			probabilities: []string{"0.1", "0.1", "0.2"},
			n:             3,
			want:          "0.002",
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf(
			"%d of %s", tc.n, strings.Join(tc.probabilities, ","),
		)
		t.Run(name, func(t *testing.T) {
			probs := make([]big.Float, len(tc.probabilities))
			for i, str := range tc.probabilities {
				_, _, err := probs[i].Parse(str, 10)
				require.NoError(t, err)
			}
			got := ProbabilityAtLeastN(tc.n, probs)
			require.Equal(t, tc.want, got.String())
		})
	}
}
