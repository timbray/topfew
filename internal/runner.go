package topfew

import (
	"fmt"
	"io"
	"os"
)

func Run(config *Config, instream io.Reader) ([]*KeyCount, error) {
	// lifted out of main.go to facilitate testing
	var kf = NewKeyFinder(config.Fields)
	var topList []*KeyCount
	var err error

	if config.Fname == "" {
		if config.Sample {
			for i, sed := range config.Filter.Seds {
				fmt.Printf("SED %d: s/%s/%s/\n", i, sed.ReplaceThis, sed.WithThat)
			}
			err = Sample(instream, &config.Filter, kf)
		} else {
			topList, err = FromStream(instream, &config.Filter, kf, config.Size)
		}
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error reading stream: %s\n", err.Error())
			return nil, err
		}
	} else {
		counter := NewCounter(config.Size)
		err = ReadFileInSegments(config.Fname, &config.Filter, counter, kf, config.Width)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", config.Fname, err.Error())
			return nil, err
		}
		topList = counter.GetTop()
	}

	return topList, err
}
