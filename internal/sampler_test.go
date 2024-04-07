package topfew

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

type customErrorReader struct {
	nonce string
}

func newCER(s string) *customErrorReader {
	return &customErrorReader{s}
}

func (r *customErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New(r.nonce)
}

func TestSampler(t *testing.T) {
	args := []string{"--sample", "--sed", "foo", "bar", "--grep", "t1", "--vgrep", "t2"}
	c, err := Configure(args)
	if err != nil {
		t.Error("CONFIG!")
	}
	cer := newCER("testing sampler")
	_, err = Run(c, cer)
	if err == nil {
		t.Error("No error on bogus reader")
	}
	if err.Error() != cer.nonce {
		t.Errorf("wanted nonce %s got %s", err.Error(), cer.nonce)
	}

	// capture stdout
	saveStdout := os.Stdout
	defer func() { os.Stdout = saveStdout }()
	readPipe, writePipe, _ := os.Pipe()
	os.Stdout = writePipe
	stash := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, readPipe)
		stash <- buf.String()
	}()

	// we're not testing the workings of the filtering, just that the input gets routed appropriately
	input := "alpha\n" +
		"t1 foo\n" +
		"t2 foo\n"
	inputReader := bufio.NewReader(strings.NewReader(input))
	_, err = Run(c, inputReader)
	if err != nil {
		t.Error("synth read 1")
	}

	args = []string{"--sample", "--fields", "3"}
	c, err = Configure(args)
	if err != nil {
		t.Error("CONFIG!")
	}
	input = "alpha\n"

	inputReader = bufio.NewReader(strings.NewReader(input))
	_, err = Run(c, inputReader)
	if err == nil {
		t.Error("Accepted -f 3 but no fields")
	}

	_ = writePipe.Close()
	written := <-stash
	lines := strings.Split(written, "\n")
	wanted := []string{
		"SED 0: s/foo/bar/",
		"REJECT: alpha",
		"ACCEPT: t1 foo",
		"KEY IN: t1 foo",
		"FILTERED: t1 bar",
		"REJECT: t2 foo",
	}
	for i, wanted := range wanted {
		matched, _ := regexp.MatchString(wanted, lines[i])
		if !matched {
			t.Errorf("%s didn't match %s", lines[i], wanted)
		}
	}
}
