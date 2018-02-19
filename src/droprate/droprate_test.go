package droprate

import (
	"testing"
)

// Test cases taken from
// https://www.reddit.com/r/runescape/comments/736utj/aod_and_telos_drop_rates_revealed/
func TestDropRate(t *testing.T) {
	tests := []struct {
		enrage int64
		streak int64
		lotd   bool
		want   int64
	}{
		{
			100,
			1,
			false,
			263,
		},
		{
			100,
			1,
			true,
			226,
		},
		{
			500,
			1,
			false,
			72,
		},
		{
			200,
			5,
			false,
			133,
		},
		{
			999,
			1,
			false,
			38,
		},
		{
			1024,
			1,
			false,
			37,
		},
	}

	for _, test := range tests {
		if got, err := DropRate(test.enrage, test.streak, test.lotd); got != test.want || err != nil {
			t.Errorf("DropRate(%v, %v, %v)= %v, %v want= %v, nil", test.enrage, test.streak, test.lotd, got, err, test.want)
		}
	}
}

func TestDropRateErrors(t *testing.T) {
	tests := []struct {
		name   string
		enrage int64
		streak int64
		lotd   bool
	}{
		{
			"Negative enrage",
			-5,
			5,
			false,
		},
		{
			"Negative streak",
			5,
			-5,
			false,
		},
		{
			"Over 4000 enrage",
			4010,
			5,
			false,
		},
	}

	for _, test := range tests {
		if got, err := DropRate(test.enrage, test.streak, test.lotd); err == nil {
			t.Errorf("DropRate(%v)= %v, %v, want= 0, err", test.name, got, err)
		}
	}
}
