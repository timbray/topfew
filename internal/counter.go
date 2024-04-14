package topfew

import (
	"sort"
)

// keyCount represents a Key's occurrence count.
type keyCount struct {
	Key   string
	Count *uint64
}

// The core idea is that when you read a large number of field values and want to find the N values which
// occur most commonly, you keep a large table of the occurrence counts for each observed value and a small
// table of the top values/counts, and remember the occurrence threshold it takes to get into the top table.
// Then for each value, you increment its count and see if the new count gets it into the current top-values
// list, you add it if it's not already there. The top-values table will grow, so every so often you trim it
// back to size N. After a while, in a large dataset the overwhelming majority of values will either already
// be in the top-values table or not belong there, so that table's membership will be increasingly stable
// and require neither growing nor trimming. When you reach the end of the data, you sort the top-values
// table (trivial, because it's small) and return that. I haven't done a formal analysis but I'm pretty sure
// the computation trends to O(N) in the size of the number of records. Also it's "embarrassingly parallel"
// in the sense that *if* you can access the records in parallel you can do the top-values computation in as
// many parallel threads as the underlying computer can offer.

// counter represents a bunch of keys and their occurrence counts, with the highest counts tracked.
// threshold represents the minimum count value to qualify for consideration as a top count
// the "top" map represents the keys & counts encountered so far which are higher than threshold
// The hash values are pointers not integers for efficiency reasons, so you don't have to update the
// map[string] mapping, you just update the number the Key maps to.
type counter struct {
	counts    map[string]*uint64
	top       map[string]*uint64
	threshold uint64
	size      int
}

// newCounter creates a new empty counter, ready for use. size controls how many top items to track.
func newCounter(size int) *counter {
	t := new(counter)
	t.size = size
	t.counts = make(map[string]*uint64, 1024)
	t.top = make(map[string]*uint64, size*2)
	return t
}

// add one occurrence to the counts for the indicated Key.
func (t *counter) add(bytes []byte) {
	// note the call with a byte slice rather than the string because of
	//  https://github.com/golang/go/commit/f5f5a8b6209f84961687d993b93ea0d397f5d5bf
	//  which recognizes the idiom foo[string(someByteSlice)] and bypasses constructing the string;
	//  of course we'd rather just say foo[someByteSlice] but that's not legal because Reasons.

	// have we seen this Key?
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

func (t *counter) compact() {
	// sort the top candidates, shrink the list to the top t.size, put them back in a map
	var topList = t.topAsSortedList()
	topList = topList[0:t.size]
	t.threshold = *(topList[len(topList)-1].Count)
	t.top = make(map[string]*uint64, t.size*2)
	for _, kc := range topList {
		t.top[kc.Key] = kc.Count
	}
}

func (t *counter) topAsSortedList() []*keyCount {
	topList := make([]*keyCount, 0, len(t.top))
	for key, count := range t.top {
		topList = append(topList, &keyCount{key, count})
	}
	sort.Slice(topList, func(k1, k2 int) bool {
		return *topList[k1].Count > *topList[k2].Count
	})
	return topList
}

// getTop returns the top occurring keys & counts in order of descending count
func (t *counter) getTop() []*keyCount {
	topList := t.topAsSortedList()
	if len(topList) > t.size {
		return topList[0:t.size]
	}
	return topList
}

// merge applies the counts from the SegmentCounter into the counter.
// Once merged, the SegmentCounter should be discarded.
func (t *counter) merge(segCounter segmentCounter) {
	for segKey, segCount := range segCounter {
		// Annoyingly we can't efficiently call add here because we have
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
			// if it wasn't in t.counts then we already know it's not in
			// t.top
			var topKey bool
			if existingKey {
				_, topKey = t.top[segKey]
			}
			if !topKey {
				t.top[segKey] = count
				// has the top set grown enough to compress?
				if len(t.top) >= (t.size * 2) {
					t.compact()
				}
			}
		}
	}
}

// SegmentCounter tracks Key occurrence counts for a single segment.
type segmentCounter map[string]*uint64

func newSegmentCounter() segmentCounter {
	return make(segmentCounter, 1024)
}

func (s segmentCounter) add(key []byte) {
	count, ok := s[string(key)]
	if !ok {
		var one uint64 = 1
		count = &one // a little surprised this works, i.e. you can give a local variable permanent life…
		s[string(key)] = count
	} else {
		*count++
	}
}
