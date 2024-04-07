package topfew

import (
	"bufio"
	"os"
	"regexp"
	"testing"
)

func Test1KLines(t *testing.T) {
	file, err := os.Open("../test/data/small")
	if err != nil {
		t.Fatalf("Can't open file %v", err)
	}
	//noinspection ALL
	defer file.Close()
	table := NewCounter(5)
	re := regexp.MustCompile(`\s+`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := re.Split(scanner.Text(), 2)
		table.Add([]byte(fields[0]))
	}

	if err := scanner.Err(); err != nil {
		t.Error(err.Error())
	}
	x := table.GetTop()

	var wanted = map[string]int{
		"96.48.229.116":   74,
		"71.227.232.164":  24,
		"122.169.54.96":   13,
		"185.156.175.199": 13,
		"203.189.152.127": 13,
	}

	if len(x) != len(wanted) {
		t.Error("lengths don't match")
	}

	for _, kc := range x {
		if *kc.Count != uint64(wanted[kc.Key]) {
			t.Error("Wrong count for key: " + kc.Key)
		}
	}
}

func TestTable_Add(t *testing.T) {
	table := NewCounter(5)
	keys := []string{
		"a", "b", "c", "d", "e", "f", "g", "h",
		"a", "b", "c", "d", "e", "f", "g",
		"a", "c", "d", "e", "f", "g",
		"a", "c", "e", "f", "g",
		"c", "e", "f", "g",
		"c", "e", "g",
		"c", "g",
		"c"}
	for _, key := range keys {
		table.Add([]byte(key))
	}
	n4 := uint64(4)
	n5 := uint64(5)
	n6 := uint64(6)
	n7 := uint64(7)
	n8 := uint64(8)

	wanted := []*KeyCount{
		{"c", &n8},
		{"g", &n7},
		{"e", &n6},
		{"f", &n5},
		{"a", &n4},
	}
	assertKeyCountsEqual(t, wanted, table.GetTop())

	table = NewCounter(3)
	for _, key := range keys {
		table.Add([]byte(key))
	}
	wanted = []*KeyCount{
		{"c", &n8},
		{"g", &n7},
		{"e", &n6},
	}
	assertKeyCountsEqual(t, wanted, table.GetTop())
}

func Test_newTable(t *testing.T) {
	table := NewCounter(333)
	top := table.GetTop()
	if len(top) != 0 {
		t.Error("new table should be empty")
	}
}

func Test_Merge(t *testing.T) {
	a := NewCounter(10)
	b := newSegmentCounter()
	c := newSegmentCounter()
	for i := 0; i < 50; i++ {
		b.Add([]byte{byte('A')})
		b.Add([]byte{byte('B')})
		c.Add([]byte{byte('C')})
		c.Add([]byte{byte('A')})
	}
	c.Add([]byte{byte('C')})
	a.merge(b)
	a.merge(c)
	exp := []*KeyCount{
		{"A", pv(100)}, {"C", pv(51)}, {"B", pv(50)},
	}
	assertKeyCountsEqual(t, exp, a.GetTop())
}

func pv(v uint64) *uint64 {
	return &v
}

func assertKeyCountsEqual(t *testing.T, exp []*KeyCount, act []*KeyCount) {
	t.Helper()
	if len(exp) != len(act) {
		t.Errorf("Expecting %d results, but got %d", len(exp), len(act))
	}
	for i := 0; i < min(len(exp), len(act)); i++ {
		if exp[i].Key != act[i].Key {
			t.Errorf("Unexpected key %v at index %d, expecting %v", act[i].Key, i, exp[i].Key)
		}
		if *exp[i].Count != *act[i].Count {
			t.Errorf("Unexpected count of %d at index %d, expecting %d", *act[i].Count, i, *exp[i].Count)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
