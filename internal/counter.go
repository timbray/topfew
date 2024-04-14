package topfew

import (
	"math/bits"
	"slices"
)

// KeyCount represents a key's occurrence count.
type KeyCount struct {
	Key   string
	Count *uint64
}

// Counter represents a bunch of keys and their occurrence counts, with the highest counts tracked.
// threshold represents the minimum count value to qualify for consideration as a top count
// the "top" map represents the keys & counts encountered so far which are higher than threshold
// The hash values are pointers not integers for efficiency reasons, so you don't have to update the
// map[string] mapping, you just update the number the key maps to.
type Counter struct {
	counts map[string]uint64
	size   int
}

// NewCounter creates a new empty counter, ready for use. Size controls how many top items to track.
func NewCounter(size int) *Counter {
	t := new(Counter)
	t.size = size
	t.counts = make(map[string]uint64, 1024)
	return t
}

// Add one occurrence to the counts for the indicated key.
func (t *Counter) Add(bytes []byte) {
	// note the call with a byte slice rather than the string because of
	//  https://github.com/golang/go/commit/f5f5a8b6209f84961687d993b93ea0d397f5d5bf
	//  which recognizes the idiom foo[string(someByteSlice)] and bypasses constructing the string;
	//  of course we'd rather just say foo[someByteSlice] but that's not legal because Reasons.

	// have we seen this key?
	t.counts[string(bytes)]++
}

// GetTop returns the top occurring keys & counts in order of descending count
func (t *Counter) GetTop() []*KeyCount {
	keys := make([]string, 0, len(t.counts))
	packedCounts := make([]uint64, 0, len(t.counts))

	// calculate how many bits of the counter space should be dedicated to key space
	shift := 64 - bits.LeadingZeros(uint(len(t.counts)))
	mask := uint64(1<<shift - 1)

	id := 0
	for key, count := range t.counts {
		keys = append(keys, key)
		packed := (count << shift) + uint64(id)
		packedCounts = append(packedCounts, packed)
		id++
	}

	slices.Sort(packedCounts)

	results := make([]*KeyCount, 0, t.size)
	start := len(packedCounts) - 1
	end := start - t.size
	// clamp in case we don't have t.size top values
	if end < -1 {
		end = -1
	}
	for i := start; i > end; i-- {
		packed := packedCounts[i]
		count := packed >> shift
		key := keys[packed&mask]
		results = append(results, &KeyCount{key, &count})
	}
	return results
}

// merge applies the counts from the SegmentCounter into the Counter.
// Once merged, the SegmentCounter should be discarded.
func (t *Counter) merge(segCounter segmentCounter) {
	for segKey, segCount := range segCounter {
		t.counts[segKey] += segCount
	}
}

// SegmentCounter tracks key occurrence counts for a single segment.
type segmentCounter map[string]uint64

func newSegmentCounter() segmentCounter {
	return make(segmentCounter, 1024)
}

func (s segmentCounter) Add(key []byte) {
	s[string(key)]++
}
