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
		t.Error("Can't open file")
	}
	//noinspection ALL
	defer file.Close()
	table := NewCounter(5)
	re := regexp.MustCompile("\\s+")

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
	x := table.GetTop()
	n4 := uint64(4)
	n5 := uint64(5)
	n6 := uint64(6)
	n7 := uint64(7)
	n8 := uint64(8)

	wanted := []KeyCount{
		{"c", &n8},
		{"g", &n7},
		{"e", &n6},
		{"f", &n5},
		{"a", &n4},
	}
	if len(x) != len(wanted) {
		t.Error("lengths don't match")
	}
	// shouldn't deepEqual do this?
	for i := 0; i < len(wanted); i++ {
		if x[i].Key != wanted[i].Key || *x[i].Count != *wanted[i].Count {
			t.Errorf("Mismatch at index %d", i)
		}
	}

	table = NewCounter(3)
	for _, key := range keys {
		table.Add([]byte(key))
	}
	x = table.GetTop()
	wanted = []KeyCount{
		{"c", &n8},
		{"g", &n7},
		{"e", &n6},
	}
	for i := 0; i < len(wanted); i++ {
		if x[i].Key != wanted[i].Key || *x[i].Count != *wanted[i].Count {
			t.Errorf("Mismatch at index %d", i)
		}
	}

}

func Test_newTable(t *testing.T) {
	table := NewCounter(333)
	top := table.GetTop()
	if len(top) != 0 {
		t.Error("new table should be empty")
	}
}
