package topfew

import (
	"os"
	"testing"
)

func Test1KLinesStream(t *testing.T) {
	file, err := os.Open("../test/data/access-1k")
	if err != nil {
		t.Error("Can't open file")
	}
	//noinspection ALL
	defer file.Close()

	kf := NewKeyFinder([]uint{1})
	x, err := TopFewFromStream(file, kf, 5)
	if err != nil {
		t.Error("OUCH: " + err.Error())
	}
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
		if kc.Count != uint64(wanted[kc.Key]) {
			t.Error("Wrong count for key: " + kc.Key)
		}
	}
}
