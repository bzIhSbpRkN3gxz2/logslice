// Package filter provides time-range filtering for log lines.
//
// A Filter is constructed with a Range specifying an inclusive Start
// and exclusive End boundary. Either boundary may be zero to indicate
// an open (unbounded) range on that side.
//
// Typical usage:
//
//	f, err := filter.New(filter.Range{
//		Start: startTime,
//		End:   endTime,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, line := range lines {
//		if f.Past(line.Time) {
//			break // sorted input: no more matches possible
//		}
//		if f.Match(line.Time) {
//			output(line)
//		}
//	}
package filter
