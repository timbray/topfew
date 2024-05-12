package topfew

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestFieldSeparator(t *testing.T) {
	args := []string{"-p", "tt*", "-f", "2,4"}

	c, err := Configure(args)
	if err != nil {
		t.Error("Config!")
	}

	records := []string{
		"atbttctttttdtttte",
	}
	wanted := []string{
		"b d",
	}
	kf := newKeyFinder(c.fields, c.fieldSeparator, false)
	for i, record := range records {
		got, err := kf.getKey([]byte(record))
		if err != nil {
			t.Error("getKey: " + err.Error())
		}
		if string(got) != wanted[i] {
			t.Errorf("wanted %s got %s", wanted[i], string(got))
		}
	}
	_, err = kf.getKey([]byte("atbtc"))
	if err == nil || err.Error() != NER {
		t.Error("bad error value")
	}
}

func TestQuotedFields(t *testing.T) {
	lines := []string{
		`i577a483c.versanet.de - - [12/Mar/2007:08:03:37 -0800] "GET /ongoing/ongoing.atom HTTP/1.1" 304 - "-" "NetNewsWire/2.1 (Mac OS X; http://ranchero.com/netnewswire/)"`,
		`105.66.1.178 - - [19/Apr/2020:06:38:44 -0700] "-" 408 156 "-" "-"`,
	}
	kf := newKeyFinder([]uint{6}, nil, true)
	for i, line := range lines {
		k, err := kf.getKey([]byte(line))
		if err != nil {
			t.Error("getKey: " + err.Error())
		}
		if i == 0 && string(k) != "GET /ongoing/ongoing.atom HTTP/1.1" {
			t.Errorf("at 0 got %s", string(k))
		}
		if i == 1 && string(k) != "-" {
			t.Errorf("at 1 got %s", string(k))
		}
	}
	// test non-quoted fields with -q

	n5 := uint64(5)
	n3 := uint64(3)
	wanted := []*keyCount{
		{"[12/Mar/2007:08:03:42", &n5},
		{"[12/Mar/2007:08:03:37", &n3},
	}
	args := []string{"-q", "-f", "4", "-n", "2", "../test/data/10lines"}
	c, err := Configure(args)
	if err != nil {
		t.Error("Config!")
	}
	kc, _ := Run(c, nil)
	assertKeyCountsEqual(t, wanted, kc)

	f, err := os.Open("../test/data/10lines")
	if err != nil {
		t.Error("Open: " + err.Error())
	}
	args = []string{"-q", "-f", "7"}
	c, err = Configure(args)
	if err != nil {
		t.Error("Config!")
	}
	kc, _ = Run(c, f)
	fmt.Printf("kc %d\n", len(kc))
}

