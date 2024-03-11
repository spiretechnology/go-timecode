package timecode_test

import (
	"testing"

	"github.com/spiretechnology/go-timecode"
	"github.com/stretchr/testify/require"
)

func TestRateFromFraction(t *testing.T) {
	type testCase struct {
		num, den int
		rate     timecode.Rate
	}
	cases := []testCase{
		{24000, 1001, timecode.Rate_23_976},
		{24, 1, timecode.Rate_24},
		{30, 1, timecode.Rate_30},
		{30000, 1001, timecode.Rate_29_97},
		{60, 1, timecode.Rate_60},
		{60000, 1001, timecode.Rate_59_94},
	}
	for _, tc := range cases {
		rate := timecode.RateFromFraction(tc.num, tc.den)
		require.Equal(t, tc.rate, rate)
	}
}

func TestRateFromFractionBuiltinRates(t *testing.T) {
	cases := []timecode.Rate{
		timecode.Rate_23_976,
		timecode.Rate_24,
		timecode.Rate_30,
		timecode.Rate_29_97,
		timecode.Rate_60,
		timecode.Rate_59_94,
	}
	for _, rate := range cases {
		newRate := timecode.RateFromFraction(rate.Num, rate.Den)
		require.Equal(t, rate, newRate)
	}
}
