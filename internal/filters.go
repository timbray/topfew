package topfew

import (
	"regexp"
)

// Sed represents a sed(1) s/a/b/g operation.
type Sed struct {
	ReplaceThis *regexp.Regexp
	WithThat    []byte
}

// Filters contains the filters to be applied prior to top-few computation.
type Filters struct {
	Greps  []*regexp.Regexp
	VGreps []*regexp.Regexp
	Seds   []*Sed
}

// AddSed appends a new Sed operation to the filters.
func (f *Filters) AddSed(replaceThis string, withThat string) error {
	re, err := regexp.Compile(replaceThis)
	if err == nil {
		f.Seds = append(f.Seds, &Sed{re, []byte(withThat)})
	}
	return err
}

// AddGrep appends a new grep/regex to the filters. Only items that match
// this regex will be counted.
func (f *Filters) AddGrep(s string) error {
	re, err := regexp.Compile(s)
	if err == nil {
		f.Greps = append(f.Greps, re)
	}
	return err
}

// AddVgrep appends a new inverse grep/regex to the filters (ala grep -v).
// Only items that don't match the regex will be counted.
func (f *Filters) AddVgrep(s string) error {
	re, err := regexp.Compile(s)
	if err == nil {
		f.VGreps = append(f.VGreps, re)
	}
	return err
}

// FilterRecord returns true if the supplied record passes all the Filter
// criteria.
func (f *Filters) FilterRecord(bytes []byte) bool {
	if f.Greps == nil && f.VGreps == nil {
		return true
	}
	for _, re := range f.Greps {
		if !re.Match(bytes) {
			return false
		}
	}
	for _, re := range f.VGreps {
		if re.Match(bytes) {
			return false
		}
	}
	return true
}

// FilterField returns a key that has had all the sed operations applied to it.
func (f *Filters) FilterField(bytes []byte) []byte {
	for _, sed := range f.Seds {
		bytes = sed.ReplaceThis.ReplaceAll(bytes, sed.WithThat)
	}
	return bytes
}
