package topfew

import (
	"regexp"
)

// sed represents a sed(1) s/a/b/g operation.
type sed struct {
	ReplaceThis *regexp.Regexp
	WithThat    []byte
}

// filters contains the filters to be applied prior to top-few computation.
type filters struct {
	greps  []*regexp.Regexp
	vgreps []*regexp.Regexp
	seds   []*sed
}

// addSed appends a new sed operation to the filters.
func (f *filters) addSed(replaceThis string, withThat string) error {
	re, err := regexp.Compile(replaceThis)
	if err == nil {
		f.seds = append(f.seds, &sed{re, []byte(withThat)})
	}
	return err
}

// addGrep appends a new grep/regex to the filters. Only items that match
// this regex will be counted.
func (f *filters) addGrep(s string) error {
	re, err := regexp.Compile(s)
	if err == nil {
		f.greps = append(f.greps, re)
	}
	return err
}

// addVgrep appends a new inverse grep/regex to the filters (ala grep -v).
// Only items that don't match the regex will be counted.
func (f *filters) addVgrep(s string) error {
	re, err := regexp.Compile(s)
	if err == nil {
		f.vgreps = append(f.vgreps, re)
	}
	return err
}

// filterRecord returns true if the supplied record passes all the filter
// criteria.
func (f *filters) filterRecord(bytes []byte) bool {
	if f.greps == nil && f.vgreps == nil {
		return true
	}
	for _, re := range f.greps {
		if !re.Match(bytes) {
			return false
		}
	}
	for _, re := range f.vgreps {
		if re.Match(bytes) {
			return false
		}
	}
	return true
}

// filterField returns a Key that has had all the sed operations applied to it.
func (f *filters) filterField(bytes []byte) []byte {
	for _, sed := range f.seds {
		bytes = sed.ReplaceThis.ReplaceAll(bytes, sed.WithThat)
	}
	return bytes
}
