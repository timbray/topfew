package topfew

import (
	"os"
	"testing"
)

func TestStreamAndFile(t *testing.T) {
	args := []string{"-f", "7", "../test/data/apache-50k.txt"}

	type results struct {
		count uint64
		key   string
	}
	wanted1 := []results{
		{9946, "/ongoing/ongoing.atom"},
		{1460, "/resources/image-set-1x.png"},
		{722, "/ongoing/resources/green-24x24.jpg"},
		{585, "/ongoing/When/202x/2024/04/01/OSQI"},
		{565, "/ongoing/"},
		{532, "/media/content/test.mp4"},
		{458, "/html/syntax/speculative-parsing/resources/stash.py?action=put&uuid={{GET[uuid]}}&encodingcheck=%C4%9E"},
		{439, "/media/content/test.wav"},
		{418, "/images/mozilla-banner.gif"},
		{417, "/ongoing/support/cat.png"},
	}

	config, err := Configure(args)
	if err != nil {
		t.Error("Configure!")
	}
	kc, err := Run(config, nil)
	if err != nil {
		t.Error("Run! " + err.Error())
	}
	for i, res := range kc {
		if *res.Count != wanted1[i].count || res.Key != wanted1[i].key {
			t.Errorf("Missmatch at %d: k %s/%s, c %d/%d", i, res.Key, wanted1[i].key,
				*res.Count, wanted1[i].count)
		}
	}

	args = []string{"-f", "7"}
	config, err = Configure(args)
	if err != nil {
		t.Error("Configure!")
	}
	str, err := os.Open("../test/data/apache-50k.txt")
	if err != nil {
		t.Error("Open!")
	}
	kc, err = Run(config, str)
	if err != nil {
		t.Error("Run! " + err.Error())
	}
	for i, res := range kc {
		if *res.Count != wanted1[i].count || res.Key != wanted1[i].key {
			t.Errorf("Missmatch at %d: k %s/%s, c %d/%d", i, res.Key, wanted1[i].key,
				*res.Count, wanted1[i].count)
		}
	}
}

func TestBadFile(t *testing.T) {
	args := []string{"/nosuch"}
	c, err := Configure(args)
	if err != nil {
		t.Error("config!")
	}
	_, err = Run(c, nil)
	if err == nil {
		t.Error("Accepted bogus file")
	}
}
