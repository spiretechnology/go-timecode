package timecode_test

import (
	"testing"

	"github.com/spiretechnology/go-timecode"
)

func TestTimecode_FrameToString_DF(t *testing.T) {
	cases := map[int64]string{
		2878: "00:02:00;02",
	}
	for f, tcode := range cases {
		tc := timecode.FromFrame(f, timecode.Rate_23_976)
		if str := tc.String(); str != tcode {
			t.Errorf("Frame %d should be equivalent to timecode %s. Got %s\n", f, tcode, str)
		} else {
			t.Logf("Success, frame %d equals timecode %s\n", f, tcode)
		}
	}
}

func TestTimecode_Identity_DF(t *testing.T) {
	cases := []string{
		"00:02:00;02",
		"00:00:00;00",
		"00:00:59;23",
		"00:01:00;02",
		"00:03:59;23",
		"00:04:00;02",
		"00:01:59;23",
		"00:09:59;23",
		"00:10:00;00",
	}
	for _, tcode := range cases {
		tc, _ := timecode.Parse(tcode, timecode.Rate_23_976)
		if str := tc.String(); str != tcode {
			t.Errorf("Timecode %s became %s during parsing and printing\n", tcode, str)
		} else {
			t.Logf("Success, identity valid for %s\n", tcode)
		}
	}
}

func TestTimecode_AddOne_DF(t *testing.T) {
	sequences := map[string]string{
		"00:00:59;23": "00:01:00;02",
		"00:03:59;23": "00:04:00;02",
		"00:01:59;23": "00:02:00;02",
		"00:09:59;23": "00:10:00;00",
	}
	for fromTC, toTC := range sequences {
		tc, _ := timecode.Parse(fromTC, timecode.Rate_23_976)
		next := tc.Add(timecode.Frame(1))
		if str := next.String(); str != toTC {
			t.Errorf("Expected %s => %s, got %s\n", fromTC, toTC, str)
		} else {
			t.Logf("Success, got %s => %s\n", fromTC, toTC)
		}
	}
}

func bruteForceAdd1(c timecode.Components) timecode.Components {
	c.Frames++
	if c.Frames >= 24 {
		c.Frames -= 24
		c.Seconds++
		if c.Seconds >= 60 {
			c.Seconds -= 60
			c.Minutes++
			if c.Minutes >= 60 {
				c.Minutes -= 60
				c.Hours++
			}
		}
	}
	return c
}

// TestTimecodeSequenceNDF jumps to a starting point and then seeks through the frames 1 by 1 to make sure
// the generated timecodes match what the result would be if we brute forced. Brute forcing is much slower
// if we're adding multiple frames, but adding just 1 frame allows us to put it up head-to-head against
// our timecode implementation to ensure correctness.
func TestTimecodeSequenceNDF(t *testing.T) {
	startTimecodes := map[string]int{
		"00:00:00:00": 100000,
		"03:59:59:00": 100000,
	}
	for startTimecodeStr, iterations := range startTimecodes {
		prevTc, _ := timecode.Parse(startTimecodeStr, timecode.Rate_24)
		prevComp := prevTc.Components()

		// Run through all the iterations for this sample
		for i := 0; i < iterations; i++ {
			tc := prevTc.Add(timecode.Frame(1))
			comp := tc.Components()
			expectedComp := bruteForceAdd1(prevComp)
			if !comp.Equals(expectedComp) {
				t.Errorf("Add 1 frame, skipped from %s to %s\n", prevTc.String(), tc.String())
				// } else {
				// 	t.Logf("Success, %s + 1 frame = %s", prevTc.String(), tc.String())
			}
			prevTc = tc
			prevComp = comp
		}
	}
}
