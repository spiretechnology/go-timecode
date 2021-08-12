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