func TestSpacedFieldSelection(t *testing.T) {
	spacedFields := []string{
		`i577a483c.versanet.de - - [12/Mar/2007:08:03:37 -0800] "GET /ongoing/ongoing.atom HTTP/1.1" 304 - "-" "NetNewsWire/2.1 (Mac OS X; http://ranchero.com/netnewswire/)"`,
		`ln-bas00.csfb.com - - [12/Mar/2007:08:03:37 -0800] "GET /ongoing/ongoing.atom HTTP/1.0" 304 - "-" "Mozilla/3.01 (compatible;)"`,
		`host100.newsgator.com - - [12/Mar/2007:08:03:37 -0800] "GET /ongoing/comments.atom HTTP/1.1" 200 134097 "-" "NewsGatorOnline/2.0 (http://www.newsgator.com; 3 subscribers)"`,
		`pool-141-152-243-156.phil.east.verizon.net - - [12/Mar/2007:08:03:38 -0800] "GET /ongoing/ongoing.atom HTTP/1.1" 304 - "-" "NetNewsWire/2.1.1 (Mac OS X; Lite; http://ranchero.com/netnewswire/)"`,
		`ertpg6e1.nortelnetworks.com - - [12/Mar/2007:08:03:42 -0800] "GET /ongoing/ongoing.atom HTTP/1.1" 304 - "-" "NetNewsWire/2.1.1 (Mac OS X; Lite; http://ranchero.com/netnewswire/)"`,
		`82.153.22.192 - - [12/Mar/2007:08:03:42 -0800] "GET /ongoing/ongoing.rss HTTP/1.1" 301 322 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.0; en-GB; rv:1.8.1.2) Gecko/20070219 Firefox/2.0.0.2"`,
		`css6.csee.usf.edu - - [12/Mar/2007:08:03:42 -0800] "GET /ongoing/When/200x/2007/03/11/Ramirez.png HTTP/1.1" 200 67663 "http://www.tbray.org/ongoing/When/200x/2007/03/11/Misa-Criolla" "endo/1.0 (Mac OS X; ppc i386; http://kula.jp/endo)"`,
		`82.153.22.192 - - [12/Mar/2007:08:03:42 -0800] "GET /ongoing/ongoing.rss HTTP/1.1" 301 327 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.0; en-GB; rv:1.8.1.2) Gecko/20070219 Firefox/2.0.0.2"`,
		`82.153.22.192 - - [12/Mar/2007:08:03:42 -0800] "GET /ongoing/ongoing.atom HTTP/1.1" 304 - "-" "Mozilla/5.0 (Windows; U; Windows NT 5.0; en-GB; rv:1.8.1.2) Gecko/20070219 Firefox/2.0.0.2"`,
		`css6.csee.usf.edu - - [12/Mar/2007:08:03:43 -0800] "GET /ongoing/When/200x/2007/03/11/Misa-Criolla.png HTTP/1.1" 200 79064 "http://www.tbray.org/ongoing/When/200x/2007/03/11/Misa-Criolla" "endo/1.0 (Mac OS X; ppc i386; http://kula.jp/endo)"`,
	}
	for recNum, recString := range spacedFields {
		fields := strings.Split(recString, " ")
		for fieldNum, field := range fields {
			kf := newKeyFinder([]uint{uint(fieldNum + 1)}, nil, false)
			k, err := kf.getKey([]byte(recString))
			if err != nil {
				t.Errorf("getKey! rec %d field %d", recNum, fieldNum)
			}
			if string(k) != field {
				t.Errorf("bad match rec %d field %d wanted %s got %s", recNum, fieldNum, field, string(k))
			}
		}
	}
}
func TestQuotedFieldSelection(t *testing.T) {
	quotedFields := [][]string{
		{"i577a483c.versanet.de", "-", "-", "[12/Mar/2007:08:03:37", "-0800]",
			"GET /ongoing/ongoing.atom HTTP/1.1", "304", "-", "-",
			"NetNewsWire/2.1 (Mac OS X; http://ranchero.com/netnewswire/)"},
		{"ln-bas00.csfb.com", "-", "-", "[12/Mar/2007:08:03:37", "-0800]",
			"GET /ongoing/ongoing.atom HTTP/1.0", "304", "-", "-",
			"Mozilla/3.01 (compatible;)"},
		{"host100.newsgator.com", "-", "-", "[12/Mar/2007:08:03:37", "-0800]",
			"GET /ongoing/comments.atom HTTP/1.1", "200", "134097", "-",
			"NewsGatorOnline/2.0 (http://www.newsgator.com; 3 subscribers)"},
		{"pool-141-152-243-156.phil.east.verizon.net", "-", "-", "[12/Mar/2007:08:03:38", "-0800]",
			"GET /ongoing/ongoing.atom HTTP/1.1", "304", "-", "-",
			"NetNewsWire/2.1.1 (Mac OS X; Lite; http://ranchero.com/netnewswire/)"},
		{"ertpg6e1.nortelnetworks.com", "-", "-", "[12/Mar/2007:08:03:42", "-0800]",
			"GET /ongoing/ongoing.atom HTTP/1.1", "304", "-", "-",
			"NetNewsWire/2.1.1 (Mac OS X; Lite; http://ranchero.com/netnewswire/)"},
		{"82.153.22.192", "-", "-", "[12/Mar/2007:08:03:42", "-0800]",
			"GET /ongoing/ongoing.rss HTTP/1.1", "301", "322", "-",
			"Mozilla/5.0 (Windows; U; Windows NT 5.0; en-GB; rv:1.8.1.2) Gecko/20070219 Firefox/2.0.0.2"},
		{"css6.csee.usf.edu", "-", "-", "[12/Mar/2007:08:03:42", "-0800]",
			"GET /ongoing/When/200x/2007/03/11/Ramirez.png HTTP/1.1", "200", "67663",
			"http://www.tbray.org/ongoing/When/200x/2007/03/11/Misa-Criolla",
			"endo/1.0 (Mac OS X; ppc i386; http://kula.jp/endo)"},
		{"82.153.22.192", "-", "-", "[12/Mar/2007:08:03:42", "-0800]",
			"GET /ongoing/ongoing.rss HTTP/1.1", "301", "327", "-",
			"Mozilla/5.0 (Windows; U; Windows NT 5.0; en-GB; rv:1.8.1.2) Gecko/20070219 Firefox/2.0.0.2"},
		{"82.153.22.192", "-", "-", "[12/Mar/2007:08:03:42", "-0800]",
			"GET /ongoing/ongoing.atom HTTP/1.1", "304", "-", "-",
			"Mozilla/5.0 (Windows; U; Windows NT 5.0; en-GB; rv:1.8.1.2) Gecko/20070219 Firefox/2.0.0.2"},
		{"css6.csee.usf.edu", "-", "-", "[12/Mar/2007:08:03:43", "-0800]",
			"GET /ongoing/When/200x/2007/03/11/Misa-Criolla.png HTTP/1.1", "200", "79064",
			"http://www.tbray.org/ongoing/When/200x/2007/03/11/Misa-Criolla",
			"endo/1.0 (Mac OS X; ppc i386; http://kula.jp/endo)"},
	}

	var fieldNum uint
	for fieldNum = 1; fieldNum <= 10; fieldNum++ {
		kf := newKeyFinder([]uint{fieldNum}, nil, true)
		kf11 := newKeyFinder([]uint{11}, nil, true)
		f, err := os.Open("../test/data/10lines")
		br := bufio.NewReader(f)
		if err != nil {
			t.Error("Open: " + err.Error())
		}
		for recordNum := 0; recordNum < 10; recordNum++ {
			record, err := br.ReadBytes('\n')
			if err != nil {
				t.Error("readBytes: " + err.Error())
			}
			k, err := kf11.getKey(record)
			if err == nil {
				t.Errorf("r11 OK! key <%s>\n", string(k))
			}
			k, _ = kf.getKey(record)
			sk := string(k)
			wanted := quotedFields[recordNum][fieldNum-1]
			if sk != wanted {
				t.Errorf("r[%d] f[%d] wanted %s got %s", recordNum, fieldNum, wanted, sk)
			}
		}
		_ = f.Close()
	}
}

