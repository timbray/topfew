package topfew

import (
	"testing"
)

func TestSeds(t *testing.T) {
	lines := []string{
		/*0*/ "96.48.229.116 - - [04/May/2020:06:36:20 -0700] \"GET /ongoing/in-feed.xml HTTP/1.1\" 200 781 \"https://old.tbray.org/ongoing/When/202x/2020/04/29/Leaving-Amazon\" \"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36\"\n",
		/*1*/ "151.225.84.185 - - [04/May/2020:06:47:04 -0700] \"GET /ongoing/ongoing.js HTTP/1.1\" 200 2477 \"https://www.tbray.org/ongoing/When/202x/2020/04/29/Leaving-Amazon\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.5 Safari/605.1.15\"\n",
		/*2*/ "173.173.23.87 - - [04/May/2020:06:47:09 -0700] \"GET /ongoing/When/202x/2020/04/29/Leaving-Amazon HTTP/1.1\" 200 10465 \"https://t.co/oShy4TQisN?amp=1\" \"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.5 Mobile/15E148 Safari/604.1\"\n",
		/*3*/ "203.189.152.127 - - [04/May/2020:06:47:12 -0700] \"GET /favicon.ico HTTP/1.1\" 200 6958 \"-\" \"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0\"\n",
		/*4*/ "54.38.222.160 - - [04/May/2020:06:47:14 -0700] \"GET /ongoing/When/201x/2017/10/26/Working-at-Amazon HTTP/1.1\" 200 10944 \"https://www.tbray.org/ongoing/When/201x/2017/10/26/Working-at-Amazon\" \"WordPress/5.4.1; https://icdomainnames.com\"\n",
		/*5*/ "96.44.24.65 - - [04/May/2020:06:47:32 -0700] \"GET /ongoing/in-feed.xml HTTP/1.1\" 200 781 \"https://www.tbray.org/ongoing/When/202x/2020/04/29/Leaving-Amazon\" \"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148\"\n",
		/*6*/ "172.124.211.165 - - [04/May/2020:06:47:40 -0700] \"GET /ongoing/serif.css HTTP/1.1\" 200 2177 \"https://www.tbray.org/ongoing/When/202x/2020/04/29/Leaving-Amazon\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:74.0) Gecko/20100101 Firefox/74.0\"\n",
	}

	var filter Filters
	err := filter.AddSed("^.*\\[04/May/2020:", "")
	if err != nil {
		t.Error("remove prefix" + err.Error())
	}
	err = filter.AddSed(" .*\n", "")
	if err != nil {
		t.Error("isolate time " + err.Error())
	}
	wanted := []string{"06:36:20", "06:47:04", "06:47:09", "06:47:12", "06:47:14", "06:47:32", "06:47:40"}
	for i, line := range lines {
		got := string(filter.FilterField([]byte(line)))
		if got != wanted[i] {
			t.Errorf("Wanted [%s], got [%s]", wanted[i], got)
		}
	}
}

func TestFilterCombos(t *testing.T) {

	lines := []string{
		/*0*/ "96.48.229.116 - - [04/May/2020:06:36:20 -0700] \"GET /ongoing/in-feed.xml HTTP/1.1\" 200 781 \"https://old.tbray.org/ongoing/When/202x/2020/04/29/Leaving-Amazon\" \"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36\"\n",
		/*1*/ "151.225.84.185 - - [04/May/2020:06:47:04 -0700] \"GET /ongoing/ongoing.js HTTP/1.1\" 200 2477 \"https://www.tbray.org/ongoing/When/202x/2020/04/29/Leaving-Amazon\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.5 Safari/605.1.15\"\n",
		/*2*/ "173.173.23.87 - - [04/May/2020:06:47:09 -0700] \"GET /ongoing/When/202x/2020/04/29/Leaving-Amazon HTTP/1.1\" 200 10465 \"https://t.co/oShy4TQisN?amp=1\" \"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.5 Mobile/15E148 Safari/604.1\"\n",
		/*3*/ "203.189.152.127 - - [04/May/2020:06:47:12 -0700] \"GET /favicon.ico HTTP/1.1\" 200 6958 \"-\" \"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0\"\n",
		/*4*/ "54.38.222.160 - - [04/May/2020:06:47:14 -0700] \"GET /ongoing/When/201x/2017/10/26/Working-at-Amazon HTTP/1.1\" 200 10944 \"https://www.tbray.org/ongoing/When/201x/2017/10/26/Working-at-Amazon\" \"WordPress/5.4.1; https://icdomainnames.com\"\n",
		/*5*/ "96.44.24.65 - - [04/May/2020:06:47:32 -0700] \"GET /ongoing/in-feed.xml HTTP/1.1\" 200 781 \"https://www.tbray.org/ongoing/When/202x/2020/04/29/Leaving-Amazon\" \"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148\"\n",
		/*6*/ "172.124.211.165 - - [04/May/2020:06:47:40 -0700] \"GET /ongoing/serif.css HTTP/1.1\" 200 2177 \"https://www.tbray.org/ongoing/When/202x/2020/04/29/Leaving-Amazon\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:74.0) Gecko/20100101 Firefox/74.0\"\n",
	}

	/*
		fields := []string{
			"foo",         // 0
			"bar",         // 1
			"donkey",      // 2
			"baz",         // 3
			"risk",        // 4
			"failure",     // 5
			"dunk",        // 6
		}
	*/

	wantCSS := "\"GET \\S+\\.css "
	var err error

	var recordFilter Filters
	err = recordFilter.AddGrep(wantCSS)
	if err != nil {
		t.Error("AddGrep " + err.Error())
	}

	for i, line := range lines {
		matched := recordFilter.FilterRecord([]byte(line))
		if matched {
			if i != 6 {
				t.Error("Matched " + lines[i])
			}
		} else {
			if i == 6 {
				t.Error("Didn't match " + lines[6])
			}
		}
	}

	recordFilter = Filters{nil, nil, nil}
	err = recordFilter.AddVgrep(wantCSS)
	if err != nil {
		t.Error("AddVgrep" + err.Error())
	}
	for i, line := range lines {
		matched := recordFilter.FilterRecord([]byte(line))
		if !matched {
			if i != 6 {
				t.Error("Didn't match " + lines[i])
			}
		} else {
			if i == 6 {
				t.Error("Matched " + lines[6])
			}
		}
	}

	recordFilter = Filters{nil, nil, nil}
	err = recordFilter.AddGrep("\"GET \\S*-Amazon ")
	if err != nil {
		t.Error("AddGrep " + err.Error())
	}
	err = recordFilter.AddVgrep("Leaving-")
	if err != nil {
		t.Error("AddVgrep " + err.Error())
	}
	for i, line := range lines {
		matched := recordFilter.FilterRecord([]byte(line))
		if matched {
			if i != 4 {
				t.Error("Matched " + lines[i])
			}
		} else {
			if i == 4 {
				t.Error("Didn't match " + lines[6])
			}
		}
	}

	recordFilter = Filters{nil, nil, nil}
	err = recordFilter.AddGrep("\"GET \\S+-Amazon ")
	if err != nil {
		t.Error("AddGrep " + err.Error())
	}
	err = recordFilter.AddGrep("^54.38.222.160 ")
	if err != nil {
		t.Error("AddGrep " + err.Error())
	}
	matched := 0
	for _, line := range lines {
		if recordFilter.FilterRecord([]byte(line)) {
			matched++
		}
	}
	if matched != 1 {
		t.Errorf("Matched was %d wanted 4", matched)
	}

}
