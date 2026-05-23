package index

import (
	"testing"
	"time"
)

func buildIndex(entries []Entry) *Index {
	return &Index{Entries: entries}
}

func TestSearchBothBounds(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	idx := buildIndex([]Entry{
		{Time: base, Offset: 0},
		{Time: base.Add(10 * time.Minute), Offset: 1000},
		{Time: base.Add(20 * time.Minute), Offset: 2000},
		{Time: base.Add(30 * time.Minute), Offset: 3000},
	})

	result := Search(idx, base.Add(5*time.Minute), base.Add(25*time.Minute))

	if result.StartOffset != 0 {
		t.Errorf("StartOffset: got %d, want 0", result.StartOffset)
	}
	if result.EndOffset != 3000 {
		t.Errorf("EndOffset: got %d, want 3000", result.EndOffset)
	}
}

func TestSearchUnboundedStart(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	idx := buildIndex([]Entry{
		{Time: base, Offset: 0},
		{Time: base.Add(10 * time.Minute), Offset: 500},
	})

	result := Search(idx, time.Time{}, base.Add(5*time.Minute))

	if result.StartOffset != 0 {
		t.Errorf("StartOffset: got %d, want 0", result.StartOffset)
	}
	if result.EndOffset != 500 {
		t.Errorf("EndOffset: got %d, want 500", result.EndOffset)
	}
}

func TestSearchUnboundedEnd(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	idx := buildIndex([]Entry{
		{Time: base, Offset: 0},
		{Time: base.Add(10 * time.Minute), Offset: 800},
	})

	result := Search(idx, base.Add(5*time.Minute), time.Time{})

	if result.StartOffset != 0 {
		t.Errorf("StartOffset: got %d, want 0", result.StartOffset)
	}
	if result.EndOffset != -1 {
		t.Errorf("EndOffset: got %d, want -1", result.EndOffset)
	}
}

func TestSearchEmptyIndex(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	idx := buildIndex(nil)

	result := Search(idx, base, base.Add(time.Hour))

	if result.StartOffset != 0 {
		t.Errorf("StartOffset: got %d, want 0", result.StartOffset)
	}
	if result.EndOffset != -1 {
		t.Errorf("EndOffset: got %d, want -1", result.EndOffset)
	}
}

func TestSearchEndPastAllEntries(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	idx := buildIndex([]Entry{
		{Time: base, Offset: 0},
		{Time: base.Add(10 * time.Minute), Offset: 1000},
	})

	result := Search(idx, base, base.Add(time.Hour))

	if result.EndOffset != -1 {
		t.Errorf("EndOffset: got %d, want -1 (read to EOF)", result.EndOffset)
	}
}