func TestMultiFields(t *testing.T) {
	records12 := []string{
		`a b c d`,
		`a "b" c d`,
		`"a" b "c" d`,
		`"a" "b" "c" d`,
	}
	records13 := []string{
		`a b c d`,
		`"a" b c d`,
		`a b "c" d`,
		`"a" b "c" d`,
	}
	records23 := []string{
		`a b c d`,
		`a "b" c d`,
		`a b "c" d`,
		`a "b" "c" d`,
	}
	records24 := []string{
		`a b c d`,
		`a "b" c d`,
		`a b c "d"`,
		`a "b" "c" "d"`,
		`a b c d    `,
		`a b c "d"  `,
	}
	records34 := []string{
		`a b c d`,
		`a b "c" d`,
		`a b c "d"`,
		`a b "c" "d"`,
		`a b c d    `,
		`a b c "d"  `,
	}
	kf12 := newKeyFinder([]uint{uint(1), uint(2)}, nil, true)
	kf13 := newKeyFinder([]uint{uint(1), uint(3)}, nil, true)
	kf23 := newKeyFinder([]uint{uint(2), uint(3)}, nil, true)
	kf24 := newKeyFinder([]uint{uint(2), uint(4)}, nil, true)
	kf34 := newKeyFinder([]uint{uint(3), uint(4)}, nil, true)

	for _, record := range records12 {
		k, err := kf12.getKey([]byte(record))
		if err != nil {
			t.Errorf("kf12 err on <%s>: %s", record, err.Error())
		} else if string(k) != "a b" {
			t.Errorf("kf12 key on <%s> = <%s>", record, string(k))
		}
	}
	for _, record := range records13 {
		k, err := kf13.getKey([]byte(record))
		if err != nil {
			t.Errorf("kf13 err on <%s>: %s", record, err.Error())
		} else if string(k) != "a c" {
			t.Errorf("kf13 key on <%s> = <%s>", record, string(k))
		}
	}
	for _, record := range records23 {
		k, err := kf23.getKey([]byte(record))
		if err != nil {
			t.Errorf("kf23 err on <%s>: %s", record, err.Error())
		} else if string(k) != "b c" {
			t.Errorf("kf23 key on <%s> = <%s>", record, string(k))
		}
	}
	for _, record := range records24 {
		k, err := kf24.getKey([]byte(record))
		if err != nil {
			t.Errorf("kf24 err on <%s>: <%s>", record, err.Error())
		} else if string(k) != "b d" {
			t.Errorf("kf24 key on <%s> = <%s>", record, string(k))
		}
	}
	for _, record := range records34 {
		k, err := kf34.getKey([]byte(record))
		if err != nil {
			t.Errorf("kf34 err on <%s>: <%s>", record, err.Error())
		}
		if string(k) != "c d" {
			t.Errorf("kf34 key on <%s> = <%s>", record, string(k))
		}
	}
}

