package topfew

import (
	"sort"
)

// KeyCount represents a key's occurrence count.
type KeyCount struct {
	Key   string
	Count *uint64
}

// Counter represents a bunch of keys and their occurrence counts, with the highest counts tracked.
// threshold represents the minimum count value to qualify for consideration as a top count
// the "top" map represents the keys & counts encountered so far which are higher than threshold
// The hash values are pointers not integers for efficiency reasons so you don't have to update the
// map[string] mapping, you just update the number the key maps to.
type Counter struct {
	counts    map[string]*uint64
	top       map[string]*uint64
	threshold uint64
	size      int
}

// NewCounter creates a new empty counter, ready for use. size controls how many top items to track.
func NewCounter(size int) *Counter {
	t := new(Counter)
	t.size = size
	t.counts = make(map[string]*uint64, 1024)
	t.top = make(map[string]*uint64, size*2)
	return t
}

// Add one occurrence to the counts for the indicated key.
func (t *Counter) Add(bytes []byte) {
	// note the call with a byte slice rather than the string because of
	//  https://github.com/golang/go/commit/f5f5a8b6209f84961687d993b93ea0d397f5d5bf
	//  which recognizes the idiom foo[string(someByteSlice)] and bypasses constructing the string;
	//  of course we'd rather just say foo[someByteSlice] but that's not legal because Reasons.

	// have we seen this key?
	count, ok := t.counts[string(bytes)]
	if !ok {
		var one uint64 = 1
		count = &one // a little surprised this works, i.e. you can give a local variable permanent life…
		t.counts[string(bytes)] = count
	} else {
		*count++
	}

	// big enough to be a top candidate?
	if *count < t.threshold {
		return
	}
	_, ok = t.top[string(bytes)]
	if !ok {
		t.top[string(bytes)] = count
		if len(t.top) >= (t.size * 2) {
			t.compact()
		}
	}
}

func (t *Counter) compact() {
	// sort the top candidates, shrink the list to the top t.size, put them back in a map
	var topList = t.topAsSortedList()
	topList = topList[0:t.size]
	t.threshold = *(topList[len(topList)-1].Count)
	t.top = make(map[string]*uint64, t.size*2)
	for _, kc := range topList {
		t.top[kc.Key] = kc.Count
	}
}

func (t *Counter) topAsSortedList() []*KeyCount {
	topList := make([]*KeyCount, 0, len(t.top))
	for key, count := range t.top {
		topList = append(topList, &KeyCount{key, count})
	}
	sort.Slice(topList, func(k1, k2 int) bool {
		return *topList[k1].Count > *topList[k2].Count
	})
	return topList
}

// GetTop returns the top occuring keys & counts in order, with highest count first.
func (t *Counter) GetTop() []*KeyCount {
	topList := t.topAsSortedList()
	if len(topList) > t.size {
		return topList[0:t.size]
	}
	return topList
}

// merge applies the counts from the SegmentCounter into the Counter.
// Once merged, the SegmentCounter should be discarded.
func (t *Counter) merge(segCounter segmentCounter) {
	for segKey, segCount := range segCounter {
		// Annoyingly we can't efficiently call Add here because we have
		// a string not a []byte
		count, existingKey := t.counts[segKey]
		if !existingKey {
			count = segCount
			t.counts[segKey] = segCount
		} else {
			*count += *segCount
		}

		// big enough to be a top candidate?
		if *count >= t.threshold {
			// if it wasn't in t.counts then we already know its not in
			// t.top
			if existingKey {
				_, existingKey = t.top[segKey]
			}
			if !existingKey {
				t.top[segKey] = count
				// has the top set grown enough to compress?
				if len(t.top) >= (t.size * 2) {
					t.compact()
				}
			}
		}
	}
}

// SegmentCounter tracks key occurrence counts for a single segment.
type segmentCounter map[string]*uint64

func newSegmentCounter() segmentCounter {
	return make(segmentCounter, 1024)
}

func (s segmentCounter) Add(key []byte) {
	count, ok := s[string(key)]
	if !ok {
		var one uint64 = 1
		count = &one // a little surprised this works, i.e. you can give a local variable permanent life…
		s[string(key)] = count
	} else {
		*count++
	}
}
