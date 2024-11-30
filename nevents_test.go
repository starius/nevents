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
			// 1 - 0.9*0.9.
			want: "0.19",
		},
		{
			probabilities: []string{"0.1", "0.1"},
			n:             0,
			want:          "1",
		},
		{
			probabilities: []string{"0.1", "0.1"},
			n:             2,
			// 0.1*0.1.
			want: "0.01",
		},
		{
			probabilities: []string{"0.1", "0.1"},
			n:             3,
			want:          "0",
		},
		{
			probabilities: []string{"0.1", "0.1"},
			n:             10,
			want:          "0",
		},
		{
			probabilities: []string{"0.1", "0.1"},
			n:             -1,
			want:          "1",
		},
		{
			probabilities: []string{"0.1", "0.1"},
			n:             -100,
			want:          "1",
		},
		{
			probabilities: []string{"0.1", "0.1", "0.2"},
			n:             1,
			// 1 - 0.9*0.9*0.8.
			want: "0.352",
		},
		{
			probabilities: []string{"0.1", "0.1", "0.2"},
			n:             2,
			// 1 - 0.9*0.9*0.8 - 2*0.1*0.9*0.8 - 0.2*0.9*0.9.
			want: "0.046",
		},
		{
			probabilities: []string{"0.1", "0.1", "0.2"},
			n:             3,
			// 0.1*0.1*0.2.
			want: "0.002",
		},
		{
			probabilities: []string{
				"0.1", "0.1", "0.1", "0.1", "0.1",
				"0.1", "0.1", "0.1", "0.1", "0.1",
			},
			n: 4,
			// 1 - 0.9**10 - 10 * 0.1 * 0.9**9 -
			// - 10*9/2 * 0.1**2 * 0.9**8 -
			// - 10*9*8/(2*3) * 0.1**3 * 0.9**7
			want: "0.0127951984",
		},
		{
			probabilities: []string{
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001",
			},
			n: 4,
			// 1 - 0.999999**10 - 10 * 0.000001 * 0.999999**9 -
			// - 10*9/2 * 0.000001**2 * 0.999999**8 -
			// - 10*9*8/(2*3) * 0.000001**3 * 0.999999**7
			want: "2.09998992e-22",
		},
		{
			probabilities: []string{
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
			},
			n: 30,
			// sum(c(40,i) * 0.000001^i * 0.999999^(40-i)) i=30 to 40
			want: "8.476523249e-172",
		},
		{
			probabilities: []string{
				"0.01", "0.01", "0.01", "0.01", "0.01",
				"0.01", "0.01", "0.01", "0.01", "0.01",
				"0.01", "0.01", "0.01", "0.01", "0.01",
				"0.01", "0.01", "0.01", "0.01", "0.01",
				"0.01", "0.01", "0.01", "0.01", "0.01",
				"0.01", "0.01", "0.01", "0.01", "0.01",
				"0.01", "0.01", "0.01", "0.01", "0.01",
				"0.01", "0.01", "0.01", "0.01", "0.01",
			},
			n: 30,
			// sum(c(40,i) * 0.01^i * 0.99^(40-i)) i=30 to 40
			want: "7.691140123e-52",
		},
		{
			probabilities: []string{
				"0.5", "0.5", "0.5", "0.5", "0.5",
				"0.5", "0.5", "0.5", "0.5", "0.5",
				"0.5", "0.5", "0.5", "0.5", "0.5",
				"0.5", "0.5", "0.5", "0.5", "0.5",
				"0.5", "0.5", "0.5", "0.5", "0.5",
				"0.5", "0.5", "0.5", "0.5", "0.5",
				"0.5", "0.5", "0.5", "0.5", "0.5",
				"0.5", "0.5", "0.5", "0.5", "0.5",
			},
			n: 30,
			// sum(c(40,i) * 0.5^40) i=30 to 40
			want: "0.001110716887",
		},
		{
			probabilities: []string{
				"0.9", "0.9", "0.9", "0.9", "0.9",
				"0.9", "0.9", "0.9", "0.9", "0.9",
				"0.9", "0.9", "0.9", "0.9", "0.9",
				"0.9", "0.9", "0.9", "0.9", "0.9",
				"0.9", "0.9", "0.9", "0.9", "0.9",
				"0.9", "0.9", "0.9", "0.9", "0.9",
				"0.9", "0.9", "0.9", "0.9", "0.9",
				"0.9", "0.9", "0.9", "0.9", "0.9",
			},
			n: 30,
			// sum(c(40,i) * 0.9^i * 0.1^(40-i)) i=30 to 40
			want: "0.9985302775",
		},
		{
			probabilities: []string{
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
			},
			n: 30,
			// sum(c(120,i) * 0.000001^i * 0.999999^(120-i))
			// i=30 to 120
			want: "1.69730604e-152",
		},
		{
			probabilities: []string{
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
			},
			n: 41,
			// sum(c(120,i) * 0.000001^i * 0.999999^(120-i))
			// i=41 to 120
			want: "2.235083201e-214",
		},
		{
			probabilities: []string{
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",

				"1", "1", "1", "1", "1", "1", "1", "1",
				"1", "1", "1", "1", "1", "1", "1", "1",
				"1", "1", "1", "1", "1", "1", "1", "1",
				"1", "1", "1", "1", "1", "1", "1", "1",
				"1", "1", "1", "1", "1", "1", "1", "1",
			},
			n:    40,
			want: "1",
		},
		{
			probabilities: []string{
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",
				"0.000001", "0.000001", "0.000001", "0.000001",

				"1", "1", "1", "1", "1", "1", "1", "1",
				"1", "1", "1", "1", "1", "1", "1", "1",
				"1", "1", "1", "1", "1", "1", "1", "1",
				"1", "1", "1", "1", "1", "1", "1", "1",
				"1", "1", "1", "1", "1", "1", "1", "1",
			},
			n:    41,
			// 1 - (1-0.000001)**80.
			want: "7.999684008e-05",
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

func BenchmarkProbabilityAtLeastN(b *testing.B) {
	const (
		M = 120
		N = 41
	)

	var p big.Float
	p.SetFloat64(0.000001)

	probabilities := make([]big.Float, M)
	for i := range probabilities {
		probabilities[i] = p
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ProbabilityAtLeastN(N, probabilities)
	}
}
