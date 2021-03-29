package topfew

import (
	"regexp"
)

type Sed struct {
	ReplaceThis *regexp.Regexp
	WithThat    []byte
}
type Filters struct {
	Greps  []*regexp.Regexp
	VGreps []*regexp.Regexp
	Seds   []*Sed
}

func (f *Filters) AddSed(replaceThis string, withThat string) error {
	re, err := regexp.Compile(replaceThis)
	if err == nil {
		f.Seds = append(f.Seds, &Sed{re, []byte(withThat)})
	}
	return err
}
func (f *Filters) AddGrep(s string) error {
	re, err := regexp.Compile(s)
	if err == nil {
		f.Greps = append(f.Greps, re)
	}
	return err
}

func (f *Filters) AddVgrep(s string) error {
	re, err := regexp.Compile(s)
	if err == nil {
		f.VGreps = append(f.VGreps, re)
	}
	return err
}

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

func (f *Filters) FilterField(bytes []byte) []byte {
	for _, sed := range f.Seds {
		bytes = sed.ReplaceThis.ReplaceAll(bytes, sed.WithThat)
	}
	return bytes
}
