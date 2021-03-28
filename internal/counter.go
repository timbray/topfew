package topfew

import (
	"sort"
	"sync"
)

// represents a key's occurrence count
type KeyCount struct {
	Key   string
	Count uint64
}

// represents a bunch of keys and their occurrence counts, with the highest counts tracked.
// threshold represents the minimum count value to qualify for consideration as a top count
// the "top" map represents the keys & counts encountered so far which are higher than threshold
// The hash values are pointers not integers for efficiency reasons so you don't have to update the
//  map[string] mapping, you just update the number the key maps to.

type Counter struct {
	counts    map[string]*uint64
	top       map[string]*uint64
	threshold uint64
	size      int
	lock      sync.Mutex
}

func NewCounter(size int) *Counter {
	t := new(Counter)
	t.size = size
	t.counts = make(map[string]*uint64)
	t.top = make(map[string]*uint64)
	return t
}

func (t *Counter) ConcurrentAddKeys(keys [][]byte) {
	t.lock.Lock()
	defer t.lock.Unlock()
	for _, key := range keys {
		t.Add(key)
	}
}

// note the call with a byte slice rather than the string because of
//  https://github.com/golang/go/commit/f5f5a8b6209f84961687d993b93ea0d397f5d5bf
//  which recognizes the idiom foo[string(someByteSlice)] and bypasses constructing the string;
//  of course we'd rather just say foo[someByteSlice] but that's not legal because Reasons.
func (t *Counter) Add(bytes []byte) {

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
	t.top[string(bytes)] = count

	// has the top set grown enough to compress?
	if len(t.top) < (t.size * 2) {
		return
	}

	// sort the top candidates, shrink the list to the top t.size, put them back in a map
	var topList = t.topAsSortedList()
	topList = topList[0:t.size]
	t.threshold = topList[len(topList)-1].Count
	t.top = make(map[string]*uint64)
	for _, kc := range topList {
		t.top[kc.Key] = &kc.Count
	}
}

func (t *Counter) topAsSortedList() []*KeyCount {
	var topList []*KeyCount
	for key, count := range t.top {
		topList = append(topList, &KeyCount{key, *count})
	}
	sort.Slice(topList, func(k1, k2 int) bool {
		return topList[k1].Count > topList[k2].Count
	})
	return topList
}

func (t *Counter) GetTop() []*KeyCount {
	t.lock.Lock()
	defer t.lock.Unlock()
	topList := t.topAsSortedList()
	if len(topList) > t.size {
		return topList[0:t.size]
	} else {
		return topList
	}
}
