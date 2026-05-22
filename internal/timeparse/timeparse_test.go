package timeparse_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/timeparse"
)

func TestParserKnownFormats(t *testing.T) {
	cases := []struct {
		input    string
		wantYear int
		wantHour int
	}{
		{"2024-03-15T14:22:01Z", 2024, 14},
		{"2024-03-15T14:22:01.123456789Z", 2024, 14},
		{"2024-03-15 14:22:01", 2024, 14},
		{"2024-03-15 14:22:01.999", 2024, 14},
		{"15/Mar/2024:14:22:01 +0000", 2024, 14},
	}

	p := timeparse.NewParser()
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := p.Parse(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Year() != tc.wantYear {
				t.Errorf("year: got %d, want %d", got.Year(), tc.wantYear)
			}
			if got.Hour() != tc.wantHour {
				t.Errorf("hour: got %d, want %d", got.Hour(), tc.wantHour)
			}
		})
	}
}

func TestParserCachesFormat(t *testing.T) {
	p := timeparse.NewParser()
	inputs := []string{
		"2024-01-01T00:00:00Z",
		"2024-06-15T12:30:00Z",
		"2024-12-31T23:59:59Z",
	}
	for _, raw := range inputs {
		if _, err := p.Parse(raw); err != nil {
			t.Fatalf("parse %q: %v", raw, err)
		}
	}
}

func TestParserUnknownFormat(t *testing.T) {
	p := timeparse.NewParser()
	_, err := p.Parse("not-a-timestamp")
	if err == nil {
		t.Fatal("expected error for unrecognized format, got nil")
	}
}

func TestParseWithFormat(t *testing.T) {
	raw := "2024-03-15T14:22:01Z"
	got, err := timeparse.ParseWithFormat(raw, time.RFC3339)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Month() != time.March {
		t.Errorf("expected March, got %v", got.Month())
	}
}

func TestParseWithFormatInvalid(t *testing.T) {
	_, err := timeparse.ParseWithFormat("garbage", time.RFC3339)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
