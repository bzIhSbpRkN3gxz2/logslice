package index

import (
	"testing"
	"time"
)

func ts(sec int64) time.Time {
	return time.Unix(sec, 0).UTC()
}

func TestMergeDisjoint(t *testing.T) {
	a := Index{Entries: []Entry{
		{Time: ts(100), Offset: 0},
		{Time: ts(200), Offset: 50},
	}}
	b := Index{Entries: []Entry{
		{Time: ts(300), Offset: 100},
		{Time: ts(400), Offset: 150},
	}}

	got := Merge(a, b)
	if len(got.Entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(got.Entries))
	}
	for i := 1; i < len(got.Entries); i++ {
		if !got.Entries[i].Time.After(got.Entries[i-1].Time) {
			t.Errorf("entries not sorted at index %d", i)
		}
	}
}

func TestMergeDuplicateTimestampTakesHigherOffset(t *testing.T) {
	a := Index{Entries: []Entry{{Time: ts(100), Offset: 10}}}
	b := Index{Entries: []Entry{{Time: ts(100), Offset: 99}}}

	got := Merge(a, b)
	if len(got.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got.Entries))
	}
	if got.Entries[0].Offset != 99 {
		t.Errorf("expected offset 99, got %d", got.Entries[0].Offset)
	}
}

func TestMergeEmptyA(t *testing.T) {
	a := Index{}
	b := Index{Entries: []Entry{{Time: ts(1), Offset: 5}}}
	got := Merge(a, b)
	if len(got.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got.Entries))
	}
}

func TestMergeEmptyB(t *testing.T) {
	a := Index{Entries: []Entry{{Time: ts(1), Offset: 5}}}
	b := Index{}
	got := Merge(a, b)
	if len(got.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got.Entries))
	}
}

func TestMergeResultIsSorted(t *testing.T) {
	a := Index{Entries: []Entry{
		{Time: ts(300), Offset: 200},
		{Time: ts(100), Offset: 0},
	}}
	b := Index{Entries: []Entry{
		{Time: ts(200), Offset: 100},
	}}

	got := Merge(a, b)
	if len(got.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got.Entries))
	}
	for i := 1; i < len(got.Entries); i++ {
		if !got.Entries[i].Time.After(got.Entries[i-1].Time) {
			t.Errorf("entries not sorted at index %d", i)
		}
	}
}
