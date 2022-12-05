package overlay

import (
	"io"
	"sort"
)

// Interval represents a 2er set of integers modelling a range.
type Interval interface {
	// Range returns the minimum and maximum of the interval.
	// Minimum value is inclusive, maximum value exclusive.
	// In other notation: [min, max)
	Range() (int64, int64)

	// Merge merges the interval `i` to this interval.
	// The range borders should be fixed accordingly,
	// so that [min(i.min, self.min), max(i.max, self.max)] applies.
	Merge(i Interval)
}

// Modification represents a single write
type Modification struct {
	// Offset where the modification started:
	offset int64

	// Data that was changed:
	// This might be changed to a mmap'd byte slice later.
	data []byte
}

// Range returns a fitting integer interval
func (n *Modification) Range() (int64, int64) {
	return n.offset, n.offset + int64(len(n.data))
}

// Merge adds the data of another interval where they intersect.
// The overlapping parts are taken from `n` always.
// Note: `i` shall not be used after calling Merge.
func (n *Modification) Merge(i Interval) {
	// interracial merges are forbidden
	other, ok := i.(*Modification)
	if !ok {
		return
	}

	oMin, oMax := other.Range()
	nMin, nMax := n.Range()

	// check if the intervals overlap
	// if not, there's nothing left to do
	if oMin > nMax || nMin > oMax {
		return
	}

	// Prepend non-overlapping data from `other`:
	if nMin > oMin {
		n.data = append(other.data[:nMin-oMin], n.data...)
		n.offset = other.offset
	}

	// Append non-overlapping data from `other`:
	if nMax < oMax {
		// Append other.data[(other.Max - n.Max):]
		n.data = append(n.data, other.data[(oMax-nMax-1):]...)
	}

	// Make sure old data gets invalidated quickly:
	other.data = nil
}

// IntervalIndex represents a continuous array of sorted intervals.
// When adding intervals to the index, it will merge them overlapping areas.
// Holes between the intervals are allowed.
type IntervalIndex struct {
	r []Interval

	// Max is the maximum interval offset given to Add()
	Max int64
}

func (ivl *IntervalIndex) Add(n Interval) {
	Min, Max := n.Range()
	if Max < Min {
		panic("Max > Min!")
	}

	// Initial case: Add as single element.
	if ivl.r == nil {
		ivl.r = []Interval{n}
		ivl.Max = Max
		return
	}

	// find the lowest fitting interval
	minIdx := sort.Search(len(ivl.r), func(i int) bool {
		_, iMax := ivl.r[i].Range()
		return Min <= iMax
	})

	// find the highest fitting interval
	maxIdx := sort.Search(len(ivl.r), func(i int) bool {
		iMin, _ := ivl.r[i].Range()
		return Max <= iMin
	})

	// remember the biggest offset:
	if Max > ivl.Max {
		ivl.Max = Max
	}

	// new interval is bigger than all others:
	if minIdx >= len(ivl.r) {
		ivl.r = append(ivl.r, n)
		return
	}

	// new range fits nicely in; just insert it in between:
	if minIdx == maxIdx {
		ivl.r = insert(ivl.r, minIdx, n)
		return
	}

	// somewhere in between. Merge to continuous interval:
	for i := minIdx; i < maxIdx; i++ {
		n.Merge(ivl.r[i])
	}

	// delete old unmerged intervals and substitute with the merged:
	ivl.r[minIdx] = n
	ivl.r = cut(ivl.r, minIdx+1, maxIdx)

}

// cut deletes the a[i:j] from a and returns the new slice.
func cut(a []Interval, i, j int) []Interval {
	copy(a[i:], a[j:])
	for k, n := len(a)-j+i, len(a); k < n; k++ {
		a[k] = nil
	}
	return a[:len(a)-j+i]
}

// insert squeezes `x` at a[i] and moves the reminding elements.
// Returns the modified slice.
func insert(a []Interval, i int, x Interval) []Interval {
	a = append(a, nil)
	copy(a[i+1:], a[i:])
	a[i] = x
	return a
}

// Overlays returns all intervals that intersect with [start, end)
func (ivl *IntervalIndex) Overlays(start, end int64) []Interval {
	// Find the lowest matching interval:
	lo := sort.Search(len(ivl.r), func(i int) bool {
		_, iMax := ivl.r[i].Range()
		return start <= iMax
	})

	hi := sort.Search(len(ivl.r), func(i int) bool {
		iMin, _ := ivl.r[i].Range()
		return end <= iMin
	})

	return ivl.r[lo:hi]
}

// min returns the minimum of a and b.
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of a and b.
func max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

// Layer is an io.ReadWriter that takes an underlying Reader
// and caches Writes on top of it. To the outside it delivers
// a zipped stream of the recent writes and the underlying stream.
type Layer struct {
	index    *IntervalIndex
	r        io.ReadSeeker
	pos      int64
	limit    int64
	fileSize int64
}

// NewLayer returns a new in memory layer.
// No IO is performed on creation.
func NewLayer(r io.ReadSeeker) *Layer {
	return &Layer{
		index:    &IntervalIndex{},
		r:        r,
		limit:    -1,
		fileSize: -1,
	}
}

// SetSize sets the size of the absolute layer.
func (l *Layer) SetSize(size int64) {
	l.fileSize = size
}
