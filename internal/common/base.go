package common

import "time"

type DateRange struct {
	Start time.Time
	End   time.Time
}

type OverlapType int

const (
	NoOverlap OverlapType = iota
	Overlap
	Enclosed
	Encloses
)

func (d DateRange) Overlap(other DateRange) OverlapType {
	if d.End.Before(other.Start) || other.End.Before(d.Start) {
		return NoOverlap
	} else if d.Start.Before(other.Start) && d.End.Before(other.End) {
		return Enclosed
	} else if other.Start.Before(d.Start) && other.End.Before(d.End) {
		return Encloses
	} else {
		return Overlap
	}
}
