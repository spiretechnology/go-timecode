package timecode_test

import (
	"testing"
	"time"

	"github.com/spiretechnology/go-timecode"
)

func TestParse_Zero(t *testing.T) {
	tc, _ := timecode.Parse("00:00:00:00", timecode.Rate_23_976)
	if tc.Frame() != 0 {
		t.Error("Timecode 00:00:00:00 should be equivalent to frame 0")
	} else {
		t.Log("Success, timecode 00:00:00:00 equals frame 0")
	}
}

func TestParse_One(t *testing.T) {
	tc, _ := timecode.Parse("00:00:00:01", timecode.Rate_23_976)
	if tc.Frame() != 1 {
		t.Error("Timecode 00:00:00:01 should be equivalent to frame 1")
	} else {
		t.Log("Success, timecode 00:00:00:01 equals frame 1")
	}
}

func TestParse_DF(t *testing.T) {
	cases := map[string]int64{
		"00:02:00;02": 2878,
		"00:01:59;23": 2877,
		"00:01:59;22": 2876,
		"00:03:00;04": 4318,
	}
	for k, f := range cases {
		tc, _ := timecode.Parse(k, timecode.Rate_23_976)
		if frame := tc.Frame(); frame != f {
			t.Errorf("Timecode %s should be equivalent to frame %d. Got %d\n", k, f, frame)
		} else {
			t.Logf("Success, timecode %s equals frame %d\n", k, f)
		}
	}
}

func TestFromPresentationTime(t *testing.T) {
	cases := map[time.Duration]string{
		time.Minute * 2:             "00:02:00:00",
		time.Minute * 10:            "00:10:00:00",
		time.Second * 6:             "00:00:06:00",
		time.Hour*6 + time.Second/2: "06:00:00:12",
	}
	for pt, s := range cases {
		tc := timecode.FromPresentationTime(pt, timecode.Rate_24)
		if str := tc.String(); str != s {
			t.Errorf("Presentation time %s should be timecode %s. Got %s\n", pt.String(), s, str)
		} else {
			t.Logf("Success, presentation time %s equals timecode %s\n", pt.String(), s)
		}
	}
}
