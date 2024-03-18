package timecode_test

import (
	"testing"

	"github.com/spiretechnology/go-timecode"
)

func TestParse_NDF(t *testing.T) {
	cases := map[string]int64{
		"00:00:00:00": 0,
		"00:00:00:01": 1,
		"00:00:01:01": 61,
		"00:00:11:01": 661,
	}
	for k, f := range cases {
		tc, _ := timecode.Parse(k, timecode.Rate_60)
		if frame := tc.Frame(); frame != f {
			t.Errorf("Timecode %s should be equivalent to frame %d. Got %d\n", k, f, frame)
		} else {
			t.Logf("Success, timecode %s equals frame %d\n", k, f)
		}
	}
}

func TestParse_DF(t *testing.T) {
	cases := map[string]int64{
		"00:02:00;02": 3598,
		"00:01:59;23": 3591,
		"00:01:59;22": 3590,
		"00:03:00;04": 5398,
	}
	for k, f := range cases {
		tc, _ := timecode.Parse(k, timecode.Rate_29_97)
		if frame := tc.Frame(); frame != f {
			t.Errorf("Timecode %s should be equivalent to frame %d. Got %d\n", k, f, frame)
		} else {
			t.Logf("Success, timecode %s equals frame %d\n", k, f)
		}
	}
}

func TestParseInvalidDF(t *testing.T) {
	cases := map[string]string{
		"00:01:59;23": "00:01:59;23",
		"00:02:00;00": "00:02:00;02",
		"00:02:00;01": "00:02:00;02",
		"00:02:00;02": "00:02:00;02",
		"00:02:00;03": "00:02:00;03",
	}
	for k, s := range cases {
		tc, _ := timecode.Parse(k, timecode.Rate_29_97)
		t.Logf("%s => %d\n", k, tc.Frame())
		if str := tc.String(); str != s {
			t.Errorf("DF timecode %s should be rounded to timecode %s. Got %s\n", k, s, str)
		} else {
			t.Logf("Success, DF timecode %s rounded to %s\n", k, s)
		}
	}
}

// TestParseInvalidDFInNDFRate makes sure that the above drop frame rounding behavior
// doesn't affect parsing NDF timecodes.
func TestParseInvalidDFInNDFRate(t *testing.T) {
	cases := []string{
		"00:01:59:23",
		"00:02:00:00",
		"00:02:00:01",
		"00:02:00:02",
		"00:02:00:03",
	}
	for _, s := range cases {
		tc, _ := timecode.Parse(s, timecode.Rate_24)
		if str := tc.String(); str != s {
			t.Errorf("NDF timecode %s should NOT be rounded. Got %s\n", s, str)
		} else {
			t.Logf("Success, NDF timecode %s stayed the same\n", s)
		}
	}
}
