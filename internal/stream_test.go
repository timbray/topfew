package topfew

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestBadStreamReader(t *testing.T) {
	args := []string{}
	c, err := Configure(args)
	if err != nil {
		t.Error("config!")
	}
	cer := newCER("testing stream")
	_, err = fromStream(cer, &c.filter, nil, c.size)
	if err == nil {
		t.Error("survived err from Read")
	}
	if err.Error() != cer.nonce {
		t.Error("got wrong error: " + err.Error())
	}
}

func TestStreamProcessing(t *testing.T) {
	args := []string{"-f", "3", "--vgrep", "FOO"}
	c, err := Configure(args)
	if err != nil {
		t.Error("config: " + err.Error())
	}
	input := "FOO\nBAR\n"
	stringreader := bufio.NewReader(strings.NewReader(input))
	top, err := Run(c, stringreader)
	if err != nil {
		t.Error("Run!")
	}
	if len(top) != 0 {
		t.Error("all input records should have been ignored")
	}
}

func Test1KLinesStream(t *testing.T) {
	file, err := os.Open("../test/data/small")
	if err != nil {
		t.Error("Can't open file")
	}
	//noinspection ALL
	defer file.Close()

	kf := newKeyFinder([]uint{1}, nil)
	f := filters{nil, nil, nil}
	x, err := fromStream(file, &f, kf, 5)
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
		if *kc.Count != uint64(wanted[kc.Key]) {
			t.Error("Wrong count for Key: " + kc.Key)
		}
	}
}