func TestPassQuoted(t *testing.T) {
	record := []byte(`foo bar baz         `)
	kf := newKeyFinder([]uint{uint(5)}, nil, true)
	key, err := kf.getKey(record)
	if err == nil {
		t.Errorf("accepted, k=<%s>", key)
	}
	record = []byte(`foo bar "baz`)
	kf = newKeyFinder([]uint{uint(4)}, nil, true)
	key, err = kf.getKey(record)
	if err == nil {
		t.Errorf("accepted 2, k=<%s>", key)
	}

}

func TestGatherQuoted(t *testing.T) {
	record := []byte(`foo bar "baz`)
	kf := newKeyFinder([]uint{uint(3)}, nil, true)
	key, err := kf.getKey(record)
	if err == nil {
		t.Errorf("accepted, k=<%s>", key)
	}
}

func TestCSVSeparator(t *testing.T) {
	args := []string{"-p", ",", "-f", "11"}
	c, err := Configure(args)
	if err != nil {
		t.Error("Config!")
	}
	input, err := os.Open("../test/data/csoc.csv")
	if err != nil {
		t.Error("Open: " + err.Error())
	}
	counts, err := Run(c, input)
	if err != nil {
		t.Error("Run: " + err.Error())
	}
	if len(counts) != 5 {
		t.Errorf("Got %d results, wanted 5", len(counts))
	}
	wantCounts := []uint64{4, 2, 1, 1, 1}
	wantKeys := []string{"50", "-1.97", "amount", "-1.75", "-1.9"}
	for i, count := range counts {
		if *count.Count != wantCounts[i] {
			t.Errorf("Counts[%d] is %d wanted %d", i, *count.Count, wantCounts[i])
		}
		// because for equal values, the sort isn't stable - Counts[2,3,4] are all 1
		if i < 2 && count.Key != wantKeys[i] {
			t.Errorf("Keys[%d] is %s wanted %s", i, count.Key, wantKeys[i])
		}
	}
}

func TestKeyFinder(t *testing.T) {
	var records = []string{
		"a x c",
		"a b c",
		"a b c d e",
	}
	var kf, kf2 *keyFinder

	kf = newKeyFinder(nil, nil, false)
	kf2 = newKeyFinder([]uint{}, nil, false)

	for _, recordString := range records {
		record := []byte(recordString)
		r, err := kf.getKey(record)
		if (err != nil) || !bytes.Equal(r, record) {
			t.Errorf("bad result on nil for %s", record)
		}
		r, err = kf2.getKey(record)
		if (err != nil) || !bytes.Equal(r, record) {
			t.Errorf("bad result on empty for %s", record)
		}
	}

	singles := []string{"x", "b", "b"}
	kf = newKeyFinder([]uint{2}, nil, false)
	for i, record := range records {
		k, err := kf.getKey([]byte(record))
		if err != nil {
			t.Error("KF fail on: " + record)
		} else {
			if string(k) != singles[i] {
				t.Errorf("got '%s' wanted '%s'", string(k), singles[i])
			}
		}
	}

	kf = newKeyFinder([]uint{1, 3}, nil, false)
	for _, recordstring := range records {
		record := []byte(recordstring)
		r, err := kf.getKey(record)
		if err != nil || string(r) != "a c" {
			t.Errorf("wanted a c from %s, got %s", record, r)
		}
	}

	kf = newKeyFinder([]uint{1, 4}, nil, false)
	tooShorts := []string{"a", "a b", "a b c"}
	for _, tooShortString := range tooShorts {
		tooShort := []byte(tooShortString)
		_, err := kf.getKey(tooShort)
		if err == nil {
			t.Errorf("no error on %s", tooShort)
		}
	}
	r, err := kf.getKey([]byte("a b c d"))
	if err != nil || string(r) != "a d" {
		t.Error("border condition")
	}
}
