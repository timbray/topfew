package topfew

import (
	"sort"
)

// represents a key's occurrence count
type KeyCount struct {
	Key   string
	Count uint64
}

// represents a bunch of keys and their occurrence counts, with the highest counts tracked.
// threshold represents the minimum count value to qualify for consideration as a top count
// the "top" map represents the keys & counts encountered so far which are higher than threshold
type Counter struct {
	counts    map[string]uint64
	top       map[string]uint64
	threshold uint64
	size      int
}

func NewCounter(size uint) *Counter {
	t := new(Counter)
	t.size = int(size)
	t.counts = make(map[string]uint64)
	t.top = make(map[string]uint64)
	return t
}

func (t *Counter) Add(key string) {
	// have we seen this key?
	count, ok := t.counts[key]
	if !ok {
		count = 1
	} else {
		count++
	}
	t.counts[key] = count

	// big enough to be a top candidate?
	if count < t.threshold {
		return
	}
	t.top[key] = count

	// has the top set grown enough to compress?
	if len(t.top) < (t.size * 2) {
		return
	}

	// sort the top candidates, shrink the list to the top t.size, put them back in a map
	var topList = t.topAsSortedList()
	topList = topList[0:t.size]
	t.threshold = topList[len(topList)-1].Count
	t.top = make(map[string]uint64)
	for _, kc := range topList {
		t.top[kc.Key] = kc.Count
	}
}

func (t *Counter) topAsSortedList() []*KeyCount {
	var topList []*KeyCount
	for key, count := range t.top {
		topList = append(topList, &KeyCount{key, count})
	}
	sort.Slice(topList, func(k1, k2 int) bool {
		return topList[k1].Count > topList[k2].Count
	})
	return topList
}

func (t *Counter) GetTop() []*KeyCount {
	topList := t.topAsSortedList()
	if len(topList) > t.size {
		return topList[0:t.size]
	} else {
		return topList
	}
}
